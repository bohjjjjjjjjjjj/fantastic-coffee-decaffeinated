package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

const maxMediaSize = 8 << 20

var mediaIDRegex = regexp.MustCompile(`^med_[a-zA-Z0-9\-]+\.(png|jpg|jpeg|gif|webp)$`)

const (
	mimePNG  = "image/png"
	mimeJPEG = "image/jpeg"
	mimeGIF  = "image/gif"
	mimeWEBP = "image/webp"
)

var extByMIME = map[string]string{
	mimePNG:  "png",
	mimeJPEG: "jpg",
	mimeGIF:  "gif",
	mimeWEBP: "webp",
}

var mimeByExt = map[string]string{
	"png":  mimePNG,
	"jpg":  mimeJPEG,
	"jpeg": mimeJPEG,
	"gif":  mimeGIF,
	"webp": mimeWEBP,
}

type mediaResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (rt *_router) uploadMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if _, _, ok := rt.auth(w, r); !ok {
		return
	}

	body := http.MaxBytesReader(w, r.Body, maxMediaSize)
	data, err := io.ReadAll(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Immagine troppo grande o illeggibile (max 8 MB)")
		return
	}
	if len(data) == 0 {
		writeError(w, http.StatusBadRequest, "Nessun contenuto caricato")
		return
	}

	ext, ok := extByMIME[detectImageMIME(data)]
	if !ok {
		writeError(w, http.StatusBadRequest, "Formato non supportato (usa PNG, JPEG, GIF o WebP)")
		return
	}

	mUUID, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadMedia: UUID")
		writeError(w, http.StatusInternalServerError, "Errore interno")
		return
	}
	mediaID := "med_" + mUUID.String() + "." + ext

	if err := os.MkdirAll(rt.mediaPath, 0o750); err != nil {
		ctx.Logger.WithError(err).Error("uploadMedia: mkdir")
		writeError(w, http.StatusInternalServerError, "Errore nel salvataggio")
		return
	}
	if err := os.WriteFile(filepath.Join(rt.mediaPath, mediaID), data, 0o600); err != nil {
		ctx.Logger.WithError(err).Error("uploadMedia: write")
		writeError(w, http.StatusInternalServerError, "Errore nel salvataggio")
		return
	}

	writeJSON(w, http.StatusCreated, mediaResponse{ID: mediaID, URL: "/media/" + mediaID})
}

func (rt *_router) getMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	mediaID := ps.ByName("mediaId")
	if !mediaIDRegex.MatchString(mediaID) {
		writeError(w, http.StatusNotFound, "Immagine non trovata")
		return
	}

	data, err := os.ReadFile(filepath.Join(rt.mediaPath, mediaID))
	if err != nil {
		writeError(w, http.StatusNotFound, "Immagine non trovata")
		return
	}

	w.Header().Set("Content-Type", mimeByExt[filepath.Ext(mediaID)[1:]])
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func detectImageMIME(data []byte) string {
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return mimeWEBP
	}
	mime := http.DetectContentType(data)
	for i := 0; i < len(mime); i++ {
		if mime[i] == ';' {
			return mime[:i]
		}
	}
	return mime
}

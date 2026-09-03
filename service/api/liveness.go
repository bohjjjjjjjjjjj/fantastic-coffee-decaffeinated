package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// GET /liveness -> liveness
// Sonda di liveness: 204 se il servizio (e il database) rispondono, 500 altrimenti.
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if err := rt.db.Ping(); err != nil {
		writeError(w, http.StatusInternalServerError, "Servizio non disponibile")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

<template>
  <!-- min-height: 0 e' necessario: un flex item ha min-height auto e non si
       restringe sotto il proprio contenuto, quindi senza questo il riquadro
       cresce invece di scrollare e spinge il campo di scrittura fuori schermo. -->
  <div ref="scroll" class="flex-grow-1 p-3 overflow-auto bg-light" style="min-height: 0;">
    <p v-if="!messages.length" class="text-muted text-center mt-4">
      Nessun messaggio. Scrivi qualcosa per iniziare.
    </p>

    <div
      v-for="msg in messages"
      :key="msg.id"
      class="d-flex mb-3 gap-2 align-items-end"
      :class="isMine(msg) ? 'justify-content-end' : 'justify-content-start'"
    >
      <!-- Avatar del mittente: solo sui messaggi altrui, il proprio e' implicito -->
      <Avatar
        v-if="!isMine(msg)"
        :src="msg.senderPhotoUrl"
        :name="msg.senderUsername"
        :size="34"
      />
      <div
        class="p-2 px-3 rounded text-break position-relative"
        :class="isMine(msg) ? 'bg-primary text-white' : 'bg-white border'"
        style="max-width: 75%; min-width: 160px;"
      >
        <small
          class="d-block mb-1"
          :class="isMine(msg) ? 'text-white-50' : 'text-muted'"
        >{{ msg.senderUsername }}</small>

        <!-- Messaggio inoltrato -->
        <div
          v-if="msg.isForwarded"
          class="small fst-italic mb-1 d-flex align-items-center gap-1"
          :class="isMine(msg) ? 'text-white-50' : 'text-muted'"
        >
          <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
            <path d="M9.5 3.5v2.2C4.9 6 2.6 8.8 2 13c1.6-2.3 3.9-3.3 7.5-3.3v2.3l4.5-4.4-4.5-4.1z" />
          </svg>
          Inoltrato
        </div>

        <!-- Risposta a un altro messaggio -->
        <div
          v-if="msg.replyToMessageId"
          class="small fst-italic border-start ps-2 mb-1 opacity-75 d-flex align-items-center gap-1"
        >
          <span>&#8617;</span>
          <svg
            v-if="replyIsPhoto(msg.replyToMessageId)"
            width="13" height="13" viewBox="0 0 16 16" fill="currentColor"
            aria-label="immagine" role="img"
          >
            <path d="M1.5 2h13A1.5 1.5 0 0 1 16 3.5v9A1.5 1.5 0 0 1 14.5 14h-13A1.5 1.5 0 0 1 0 12.5v-9A1.5 1.5 0 0 1 1.5 2zm12 1h-11a.5.5 0 0 0-.5.5v7l3.2-3.1a.5.5 0 0 1 .68 0L9 11l2.1-2a.5.5 0 0 1 .68 0L14 11V3.5a.5.5 0 0 0-.5-.5zM5 5.5a1.2 1.2 0 1 1-2.4 0 1.2 1.2 0 0 1 2.4 0z" />
          </svg>
          <span>{{ replyPreview(msg.replyToMessageId) }}</span>
        </div>

        <!-- Un messaggio puo' contenere immagine, testo o entrambi -->
        <div v-if="msg.content && msg.content.photoUrl" class="mb-1">
          <img
            :src="resolveMedia(msg.content.photoUrl)"
            alt="immagine allegata"
            class="img-fluid rounded"
            style="max-height: 320px;"
          >
        </div>
        <div v-if="msg.content && msg.content.text">{{ msg.content.text }}</div>

        <div class="d-flex align-items-center gap-2 mt-1">
          <small
            class="small"
            :class="isMine(msg) ? 'text-white-50' : 'text-muted'"
          >{{ formatTime(msg.dataSent) }}</small>
          <small
            v-if="isMine(msg)"
            class="small"
            :class="msg.status.value === 'read' ? 'fw-bold' : ''"
            :title="statusLabel(msg.status.value)"
          >{{ statusIcon(msg.status.value) }}</small>
        </div>

        <!-- Reazioni: emoji + autore, cosi' si vede chi ha reagito -->
        <div v-if="msg.reaction && msg.reaction.length" class="mt-1 d-flex flex-wrap gap-1">
          <span
            v-for="r in msg.reaction"
            :key="r.id"
            class="badge bg-light text-dark border d-inline-flex align-items-center gap-1"
            :title="`Reazione di ${r.username}`"
            role="button"
            @click="onReactionClick(msg, r)"
          >
            <span>{{ r.reactionType }}</span>
            <span class="fw-normal">{{ r.username }}</span>
          </span>
        </div>

        <div class="btn-group btn-group-sm mt-2 d-block">
          <button type="button" class="btn btn-sm btn-outline-secondary py-0" @click="$emit('reply', msg)">
            Rispondi
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary py-0" @click="$emit('forward', msg)">
            Inoltra
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary py-0" @click="$emit('react', msg)">
            Reagisci
          </button>
          <button
            v-if="isMine(msg)"
            type="button"
            class="btn btn-sm btn-outline-danger py-0"
            @click="$emit('delete', msg)"
          >
            Elimina
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mediaURL } from '../services/axios.js'
import Avatar from './Avatar.vue'

export default {
  name: 'MessageThread',
  components: { Avatar },
  props: {
    messages: { type: Array, default: () => [] },
    myUsername: { type: String, default: '' }
  },
  emits: ['reply', 'forward', 'react', 'delete', 'toggle-reaction'],
  watch: {
    messages() {
      this.$nextTick(this.scrollToBottom)
    }
  },
  mounted() {
    this.scrollToBottom()
  },
  methods: {
    resolveMedia(url) {
      return mediaURL(url)
    },
    isMine(msg) {
      return msg.senderUsername === this.myUsername
    },
    findMessage(id) {
      return this.messages.find((x) => x.id === id)
    },
    replyIsPhoto(id) {
      const m = this.findMessage(id)
      return !!(m && m.content && m.content.photoUrl)
    },
    replyPreview(id) {
      const m = this.findMessage(id)
      if (!m) return 'messaggio'
      if (m.content && m.content.text) return m.content.text
      if (m.content && m.content.photoUrl) return 'Foto'
      return ''
    },
    onReactionClick(msg, reaction) {
      this.$emit('toggle-reaction', { msg, reaction })
    },
    formatTime(iso) {
      if (!iso) return ''
      const d = new Date(iso)
      return Number.isNaN(d.getTime()) ? '' : d.toLocaleString()
    },
    // Una spunta finche' il messaggio non e' stato letto da tutti i
    // destinatari, due spunte quando tutti l'hanno aperto.
    statusIcon(value) {
      if (value === 'read') return '✓✓'
      if (value === 'failed') return '⚠'
      return '✓'
    },
    statusLabel(value) {
      if (value === 'read') return 'Letto'
      if (value === 'delivered') return 'Consegnato'
      if (value === 'failed') return 'Errore'
      return 'Inviato'
    },
    scrollToBottom() {
      const el = this.$refs.scroll
      if (el) el.scrollTop = el.scrollHeight
    }
  }
}
</script>

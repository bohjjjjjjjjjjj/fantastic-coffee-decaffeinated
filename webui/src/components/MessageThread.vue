<template>
  <!-- min-height: 0 e' necessario: un flex item ha min-height auto e non si
       restringe sotto il proprio contenuto, quindi senza questo il riquadro
       cresce invece di scrollare e spinge il campo di scrittura fuori schermo. -->
  <div ref="scroll" class="flex-grow-1 p-3 overflow-auto wt-thread" style="min-height: 0;">
    <p v-if="!messages.length" class="text-muted text-center mt-4">
      Nessun messaggio. Scrivi qualcosa per iniziare.
    </p>

    <template v-for="(msg, i) in messages" :key="msg.id">
      <!-- Separatore di giornata: compare solo quando cambia il giorno. -->
      <div v-if="isNewDay(i)" class="text-center my-3 wt-day">
        <span class="badge rounded-pill wt-day-chip">{{ dayLabel(msg.dataSent) }}</span>
      </div>

      <div
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
          class="p-2 px-3 rounded text-break position-relative wt-bubble"
          :class="isMine(msg) ? 'wt-bubble-mine' : 'wt-bubble-other'"
          :style="bubbleStyle(msg)"
        >
          <button
            type="button"
            class="wt-dots btn btn-sm p-0 border-0 text-muted"
            :aria-label="'Azioni sul messaggio di ' + msg.senderUsername"
            :aria-expanded="openMenuId === msg.id ? 'true' : 'false'"
            @click.stop="toggleMenu(msg.id)"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
              <path d="M8 3a1.3 1.3 0 1 1 0-2.6A1.3 1.3 0 0 1 8 3zm0 6.3a1.3 1.3 0 1 1 0-2.6 1.3 1.3 0 0 1 0 2.6zm0 6.3a1.3 1.3 0 1 1 0-2.6 1.3 1.3 0 0 1 0 2.6z" />
            </svg>
          </button>

          <!-- Backdrop trasparente: un clic fuori chiude il menu, senza
               listener globali da registrare e rimuovere a mano. -->
          <div v-if="openMenuId === msg.id" class="wt-menu-backdrop" @click="closeMenu" />
          <div
            v-if="openMenuId === msg.id"
            ref="menu"
            class="wt-menu list-group shadow-sm"
            :class="isMine(msg) ? 'wt-menu-right' : 'wt-menu-left'"
          >
            <button type="button" class="list-group-item list-group-item-action py-1 px-3" @click="runAction('reply', msg)">
              Rispondi
            </button>
            <button type="button" class="list-group-item list-group-item-action py-1 px-3" @click="runAction('forward', msg)">
              Inoltra
            </button>
            <button type="button" class="list-group-item list-group-item-action py-1 px-3" @click="runAction('react', msg)">
              Reazione
            </button>
            <button
              v-if="isMine(msg)"
              type="button"
              class="list-group-item list-group-item-action py-1 px-3 text-danger"
              @click="runAction('delete', msg)"
            >
              Elimina
            </button>
          </div>

          <!-- Nei gruppi serve sapere chi scrive; sui propri messaggi no. -->
          <small
            v-if="isGroup && !isMine(msg)"
            class="d-block mb-1 fw-semibold text-muted"
          >{{ msg.senderUsername }}</small>

          <!-- Messaggio inoltrato -->
          <div
            v-if="msg.isForwarded"
            class="small fst-italic mb-1 d-flex align-items-center gap-1 text-muted"
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
            <small class="small text-muted">{{ formatTime(msg.dataSent) }}</small>
            <small
              v-if="isMine(msg)"
              class="small text-muted"
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
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import { mediaURL } from '../services/axios.js'
import Avatar from './Avatar.vue'

// Tinte chiare per distinguere i mittenti nei gruppi.
const BUBBLE_PALETTE = [
  '#DDE8FF', '#E4DDFB', '#E8E0F7', '#FBDDEB', '#FBDEE0',
  '#FDE8D5', '#DCEFE4', '#D8F3EC', '#D7F1F8', '#E5E7E9',
  '#EAF0D8', '#FBF0D2'
]

export default {
  name: 'MessageThread',
  components: { Avatar },
  props: {
    messages: { type: Array, default: () => [] },
    myUsername: { type: String, default: '' },
    isGroup: { type: Boolean, default: false }
  },
  emits: ['reply', 'forward', 'react', 'delete', 'toggle-reaction'],
  data() {
    return {
      // Un solo menu aperto alla volta.
      openMenuId: null
    }
  },
  computed: {
    // I colori sono assegnati per conversazione, nell'ordine in cui i mittenti
    // compaiono: così due utenti non ricevono mai la stessa tinta, cosa che un
    // hash sul nome non potrebbe garantire.
    senderColors() {
      const map = {}
      let next = 0
      for (const m of this.messages) {
        const u = m.senderUsername
        if (!u || u === this.myUsername || map[u] !== undefined) continue
        map[u] = BUBBLE_PALETTE[next % BUBBLE_PALETTE.length]
        next += 1
      }
      return map
    }
  },
  watch: {
    messages() {
      this.$nextTick(this.scrollToBottom)
    }
  },
  mounted() {
    this.scrollToBottom()
  },
  methods: {
    // Nei gruppi ogni mittente ha la sua tinta; nelle chat a due resta il
    // colore unico, e i propri messaggi non cambiano mai.
    bubbleStyle(msg) {
      if (!this.isGroup || this.isMine(msg)) return {}
      const c = this.senderColors[msg.senderUsername]
      return c ? { backgroundColor: c } : {}
    },
    toggleMenu(id) {
      this.openMenuId = this.openMenuId === id ? null : id
      if (this.openMenuId) {
        // Il contenitore dei messaggi ha overflow-auto: senza questo il menu
        // dell'ultimo messaggio verrebbe tagliato in basso.
        this.$nextTick(() => {
          const el = Array.isArray(this.$refs.menu) ? this.$refs.menu[0] : this.$refs.menu
          if (el && el.scrollIntoView) el.scrollIntoView({ block: 'nearest' })
        })
      }
    },
    closeMenu() {
      this.openMenuId = null
    },
    // Riusa gli stessi eventi di prima: la logica del padre non cambia.
    runAction(action, msg) {
      this.closeMenu()
      this.$emit(action, msg)
    },
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
    // Nella bolla resta solo l'ora: la data sta nel separatore di giornata.
    formatTime(iso) {
      if (!iso) return ''
      const d = new Date(iso)
      return Number.isNaN(d.getTime())
        ? ''
        : d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    },
    dayKey(iso) {
      const d = new Date(iso)
      return Number.isNaN(d.getTime()) ? '' : d.toDateString()
    },
    isNewDay(index) {
      if (index === 0) return true
      return this.dayKey(this.messages[index].dataSent) !== this.dayKey(this.messages[index - 1].dataSent)
    },
    dayLabel(iso) {
      const d = new Date(iso)
      if (Number.isNaN(d.getTime())) return ''
      const oggi = new Date()
      const ieri = new Date()
      ieri.setDate(ieri.getDate() - 1)
      if (d.toDateString() === oggi.toDateString()) return 'Oggi'
      if (d.toDateString() === ieri.toDateString()) return 'Ieri'
      return d.toLocaleDateString([], { day: 'numeric', month: 'long', year: 'numeric' })
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

<style scoped>
/* Sfondo del riquadro messaggi: usa il colore di tema definito in App.vue. */
.wt-thread {
  background-color: var(--wt-bg);
}

/* Chip della giornata: resta agganciato in alto mentre si scorre. */
.wt-day {
  position: sticky;
  top: 0;
  z-index: 4;
}
.wt-day-chip {
  background-color: rgba(255, 255, 255, 0.92);
  color: #444;
  border: 1px solid rgba(0, 0, 0, 0.08);
  font-weight: 500;
}

/* Spazio a destra riservato al pulsante tre puntini. */
.wt-bubble {
  max-width: 75%;
  min-width: 10rem;
  padding-right: 2rem !important;
}

/* Tre puntini: discreti, ancorati in alto a destra nella bolla. */
.wt-dots {
  position: absolute;
  top: 0.25rem;
  right: 0.4rem;
  line-height: 1;
  opacity: 0.55;
}
.wt-dots:hover,
.wt-dots[aria-expanded='true'] {
  opacity: 1;
}

/* Il menu si apre sotto il pulsante e resta dentro la larghezza della bolla,
   quindi non esce mai dallo schermo. */
.wt-menu {
  position: absolute;
  top: 1.9rem;
  z-index: 3;
  min-width: 9rem;
  border-radius: 0.5rem;
  overflow: hidden;
  background-color: #fff;
}
.wt-menu-right {
  right: 0.4rem;
}
.wt-menu-left {
  left: 0.4rem;
}
.wt-menu-backdrop {
  position: fixed;
  inset: 0;
  z-index: 2;
}

/* Bolle: entrambe su fondo chiaro, quindi testo scuro e bordo tenue. */
.wt-bubble-mine,
.wt-bubble-other {
  border: 1px solid rgba(0, 0, 0, 0.08);
}
.wt-bubble-mine {
  background-color: #F2EBDF;
}
.wt-bubble-other {
  background-color: #E6F2DF;
}
</style>

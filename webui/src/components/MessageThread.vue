<template>
  <div ref="scroll" class="flex-grow-1 p-3 overflow-auto bg-light">
    <p v-if="!messages.length" class="text-muted text-center mt-4">
      Nessun messaggio. Scrivi qualcosa per iniziare.
    </p>

    <div
      v-for="msg in messages"
      :key="msg.id"
      class="d-flex flex-column mb-3"
      :class="isMine(msg) ? 'align-items-end' : 'align-items-start'"
    >
      <div
        class="p-2 px-3 rounded text-break position-relative"
        :class="isMine(msg) ? 'bg-primary text-white' : 'bg-white border'"
        style="max-width: 75%; min-width: 140px;"
      >
        <small
          class="d-block mb-1"
          :class="isMine(msg) ? 'text-white-50' : 'text-muted'"
        >{{ msg.senderUsername }}</small>

        <div
          v-if="msg.replyToMessageId"
          class="small fst-italic border-start ps-2 mb-1 opacity-75"
        >
          &#8617; {{ replyPreview(msg.replyToMessageId) }}
        </div>

        <div v-if="msg.content && msg.content.photoUrl">
          <img :src="msg.content.photoUrl" alt="foto" class="img-fluid rounded">
        </div>
        <div v-else>{{ msg.content ? msg.content.text : '' }}</div>

        <div class="d-flex align-items-center gap-2 mt-1">
          <small
            :class="isMine(msg) ? 'text-white-50' : 'text-muted'"
            class="small"
          >{{ formatTime(msg.dataSent) }}</small>
          <small v-if="isMine(msg)" class="small" :title="msg.status.value">
            {{ statusIcon(msg.status.value) }}
          </small>
        </div>

        <div v-if="msg.reaction && msg.reaction.length" class="mt-1">
          <span
            v-for="r in msg.reaction"
            :key="r.id"
            class="badge bg-light text-dark border me-1"
            :title="r.username"
            role="button"
            @click="onReactionClick(msg, r)"
          >{{ r.reactionType }}</span>
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
export default {
  name: 'MessageThread',
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
    isMine(msg) {
      return msg.senderUsername === this.myUsername
    },
    replyPreview(id) {
      const m = this.messages.find((x) => x.id === id)
      if (!m) return 'messaggio'
      return m.content && m.content.photoUrl ? 'Foto' : (m.content ? m.content.text : '')
    },
    onReactionClick(msg, reaction) {
      this.$emit('toggle-reaction', { msg, reaction })
    },
    formatTime(iso) {
      if (!iso) return ''
      const d = new Date(iso)
      return Number.isNaN(d.getTime()) ? '' : d.toLocaleString()
    },
    statusIcon(value) {
      if (value === 'read') return '✓✓'
      if (value === 'delivered') return '✓✓'
      if (value === 'failed') return '⚠'
      return '✓'
    },
    scrollToBottom() {
      const el = this.$refs.scroll
      if (el) el.scrollTop = el.scrollHeight
    }
  }
}
</script>

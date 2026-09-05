<template>
  <div class="list-group overflow-auto flex-grow-1">
    <p v-if="!conversations.length" class="text-muted small px-2 mt-2">
      {{ emptyText }}
    </p>
    <button
      v-for="conv in sorted"
      :key="conv.id"
      type="button"
      class="list-group-item list-group-item-action d-flex align-items-center gap-2"
      :class="{ active: selectedId === conv.id }"
      @click="$emit('select', conv)"
    >
      <Avatar :src="conv.photo" :name="conv.username" :size="42" :squared="conv.isGroup" />
      <span class="flex-grow-1 min-w-0">
        <span class="d-flex justify-content-between align-items-center gap-2">
          <span class="fw-semibold text-truncate">{{ conv.username }}</span>
          <span class="d-flex align-items-center gap-1 flex-shrink-0">
            <span v-if="conv.isGroup" class="badge bg-secondary rounded-pill">Gruppo</span>
            <small v-if="conv.lastMessage" class="text-muted wt-when">
              {{ formatWhen(conv.lastMessage.dataSent) }}
            </small>
          </span>
        </span>
        <span v-if="conv.lastMessage" class="d-block small text-truncate opacity-75">
          <span v-if="conv.isGroup">{{ conv.lastMessage.senderUsername }}: </span>
          <svg
            v-if="conv.lastMessage.isPhoto"
            width="13"
            height="13"
            viewBox="0 0 16 16"
            fill="currentColor"
            class="align-text-bottom"
            aria-label="immagine"
            role="img"
          >
            <path d="M1.5 2h13A1.5 1.5 0 0 1 16 3.5v9A1.5 1.5 0 0 1 14.5 14h-13A1.5 1.5 0 0 1 0 12.5v-9A1.5 1.5 0 0 1 1.5 2zm12 1h-11a.5.5 0 0 0-.5.5v7l3.2-3.1a.5.5 0 0 1 .68 0L9 11l2.1-2a.5.5 0 0 1 .68 0L14 11V3.5a.5.5 0 0 0-.5-.5zM5 5.5a1.2 1.2 0 1 1-2.4 0 1.2 1.2 0 0 1 2.4 0z" />
          </svg>
          {{ conv.lastMessage.text }}
        </span>
      </span>
    </button>
  </div>
</template>

<script>
import Avatar from './Avatar.vue'

export default {
  name: 'ConversationList',
  components: { Avatar },
  props: {
    conversations: { type: Array, default: () => [] },
    selectedId: { type: String, default: null },
    emptyText: {
      type: String,
      default: 'Nessuna conversazione. Inizia una chat o crea un gruppo.'
    }
  },
  emits: ['select'],
  computed: {
    sorted() {
      return [...this.conversations].sort((a, b) => {
        const ta = a.lastMessage ? a.lastMessage.dataSent : ''
        const tb = b.lastMessage ? b.lastMessage.dataSent : ''
        return tb.localeCompare(ta)
      })
    }
  },
  methods: {
    formatWhen(value) {
      if (!value) return ''
      const d = new Date(value)
      if (isNaN(d.getTime())) return ''

      const startOfDay = (x) => new Date(x.getFullYear(), x.getMonth(), x.getDate())
      const days = Math.round((startOfDay(new Date()) - startOfDay(d)) / 86400000)

      if (days === 0) {
        return d.toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' })
      }
      if (days === 1) return 'Ieri'
      if (days < 7) return d.toLocaleDateString('it-IT', { weekday: 'short' })
      return d.toLocaleDateString('it-IT', {
        day: '2-digit',
        month: '2-digit',
        year: '2-digit'
      })
    }
  }
}
</script>

<style scoped>
.min-w-0 {
  min-width: 0;
}

.wt-when {
  white-space: nowrap;
  font-size: 0.75rem;
}
</style>

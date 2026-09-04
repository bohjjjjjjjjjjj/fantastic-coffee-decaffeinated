<template>
  <div class="list-group overflow-auto flex-grow-1">
    <p v-if="!conversations.length" class="text-muted small px-2 mt-2">
      Nessuna conversazione. Inizia una chat o crea un gruppo.
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
        <span class="d-flex justify-content-between align-items-start">
          <span class="fw-semibold text-truncate">{{ conv.username }}</span>
          <span v-if="conv.isGroup" class="badge bg-secondary rounded-pill ms-1">Gruppo</span>
        </span>
        <span v-if="conv.lastMessage" class="d-block small text-truncate opacity-75">
          <span v-if="conv.isGroup">{{ conv.lastMessage.senderUsername }}: </span>{{ conv.lastMessage.text }}
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
    selectedId: { type: String, default: null }
  },
  emits: ['select'],
  computed: {
    sorted() {
      // Il backend ordina già per ultimo messaggio; qui si mantiene stabile.
      return [...this.conversations].sort((a, b) => {
        const ta = a.lastMessage ? a.lastMessage.dataSent : ''
        const tb = b.lastMessage ? b.lastMessage.dataSent : ''
        return tb.localeCompare(ta)
      })
    }
  }
}
</script>

<style scoped>
.min-w-0 {
  min-width: 0;
}
</style>

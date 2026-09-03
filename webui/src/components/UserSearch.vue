<template>
  <div>
    <input
      v-model.trim="query"
      type="text"
      class="form-control mb-2"
      placeholder="Cerca utenti per username..."
      @input="debouncedSearch"
    >
    <ErrorMsg v-if="error" :msg="error" />
    <div class="list-group">
      <p v-if="!loading && searched && !results.length" class="text-muted small m-0">
        Nessun utente trovato.
      </p>
      <button
        v-for="u in results"
        :key="u.id"
        type="button"
        class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
        :class="{ active: isSelected(u) }"
        @click="pick(u)"
      >
        <span>{{ u.username }}</span>
        <span v-if="multiple && isSelected(u)">✓</span>
      </button>
    </div>
  </div>
</template>

<script>
import api from '../services/api.js'
import ErrorMsg from './ErrorMsg.vue'

export default {
  name: 'UserSearch',
  components: { ErrorMsg },
  props: {
    multiple: { type: Boolean, default: false },
    selected: { type: Array, default: () => [] }
  },
  emits: ['pick', 'update:selected'],
  data() {
    return {
      query: '',
      results: [],
      loading: false,
      searched: false,
      error: '',
      timer: null,
      pollTimer: null
    }
  },
  mounted() {
    this.search()
    // La lista contatti si aggiorna da sola: se un utente cambia username o
    // foto, la modifica compare senza dover riaprire o ridigitare la ricerca.
    this.pollTimer = setInterval(this.search, 4000)
  },
  beforeUnmount() {
    if (this.timer) clearTimeout(this.timer)
    if (this.pollTimer) clearInterval(this.pollTimer)
  },
  methods: {
    debouncedSearch() {
      if (this.timer) clearTimeout(this.timer)
      this.timer = setTimeout(this.search, 300)
    },
    async search() {
      this.loading = true
      this.error = ''
      try {
        const res = await api.searchUsers(this.query)
        this.results = res.data.users || []
        this.searched = true
      } catch (err) {
        this.error = err.response?.data?.message || 'Errore nella ricerca.'
      } finally {
        this.loading = false
      }
    },
    isSelected(u) {
      return this.selected.some((s) => s.id === u.id)
    },
    pick(u) {
      if (!this.multiple) {
        this.$emit('pick', u)
        return
      }
      const next = this.isSelected(u)
        ? this.selected.filter((s) => s.id !== u.id)
        : [...this.selected, u]
      this.$emit('update:selected', next)
    }
  }
}
</script>

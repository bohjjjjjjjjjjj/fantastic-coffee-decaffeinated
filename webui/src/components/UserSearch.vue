<template>
  <div>
    <input
      v-if="query === null"
      v-model.trim="internalQuery"
      type="text"
      class="form-control mb-2"
      placeholder="Cerca utenti per username..."
      @input="debouncedSearch"
    >
    <ErrorMsg v-if="error" :msg="error" />
    <div class="list-group">
      <p v-if="!loading && searched && !visibleResults.length" class="text-muted small m-0 px-3 py-2">
        Nessun utente trovato.
      </p>
      <button
        v-for="u in visibleResults"
        :key="u.id"
        type="button"
        class="list-group-item list-group-item-action d-flex align-items-center gap-2"
        :class="{ active: isSelected(u) }"
        @click="pick(u)"
      >
        <Avatar :src="u.photoUrl" :name="u.username" :size="34" />
        <span class="flex-grow-1 text-truncate">{{ u.username }}</span>
        <span v-if="multiple && isSelected(u)">✓</span>
      </button>
    </div>
  </div>
</template>

<script>
import api from '../services/api.js'
import ErrorMsg from './ErrorMsg.vue'
import Avatar from './Avatar.vue'

export default {
  name: 'UserSearch',
  components: { ErrorMsg, Avatar },
  props: {
    multiple: { type: Boolean, default: false },
    selected: { type: Array, default: () => [] },
    excludeIds: { type: Array, default: () => [] },
    query: { type: String, default: null }
  },
  emits: ['pick', 'update:selected'],
  data() {
    return {
      internalQuery: '',
      results: [],
      loading: false,
      searched: false,
      error: '',
      timer: null,
      pollTimer: null
    }
  },
  computed: {
    effectiveQuery() {
      return this.query === null ? this.internalQuery : this.query
    },
    visibleResults() {
      if (!this.excludeIds.length) return this.results
      return this.results.filter((u) => !this.excludeIds.includes(u.id))
    }
  },
  watch: {
    query() {
      this.debouncedSearch()
    }
  },
  mounted() {
    this.search()
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
      const asked = this.effectiveQuery
      this.loading = true
      this.error = ''
      try {
        const res = await api.searchUsers(asked)
        if (asked !== this.effectiveQuery) return
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

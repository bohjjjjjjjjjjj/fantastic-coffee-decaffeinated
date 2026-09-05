<template>
  <img
    v-if="src && !failed"
    :src="resolved"
    :alt="`Foto di ${name}`"
    class="wt-avatar"
    :class="squared ? 'rounded' : 'rounded-circle'"
    :style="boxStyle"
    @error="failed = true"
  >
  <span
    v-else
    class="wt-avatar d-inline-flex align-items-center justify-content-center text-white fw-semibold"
    :class="squared ? 'rounded' : 'rounded-circle'"
    :style="[boxStyle, { backgroundColor: color, fontSize: size * 0.42 + 'px' }]"
    :title="name"
    aria-hidden="true"
  >{{ initial }}</span>
</template>

<script>
import { mediaURL } from '../services/axios.js'

const PALETTE = [
  '#0d6efd', '#6610f2', '#6f42c1', '#d63384', '#dc3545',
  '#fd7e14', '#198754', '#20c997', '#0dcaf0', '#495057'
]

export default {
  name: 'Avatar',
  props: {
    src: { type: String, default: '' },
    name: { type: String, default: '' },
    size: { type: Number, default: 40 },
    squared: { type: Boolean, default: false }
  },
  data() {
    return { failed: false }
  },
  computed: {
    resolved() {
      return mediaURL(this.src)
    },
    initial() {
      return (this.name || '?').trim().charAt(0).toUpperCase() || '?'
    },
    color() {
      let h = 0
      for (const ch of this.name || '') h = (h * 31 + ch.charCodeAt(0)) >>> 0
      return PALETTE[h % PALETTE.length]
    },
    boxStyle() {
      return {
        width: this.size + 'px',
        height: this.size + 'px',
        flex: `0 0 ${this.size}px`
      }
    }
  },
  watch: {
    src() {
      this.failed = false
    }
  }
}
</script>

<style scoped>
.wt-avatar {
  object-fit: cover;
  user-select: none;
}
</style>

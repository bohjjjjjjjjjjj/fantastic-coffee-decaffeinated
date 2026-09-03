<script>
export default {
  name: 'LoadingSpinner',
  props: {
    loading: { type: Boolean, default: false }
  }
}
</script>

<template>
  <div v-if="loading">
    <div style="text-align: center">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
  </div>
  <div v-else>
    <slot />
  </div>
</template>

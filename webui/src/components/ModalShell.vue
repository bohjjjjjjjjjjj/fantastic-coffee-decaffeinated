<template>
  <div class="modal-backdrop-custom" @click.self="$emit('close')">
    <div class="modal-dialog-custom card shadow" :style="{ maxWidth: maxWidth + 'px' }">
      <div class="card-header d-flex justify-content-between align-items-center gap-2">
        <h5 class="m-0 text-truncate">{{ title }}</h5>
        <div class="d-flex align-items-center gap-2 flex-shrink-0">
          <slot name="actions" />
          <button type="button" class="btn-close" aria-label="Chiudi" @click="$emit('close')" />
        </div>
      </div>
      <div class="card-body">
        <slot />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ModalShell',
  props: {
    title: { type: String, default: '' },
    maxWidth: { type: Number, default: 480 }
  },
  emits: ['close']
}
</script>

<style scoped>
.modal-backdrop-custom {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1050;
  padding: 1rem;
}
.modal-dialog-custom {
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}
</style>

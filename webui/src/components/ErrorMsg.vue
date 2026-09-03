<script>
export default {
  name: 'ErrorMsg',
  props: {
    msg: { type: String, default: '' }
  }
}
</script>

<template>
  <div class="alert alert-danger" role="alert">
    {{ msg }}
  </div>
</template>

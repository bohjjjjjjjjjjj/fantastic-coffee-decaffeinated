<template>
  <div id="app">
    <RouterView />
  </div>
</template>

<script>
import { RouterView } from 'vue-router'

export default {
  name: 'App',
  components: {
    RouterView
  }
}
</script>

<style>
:root {
  --wt-bg: #DAE4F5;
}

html,
body {
  height: 100%;
  overflow: hidden;
  background-color: var(--wt-bg);
}

#app {
  height: 100vh;
  height: 100dvh;
  margin: 0;
  padding: 0;
  overflow: hidden;
  background-color: var(--wt-bg);
}
</style>

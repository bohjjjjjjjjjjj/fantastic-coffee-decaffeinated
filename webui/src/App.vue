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
html,
body {
  height: 100%;
  overflow: hidden;
}

#app {
  height: 100vh;
  /* Su mobile 100vh include la barra degli indirizzi, che compare e scompare:
     l'altezza "salta" e il contenuto viene tagliato. 100dvh segue lo spazio
     realmente visibile; i browser che non la conoscono usano la riga sopra. */
  height: 100dvh;
  margin: 0;
  padding: 0;
  overflow: hidden;
}
</style>

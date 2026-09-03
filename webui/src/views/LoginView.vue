<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div class="card p-4 shadow-sm" style="max-width: 400px; width: 100%;">
      <h3 class="card-title text-center mb-4">WASAText &mdash; Accedi</h3>
      <form @submit.prevent="handleLogin">
        <div class="mb-3">
          <label for="username" class="form-label">Nome utente</label>
          <input
            id="username"
            v-model.trim="username"
            type="text"
            class="form-control"
            placeholder="3-30 caratteri: lettere, numeri, _"
            minlength="3"
            maxlength="30"
            pattern="[a-zA-Z0-9_]+"
            required
          >
        </div>
        <ErrorMsg v-if="errorMessage" :msg="errorMessage" />
        <button type="submit" class="btn btn-primary w-100" :disabled="loading">
          {{ loading ? 'Accesso in corso...' : 'Accedi' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import api from '../services/api.js'
import ErrorMsg from '../components/ErrorMsg.vue'

export default {
  name: 'LoginView',
  components: { ErrorMsg },
  data() {
    return {
      username: '',
      errorMessage: '',
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      this.loading = true
      this.errorMessage = ''
      try {
        const response = await api.doLogin(this.username)
        const token = response.data.identifier
        if (!token) {
          throw new Error('Risposta di login non valida')
        }
        localStorage.setItem('token', token)
        localStorage.setItem('username', this.username)
        this.$router.push('/chat')
      } catch (err) {
        this.errorMessage =
          err.response?.data?.message || 'Errore durante l\'autenticazione.'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

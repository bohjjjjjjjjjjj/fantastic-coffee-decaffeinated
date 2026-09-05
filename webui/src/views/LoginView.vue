<template>
  <div class="wt-home d-flex flex-column">
    <h1 class="wt-title text-center">WASAText</h1>

    <div class="flex-grow-1 d-flex align-items-center justify-content-center px-3 pb-5">
      <div class="card p-4 shadow wt-login-card">
        <h3 class="card-title text-center mb-4">Log in</h3>
        <form @submit.prevent="handleLogin">
          <div class="mb-3">
            <label for="username" class="form-label">Nome utente</label>
            <input
              id="username"
              v-model.trim="username"
              type="text"
              class="form-control"
              minlength="3"
              maxlength="16"
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

<style scoped>
.wt-home {
  height: 100%;
  overflow: auto;
  background-color: var(--wt-bg);
  background-image: url('../assets/panda_2.jpg');
  background-repeat: repeat;
  background-size: 320px;
  background-position: center top;
}

.wt-title {
  margin: 0;
  padding: 1.5rem 1rem 1rem;
  color: #12355B;
  font-family: ui-monospace, "Cascadia Mono", "Segoe UI Mono", Consolas,
    "Courier New", monospace;
  font-weight: 700;
  font-size: clamp(2.6rem, 11vw, 4.75rem);
  letter-spacing: 0.06em;
  line-height: 1.1;
}

.wt-login-card {
  max-width: 400px;
  width: 100%;
  border: none;
  border-radius: 1rem;
  background-color: rgba(255, 255, 255, 0.96);
}
</style>

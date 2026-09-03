import axios from 'axios'

// __API_URL__ è iniettata da vite.config.js (usata anche in fase di valutazione).
// In dev/preview il fallback punta al backend locale.
const baseURL =
  typeof __API_URL__ !== 'undefined' ? __API_URL__ : 'http://localhost:3000'

const instance = axios.create({
  baseURL,
  timeout: 1000 * 10,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Inserisce il token Bearer nell'header Authorization se presente nel localStorage.
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Su 401 pulisce la sessione e torna al login.
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      if (window.location.pathname !== '/') {
        window.location.assign('/')
      }
    }
    return Promise.reject(error)
  }
)

export default instance

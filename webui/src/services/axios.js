import axios from 'axios'

// __API_URL__ è iniettata da vite.config.js (usata anche in fase di valutazione).
// In dev/preview il fallback punta al backend locale.
const configuredURL =
  typeof __API_URL__ !== 'undefined' ? __API_URL__ : 'http://localhost:3000'

// Quando la pagina viene aperta da un altro dispositivo della stessa rete
// (es. telefono su http://192.168.1.5:4173), "localhost" indicherebbe quel
// dispositivo e non la macchina che esegue il backend. In quel caso si riusa
// l'hostname da cui la pagina è stata servita, mantenendo porta e percorso
// dell'API configurata.
function resolveBaseURL() {
  const isLoopback = (h) =>
    h === 'localhost' || h === '127.0.0.1' || h === '::1' || h === '[::1]'
  try {
    const api = new URL(configuredURL, window.location.href)
    if (isLoopback(api.hostname) && !isLoopback(window.location.hostname)) {
      api.hostname = window.location.hostname
    }
    return api.origin + api.pathname.replace(/\/$/, '')
  } catch {
    return configuredURL
  }
}

const baseURL = resolveBaseURL()

// I file caricati sono referenziati con URL relativi (/media/...), serviti dal
// backend: vanno quindi risolti sulla base API, non sull'origine della pagina.
export function mediaURL(url) {
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  return baseURL + (url.startsWith('/') ? url : '/' + url)
}

export { baseURL }

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

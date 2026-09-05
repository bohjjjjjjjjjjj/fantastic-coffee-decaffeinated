import axios from 'axios'

const configuredURL =
  typeof __API_URL__ !== 'undefined' ? __API_URL__ : window.location.origin

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

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

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

import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/chat',
      name: 'chat',
      component: ChatView,
      meta: { requiresAuth: true }
    },
    {
      // Stessa view: la presenza di conversationId decide se mostrare la lista
      // o la conversazione a pagina intera. Così il tasto Indietro del browser
      // e del telefono funziona senza logica aggiuntiva.
      path: '/chat/:conversationId',
      name: 'conversation',
      component: ChatView,
      meta: { requiresAuth: true }
    }
  ]
})

// Controllo per impedire l'accesso alla chat senza login
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/')
  } else {
    next()
  }
})

export default router

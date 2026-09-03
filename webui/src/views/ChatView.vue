<template>
  <div class="container-fluid h-100">
    <div class="row h-100">
      <!-- Sidebar -->
      <div class="col-md-4 col-lg-3 border-end bg-light d-flex flex-column p-3 h-100">
        <div class="d-flex justify-content-between align-items-center mb-3">
          <h5 class="m-0">WASAText</h5>
          <button type="button" class="btn btn-outline-danger btn-sm" @click="logout">Esci</button>
        </div>

        <div class="mb-3 p-2 bg-white rounded border d-flex align-items-center gap-2">
          <img
            v-if="me.photoUrl"
            :src="me.photoUrl"
            alt="foto profilo"
            class="rounded-circle"
            style="width: 36px; height: 36px; object-fit: cover;"
          >
          <div class="flex-grow-1">
            <small class="text-muted d-block">Collegato come</small>
            <span class="fw-bold">{{ me.username }}</span>
          </div>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="openModal('profile')">
            Profilo
          </button>
        </div>

        <div class="d-grid gap-2 mb-3">
          <button type="button" class="btn btn-outline-primary btn-sm" @click="openModal('newChat')">
            + Nuova chat
          </button>
          <button type="button" class="btn btn-outline-primary btn-sm" @click="openModal('newGroup')">
            + Nuovo gruppo
          </button>
        </div>

        <ConversationList
          :conversations="conversations"
          :selected-id="selectedConv ? selectedConv.id : null"
          @select="selectConversation"
        />
      </div>

      <!-- Chat area -->
      <div class="col-md-8 col-lg-9 d-flex flex-column h-100 p-0">
        <template v-if="selectedConv">
          <div class="p-3 border-bottom bg-white d-flex justify-content-between align-items-center">
            <div>
              <h5 class="m-0">{{ selectedConv.username }}</h5>
              <small v-if="selectedConv.isGroup" class="text-muted">Gruppo</small>
            </div>
            <button
              v-if="selectedConv.isGroup"
              type="button"
              class="btn btn-outline-secondary btn-sm"
              @click="openGroupSettings"
            >
              Gestisci gruppo
            </button>
          </div>

          <ErrorMsg v-if="threadError" :msg="threadError" />

          <MessageThread
            :messages="messages"
            :my-username="me.username"
            @reply="startReply"
            @forward="startForward"
            @react="startReact"
            @delete="removeMessage"
            @toggle-reaction="toggleReaction"
          />

          <div class="p-3 bg-white border-top">
            <div v-if="replyTo" class="small text-muted mb-1 d-flex justify-content-between">
              <span>Rispondi a: {{ replyPreview }}</span>
              <button type="button" class="btn-close btn-sm" aria-label="annulla" @click="replyTo = null" />
            </div>
            <form class="d-flex gap-2" @submit.prevent="send">
              <input
                v-model="draft"
                type="text"
                class="form-control"
                placeholder="Scrivi un messaggio... (oppure incolla un URL immagine)"
                required
              >
              <button type="submit" class="btn btn-primary" :disabled="sending">Invia</button>
            </form>
          </div>
        </template>

        <div v-else class="h-100 d-flex justify-content-center align-items-center text-muted">
          Seleziona una conversazione per iniziare a chattare.
        </div>
      </div>
    </div>

    <!-- Profilo -->
    <ModalShell v-if="modal === 'profile'" title="Il mio profilo" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <form class="mb-3" @submit.prevent="updateUsername">
        <label class="form-label">Username</label>
        <div class="input-group">
          <input v-model.trim="form.username" type="text" class="form-control" pattern="[a-zA-Z0-9_]{3,30}" required>
          <button class="btn btn-primary" type="submit">Salva</button>
        </div>
      </form>
      <form @submit.prevent="updatePhoto">
        <label class="form-label">URL foto profilo</label>
        <div class="input-group">
          <input v-model.trim="form.photoUrl" type="url" class="form-control" placeholder="https://...">
          <button class="btn btn-primary" type="submit">Salva</button>
        </div>
      </form>
    </ModalShell>

    <!-- Nuova chat -->
    <ModalShell v-if="modal === 'newChat'" title="Nuova chat" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <UserSearch @pick="createDirectConversation" />
    </ModalShell>

    <!-- Nuovo gruppo -->
    <ModalShell v-if="modal === 'newGroup'" title="Nuovo gruppo" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <label class="form-label">Nome del gruppo</label>
      <input v-model.trim="form.groupName" type="text" class="form-control mb-3" pattern="[a-zA-Z0-9_ ]{1,50}" required>
      <label class="form-label">Membri ({{ form.members.length }})</label>
      <UserSearch multiple :selected="form.members" @update:selected="form.members = $event" />
      <button type="button" class="btn btn-primary w-100 mt-3" :disabled="!form.groupName" @click="createGroup">
        Crea gruppo
      </button>
    </ModalShell>

    <!-- Impostazioni gruppo -->
    <ModalShell v-if="modal === 'groupSettings'" title="Gestisci gruppo" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <form class="mb-3" @submit.prevent="renameGroup">
        <label class="form-label">Nome</label>
        <div class="input-group">
          <input v-model.trim="form.groupName" type="text" class="form-control" pattern="[a-zA-Z0-9_ ]{1,50}" required>
          <button class="btn btn-primary" type="submit">Salva</button>
        </div>
      </form>
      <form class="mb-3" @submit.prevent="updateGroupPhoto">
        <label class="form-label">URL foto gruppo</label>
        <div class="input-group">
          <input v-model.trim="form.groupPhoto" type="url" class="form-control" placeholder="https://...">
          <button class="btn btn-primary" type="submit">Salva</button>
        </div>
      </form>
      <div class="mb-3">
        <label class="form-label">Membri</label>
        <ul class="list-group mb-2">
          <li v-for="m in groupMembers" :key="m.id" class="list-group-item py-1">{{ m.username }}</li>
        </ul>
        <UserSearch multiple :selected="form.members" @update:selected="form.members = $event" />
        <button type="button" class="btn btn-outline-primary w-100 mt-2" :disabled="!form.members.length" @click="addMembers">
          Aggiungi selezionati
        </button>
      </div>
      <button type="button" class="btn btn-outline-danger w-100" @click="leaveGroup">Esci dal gruppo</button>
    </ModalShell>

    <!-- Inoltra -->
    <ModalShell v-if="modal === 'forward'" title="Inoltra messaggio" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <div class="list-group">
        <button
          v-for="c in conversations"
          :key="c.id"
          type="button"
          class="list-group-item list-group-item-action"
          @click="doForward(c)"
        >
          {{ c.username }} <span v-if="c.isGroup" class="badge bg-secondary">Gruppo</span>
        </button>
      </div>
    </ModalShell>

    <!-- Reagisci -->
    <ModalShell v-if="modal === 'react'" title="Reagisci" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <div class="d-flex flex-wrap gap-2">
        <button
          v-for="emoji in reactionChoices"
          :key="emoji"
          type="button"
          class="btn btn-outline-secondary fs-4"
          @click="doReact(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </ModalShell>
  </div>
</template>

<script>
import api from '../services/api.js'
import ConversationList from '../components/ConversationList.vue'
import MessageThread from '../components/MessageThread.vue'
import ModalShell from '../components/ModalShell.vue'
import UserSearch from '../components/UserSearch.vue'
import ErrorMsg from '../components/ErrorMsg.vue'

export default {
  name: 'ChatView',
  components: { ConversationList, MessageThread, ModalShell, UserSearch, ErrorMsg },
  data() {
    return {
      me: { username: localStorage.getItem('username') || '', photoUrl: '' },
      conversations: [],
      selectedConv: null,
      messages: [],
      draft: '',
      replyTo: null,
      sending: false,
      threadError: '',
      modal: null,
      modalError: '',
      groupMembers: [],
      activeMessage: null,
      pollTimer: null,
      reactionChoices: ['👍', '❤️', '😂', '😮', '😢', '🔥', '👏', '🎉'],
      form: {
        username: '',
        photoUrl: '',
        groupName: '',
        groupPhoto: '',
        members: []
      }
    }
  },
  computed: {
    replyPreview() {
      if (!this.replyTo) return ''
      const c = this.replyTo.content
      return c && c.photoUrl ? 'Foto' : (c ? c.text : '')
    }
  },
  async mounted() {
    await this.loadMe()
    await this.loadConversations()
    this.pollTimer = setInterval(() => {
      this.loadConversations()
      if (this.selectedConv) this.loadMessages(true)
    }, 4000)
  },
  beforeUnmount() {
    if (this.pollTimer) clearInterval(this.pollTimer)
  },
  methods: {
    async loadMe() {
      try {
        const res = await api.getMyUserInfo()
        this.me = res.data
        localStorage.setItem('username', this.me.username)
      } catch {
        /* interceptor gestisce il 401 */
      }
    },
    async loadConversations() {
      try {
        const res = await api.getMyConversations()
        this.conversations = res.data.conversations || []
        if (this.selectedConv) {
          const fresh = this.conversations.find((c) => c.id === this.selectedConv.id)
          if (fresh) this.selectedConv = fresh
        }
      } catch (err) {
        this.threadError = err.response?.data?.message || 'Errore nel caricamento delle conversazioni.'
      }
    },
    async selectConversation(conv) {
      this.selectedConv = conv
      this.replyTo = null
      this.messages = []
      await this.loadMessages()
    },
    async loadMessages(silent = false) {
      if (!this.selectedConv) return
      try {
        const res = await api.getConversation(this.selectedConv.id)
        this.messages = res.data.messages || []
        this.selectedConv = {
          ...this.selectedConv,
          username: res.data.username,
          isGroup: res.data.isGroup
        }
        this.threadError = ''
      } catch (err) {
        if (!silent) {
          this.threadError = err.response?.data?.message || 'Errore nel caricamento dei messaggi.'
        }
      }
    },
    async send() {
      if (!this.draft.trim() || !this.selectedConv) return
      this.sending = true
      try {
        const text = this.draft.trim()
        const content = /^https?:\/\/\S+$/i.test(text)
          ? { photoUrl: text }
          : { text }
        await api.sendMessage(
          this.selectedConv.id,
          content,
          this.replyTo ? this.replyTo.id : null
        )
        this.draft = ''
        this.replyTo = null
        await this.loadMessages()
        await this.loadConversations()
      } catch (err) {
        this.threadError = err.response?.data?.message || 'Errore durante l\'invio.'
      } finally {
        this.sending = false
      }
    },
    startReply(msg) {
      this.replyTo = msg
    },
    startForward(msg) {
      this.activeMessage = msg
      this.openModal('forward')
    },
    startReact(msg) {
      this.activeMessage = msg
      this.openModal('react')
    },
    async removeMessage(msg) {
      if (!window.confirm('Eliminare questo messaggio?')) return
      try {
        await api.deleteMessage(this.selectedConv.id, msg.id)
        await this.loadMessages()
        await this.loadConversations()
      } catch (err) {
        this.threadError = err.response?.data?.message || 'Impossibile eliminare il messaggio.'
      }
    },
    async toggleReaction({ msg, reaction }) {
      // Click su una reazione esistente: se è la mia, la rimuovo.
      if (reaction.userSenderId === this.me.id) {
        try {
          await api.uncommentMessage(this.selectedConv.id, msg.id, reaction.id)
          await this.loadMessages()
        } catch (err) {
          this.threadError = err.response?.data?.message || 'Errore.'
        }
      } else {
        this.startReact(msg)
      }
    },
    async doReact(emoji) {
      try {
        await api.commentMessage(this.selectedConv.id, this.activeMessage.id, emoji)
        this.closeModal()
        await this.loadMessages()
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore durante la reazione.'
      }
    },
    async doForward(conv) {
      try {
        await api.forwardMessage(this.selectedConv.id, this.activeMessage.id, conv.id)
        this.closeModal()
        await this.loadConversations()
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore durante l\'inoltro.'
      }
    },
    async createDirectConversation(user) {
      try {
        const res = await api.createConversation(user.username)
        this.closeModal()
        await this.loadConversations()
        const conv = this.conversations.find((c) => c.id === res.data.id) || {
          id: res.data.id,
          username: res.data.username,
          isGroup: false
        }
        await this.selectConversation(conv)
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Impossibile creare la chat.'
      }
    },
    async createGroup() {
      try {
        const res = await api.createGroup(
          this.form.groupName,
          this.form.members.map((m) => m.id)
        )
        this.closeModal()
        await this.loadConversations()
        const conv = this.conversations.find((c) => c.id === res.data.id)
        if (conv) await this.selectConversation(conv)
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Impossibile creare il gruppo.'
      }
    },
    async openGroupSettings() {
      this.form.groupName = this.selectedConv.username
      this.form.groupPhoto = ''
      this.form.members = []
      this.openModal('groupSettings')
      try {
        const res = await api.getGroupMembers(this.selectedConv.id)
        this.groupMembers = res.data.members || []
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nel caricamento dei membri.'
      }
    },
    async renameGroup() {
      try {
        await api.setGroupName(this.selectedConv.id, this.form.groupName)
        await this.loadConversations()
        await this.loadMessages(true)
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nel rinominare il gruppo.'
      }
    },
    async updateGroupPhoto() {
      try {
        await api.setGroupPhoto(this.selectedConv.id, this.form.groupPhoto)
        await this.loadConversations()
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nell\'aggiornare la foto.'
      }
    },
    async addMembers() {
      try {
        const res = await api.addToGroup(
          this.selectedConv.id,
          this.form.members.map((m) => m.id)
        )
        this.groupMembers = res.data.members || []
        this.form.members = []
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nell\'aggiungere i membri.'
      }
    },
    async leaveGroup() {
      if (!window.confirm('Uscire dal gruppo?')) return
      try {
        await api.leaveGroup(this.selectedConv.id)
        this.closeModal()
        this.selectedConv = null
        this.messages = []
        await this.loadConversations()
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nell\'uscita dal gruppo.'
      }
    },
    async updateUsername() {
      try {
        await api.setMyUserName(this.form.username)
        await this.loadMe()
        this.modalError = ''
      } catch (err) {
        this.modalError = err.response?.status === 409
          ? 'Username già in uso.'
          : (err.response?.data?.message || 'Errore nell\'aggiornamento.')
      }
    },
    async updatePhoto() {
      try {
        await api.setMyPhoto(this.form.photoUrl)
        await this.loadMe()
        this.modalError = ''
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nell\'aggiornamento.'
      }
    },
    openModal(name) {
      this.modalError = ''
      this.form.username = this.me.username
      this.form.photoUrl = this.me.photoUrl || ''
      if (name === 'newGroup') {
        this.form.groupName = ''
        this.form.members = []
      }
      this.modal = name
    },
    closeModal() {
      this.modal = null
      this.modalError = ''
      this.activeMessage = null
    },
    logout() {
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      this.$router.push('/')
    }
  }
}
</script>

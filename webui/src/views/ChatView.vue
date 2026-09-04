<template>
  <div class="wt-shell d-flex flex-column">
    <!-- ===== VISTA LISTA (route /chat) ===== -->
    <template v-if="!routeConvId">
      <div class="d-flex justify-content-between align-items-center p-3 border-bottom bg-white">
        <h5 class="m-0">WASAText</h5>
        <button type="button" class="btn btn-outline-danger btn-sm" @click="logout">Esci</button>
      </div>

      <div class="p-3 pb-0">
        <div class="p-2 bg-white rounded border d-flex align-items-center gap-2">
          <Avatar :src="me.photoUrl" :name="me.username" :size="38" />
          <div class="flex-grow-1 min-w-0">
            <small class="text-muted d-block">Collegato come</small>
            <span class="fw-bold text-truncate d-block">{{ me.username }}</span>
          </div>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="openModal('profile')">
            Profilo
          </button>
        </div>

        <div class="d-flex gap-2 mt-3">
          <button type="button" class="btn btn-outline-primary btn-sm flex-fill" @click="openModal('newChat')">
            + Nuova chat
          </button>
          <button type="button" class="btn btn-outline-primary btn-sm flex-fill" @click="openModal('newGroup')">
            + Nuovo gruppo
          </button>
        </div>
      </div>

      <!-- min-height:0 permette alla lista di scrollare invece di allungare la pagina -->
      <div class="flex-grow-1 d-flex flex-column p-3" style="min-height: 0;">
        <ConversationList
          :conversations="conversations"
          :selected-id="null"
          style="min-height: 0;"
          @select="selectConversation"
        />
      </div>
    </template>

    <!-- ===== VISTA CONVERSAZIONE (route /chat/:conversationId) ===== -->
    <template v-else>
      <div class="d-flex align-items-center gap-2 p-2 p-sm-3 border-bottom bg-white">
        <button
          type="button"
          class="btn btn-outline-secondary btn-sm d-flex align-items-center gap-1 flex-shrink-0"
          title="Torna alle conversazioni"
          @click="goBackToList"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
            <path d="M10.7 2.3a1 1 0 0 1 0 1.4L6.4 8l4.3 4.3a1 1 0 0 1-1.4 1.4l-5-5a1 1 0 0 1 0-1.4l5-5a1 1 0 0 1 1.4 0z" />
          </svg>
          <span class="d-none d-sm-inline">Indietro</span>
        </button>

        <Avatar
          v-if="selectedConv"
          :src="selectedConv.photo"
          :name="selectedConv.username"
          :size="40"
          :squared="selectedConv.isGroup"
        />
        <div class="flex-grow-1 min-w-0">
          <h6 class="m-0 text-truncate">{{ selectedConv && selectedConv.username ? selectedConv.username : 'Caricamento...' }}</h6>
          <small v-if="selectedConv && selectedConv.isGroup" class="text-muted">Gruppo</small>
        </div>

        <button
          v-if="selectedConv && selectedConv.isGroup"
          type="button"
          class="btn btn-outline-secondary btn-sm flex-shrink-0"
          @click="openGroupSettings"
        >
          <span class="d-none d-sm-inline">Gestisci gruppo</span>
          <span class="d-sm-none">Gruppo</span>
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

      <div class="p-2 p-sm-3 bg-white border-top">
        <div v-if="replyTo" class="small text-muted mb-1 d-flex justify-content-between">
          <span class="text-truncate">Rispondi a: {{ replyPreview }}</span>
          <button type="button" class="btn-close btn-sm flex-shrink-0" aria-label="annulla" @click="replyTo = null" />
        </div>
        <div v-if="attachmentPreview" class="mb-2 d-flex align-items-center gap-2">
          <img :src="attachmentPreview" alt="anteprima" class="rounded border" style="height: 56px;">
          <span class="small text-muted text-truncate">{{ attachmentName }}</span>
          <button type="button" class="btn-close btn-sm flex-shrink-0" aria-label="rimuovi immagine" @click="clearAttachment" />
        </div>
        <form class="d-flex gap-2" @submit.prevent="send">
          <input
            ref="fileInput"
            type="file"
            accept="image/png,image/jpeg,image/gif,image/webp"
            class="d-none"
            @change="onFileSelected"
          >
          <button
            type="button"
            class="btn btn-outline-secondary flex-shrink-0"
            title="Allega un'immagine"
            @click="$refs.fileInput.click()"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
              <path d="M1.5 2h13A1.5 1.5 0 0 1 16 3.5v9A1.5 1.5 0 0 1 14.5 14h-13A1.5 1.5 0 0 1 0 12.5v-9A1.5 1.5 0 0 1 1.5 2zm12 1h-11a.5.5 0 0 0-.5.5v7l3.2-3.1a.5.5 0 0 1 .68 0L9 11l2.1-2a.5.5 0 0 1 .68 0L14 11V3.5a.5.5 0 0 0-.5-.5zM5 5.5a1.2 1.2 0 1 1-2.4 0 1.2 1.2 0 0 1 2.4 0z" />
            </svg>
          </button>
          <input
            v-model="draft"
            type="text"
            class="form-control"
            placeholder="Scrivi un messaggio..."
          >
          <button
            type="submit"
            class="btn btn-primary flex-shrink-0"
            :disabled="sending || (!draft.trim() && !attachmentFile)"
          >
            {{ sending ? '...' : 'Invia' }}
          </button>
        </form>
      </div>
    </template>

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
      <label class="form-label">Foto profilo</label>
      <div class="d-flex align-items-center gap-2 mb-2">
        <Avatar :src="form.photoUrl" :name="me.username" :size="56" />
        <input
          type="file"
          class="form-control"
          accept="image/png,image/jpeg,image/gif,image/webp"
          :disabled="uploading"
          @change="onProfilePhotoSelected"
        >
      </div>
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
      <label class="form-label">Foto del gruppo</label>
      <div class="d-flex align-items-center gap-2 mb-3">
        <Avatar :src="form.groupPhoto" :name="form.groupName" :size="56" squared />
        <input
          type="file"
          class="form-control"
          accept="image/png,image/jpeg,image/gif,image/webp"
          :disabled="uploading"
          @change="onGroupPhotoSelected"
        >
      </div>
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
      <p v-if="conversations.length" class="form-label mb-1">Conversazioni esistenti</p>
      <div v-if="conversations.length" class="list-group mb-3">
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
      <!-- Permette di inoltrare anche a chi non ha ancora una chat con noi:
           la conversazione viene creata al momento. -->
      <p class="form-label mb-1">Oppure inoltra a un utente</p>
      <UserSearch @pick="forwardToUser" />
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
import Avatar from '../components/Avatar.vue'

export default {
  name: 'ChatView',
  components: { ConversationList, MessageThread, ModalShell, UserSearch, ErrorMsg, Avatar },
  data() {
    return {
      me: { username: localStorage.getItem('username') || '', photoUrl: '' },
      conversations: [],
      selectedConv: null,
      messages: [],
      draft: '',
      replyTo: null,
      sending: false,
      uploading: false,
      attachmentFile: null,
      attachmentPreview: '',
      attachmentName: '',
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
    },
    // La presenza del parametro nell'URL decide quale vista mostrare.
    routeConvId() {
      return this.$route.params.conversationId || null
    }
  },
  watch: {
    routeConvId() {
      this.syncWithRoute()
    }
  },
  async mounted() {
    await this.loadMe()
    await this.loadConversations()
    await this.syncWithRoute()
    this.pollTimer = setInterval(() => {
      // loadMe() serve a propagare anche le modifiche al proprio profilo
      // (username e foto) fatte da un'altra finestra o dispositivo.
      this.loadMe()
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
    // Selezionare una conversazione significa navigare: la vista dedicata
    // viene aperta dal watcher su routeConvId.
    selectConversation(conv) {
      if (this.routeConvId === conv.id) return
      this.$router.push({ name: 'conversation', params: { conversationId: conv.id } })
    },
    goBackToList() {
      this.$router.push({ name: 'chat' })
    },
    // Allinea lo stato della vista al parametro presente nell'URL.
    async syncWithRoute() {
      const id = this.routeConvId
      this.replyTo = null
      this.clearAttachment()
      this.threadError = ''
      this.messages = []
      if (!id) {
        this.selectedConv = null
        return
      }
      // Se la lista e' gia' carica si parte dai dati noti, altrimenti bastano
      // l'id e loadMessages a completare nome, foto e tipo.
      this.selectedConv = this.conversations.find((c) => c.id === id) || { id }
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
          photo: res.data.photo,
          isGroup: res.data.isGroup
        }
        this.threadError = ''
      } catch (err) {
        if (!silent) {
          this.threadError = err.response?.data?.message || 'Errore nel caricamento dei messaggi.'
        }
      }
    },
    // Carica il file e restituisce l'URL con cui referenziarlo.
    async uploadFile(file) {
      const res = await api.uploadMedia(file)
      return res.data.url
    },
    onFileSelected(event) {
      const file = event.target.files && event.target.files[0]
      if (!file) return
      this.attachmentFile = file
      this.attachmentName = file.name
      this.attachmentPreview = URL.createObjectURL(file)
    },
    clearAttachment() {
      if (this.attachmentPreview) URL.revokeObjectURL(this.attachmentPreview)
      this.attachmentFile = null
      this.attachmentPreview = ''
      this.attachmentName = ''
      if (this.$refs.fileInput) this.$refs.fileInput.value = ''
    },
    async send() {
      const text = this.draft.trim()
      if ((!text && !this.attachmentFile) || !this.selectedConv) return
      this.sending = true
      try {
        // Un messaggio puo' contenere testo, immagine o entrambi.
        const content = {}
        if (text) content.text = text
        if (this.attachmentFile) content.photoUrl = await this.uploadFile(this.attachmentFile)

        await api.sendMessage(
          this.selectedConv.id,
          content,
          this.replyTo ? this.replyTo.id : null
        )
        this.draft = ''
        this.replyTo = null
        this.clearAttachment()
        await this.loadMessages()
        await this.loadConversations()
      } catch (err) {
        this.threadError = err.response?.data?.message || 'Errore durante l' + String.fromCharCode(39) + 'invio.'
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
    // Inoltro verso un utente con cui non esiste ancora una conversazione:
    // createConversation la crea, o restituisce quella già esistente.
    async forwardToUser(user) {
      const sourceConvID = this.selectedConv.id
      const messageID = this.activeMessage.id
      try {
        const res = await api.createConversation(user.username)
        await api.forwardMessage(sourceConvID, messageID, res.data.id)
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
      this.form.groupPhoto = this.selectedConv.photo || ''
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
    async onGroupPhotoSelected(event) {
      const file = event.target.files && event.target.files[0]
      if (!file) return
      this.uploading = true
      this.modalError = ''
      try {
        const url = await this.uploadFile(file)
        await api.setGroupPhoto(this.selectedConv.id, url)
        this.form.groupPhoto = url
        await this.loadConversations()
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nel caricamento della foto.'
      } finally {
        this.uploading = false
        event.target.value = ''
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
        this.goBackToList()
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
    async onProfilePhotoSelected(event) {
      const file = event.target.files && event.target.files[0]
      if (!file) return
      this.uploading = true
      this.modalError = ''
      try {
        const url = await this.uploadFile(file)
        await api.setMyPhoto(url)
        this.form.photoUrl = url
        await this.loadMe()
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nel caricamento della foto.'
      } finally {
        this.uploading = false
        event.target.value = ''
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

<style scoped>
/* Riempie l'altezza resa disponibile da #app (100dvh) senza mai eccederla:
   min-height 0 lascia che i figli scrollabili si restringano davvero. */
.wt-shell {
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

/* Necessaria perché text-truncate funzioni dentro un contenitore flex. */
.min-w-0 {
  min-width: 0;
}
</style>

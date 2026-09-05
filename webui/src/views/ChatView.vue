<template>
  <div class="wt-shell d-flex flex-column">
    <template v-if="!routeConvId">
      <div class="d-flex justify-content-between align-items-center p-3 border-bottom bg-white gap-2">
        <h5 class="m-0">WASAText</h5>
        <button
          type="button"
          class="btn btn-light border d-flex align-items-center gap-2 py-1 px-2 min-w-0"
          title="Il mio profilo"
          @click="openModal('profile')"
        >
          <Avatar :src="me.photoUrl" :name="me.username" :size="32" />
          <span class="fw-semibold text-truncate">{{ me.username }}</span>
        </button>
      </div>

      <div class="p-3 pb-2">
        <div class="position-relative wt-search-wrap">
          <svg
            class="wt-search-icon"
            width="16"
            height="16"
            viewBox="0 0 16 16"
            fill="currentColor"
            aria-hidden="true"
          >
            <path d="M11.7 10.3a6 6 0 1 0-1.4 1.4l3 3a1 1 0 0 0 1.4-1.4l-3-3zM7 11a4 4 0 1 1 0-8 4 4 0 0 1 0 8z" />
          </svg>
          <input
            v-model.trim="search"
            type="text"
            class="form-control ps-5"
            placeholder="Cerca utenti per username..."
            @input="onSearchInput"
            @focus="onSearchInput"
          >
          <button
            v-if="search"
            type="button"
            class="btn-close wt-search-clear"
            aria-label="Annulla ricerca"
            @click="clearSearch"
          />

          <div v-if="showSearchDropdown" class="wt-search-dropdown shadow">
            <p class="small text-muted m-0 px-3 py-2 border-bottom">Utenti</p>
            <UserSearch :query="search" @pick="createDirectConversation" />
          </div>
        </div>
      </div>

      <div class="flex-grow-1 d-flex flex-column px-3 pb-3" style="min-height: 0;">
        <ConversationList
          :conversations="conversations"
          :selected-id="null"
          style="min-height: 0;"
          @select="selectConversation"
        />
      </div>

      <div v-if="showSearchDropdown" class="wt-search-backdrop" @click="searchDropdownOpen = false" />

      <div v-if="fabOpen" class="wt-fab-backdrop" @click="fabOpen = false" />
      <div class="wt-fab-wrap">
        <div v-if="fabOpen" class="wt-fab-menu list-group shadow">
          <button
            type="button"
            class="list-group-item list-group-item-action d-flex align-items-center gap-2"
            @click="openFromFab('newChat')"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
              <path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm6 5c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z" />
            </svg>
            Nuova chat
          </button>
          <button
            type="button"
            class="list-group-item list-group-item-action d-flex align-items-center gap-2"
            @click="openFromFab('newGroup')"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
              <path d="M7 14s-1 0-1-1 1-4 5-4 5 3 5 4-1 1-1 1H7zm4-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm-5.784 6A2.238 2.238 0 0 1 5 13c0-1.355.68-2.75 1.936-3.72A6.325 6.325 0 0 0 5 9c-4 0-5 3-5 4s1 1 1 1h4.216z" />
              <path d="M4.5 8a2.5 2.5 0 1 0 0-5 2.5 2.5 0 0 0 0 5z" />
            </svg>
            Nuovo gruppo
          </button>
        </div>
        <button
          type="button"
          class="btn btn-primary rounded-circle shadow wt-fab"
          :aria-expanded="fabOpen ? 'true' : 'false'"
          :title="fabOpen ? 'Chiudi' : 'Nuova conversazione'"
          aria-label="Nuova conversazione"
          @click="fabOpen = !fabOpen"
        >
          <svg width="24" height="24" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
            <path d="M2.678 11.894a1 1 0 0 1 .287.801 10.97 10.97 0 0 1-.398 2c1.395-.323 2.247-.697 2.634-.893a1 1 0 0 1 .71-.074A8.06 8.06 0 0 0 8 14c3.996 0 7-2.807 7-6 0-3.192-3.004-6-7-6S1 4.808 1 8c0 1.468.617 2.83 1.678 3.894zm-.493 3.905a21.682 21.682 0 0 1-.713.129c-.2.032-.352-.176-.273-.362a9.68 9.68 0 0 0 .244-.637l.003-.01c.248-.72.45-1.548.524-2.319C.743 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7-3.582 7-8 7a9.06 9.06 0 0 1-2.347-.306c-.52.263-1.639.742-3.468 1.105z" />
          </svg>
        </button>
      </div>
    </template>

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
        :is-group="!!(selectedConv && selectedConv.isGroup)"
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

    <ModalShell
      v-if="modal === 'profile'"
      :title="editingProfile ? 'Modifica profilo' : 'Il mio profilo'"
      :max-width="340"
      @close="closeModal"
    >
      <template #actions>
        <button
          v-if="!editingProfile"
          type="button"
          class="btn btn-sm btn-outline-secondary d-inline-flex align-items-center"
          title="Modifica profilo"
          aria-label="Modifica profilo"
          @click="startEditProfile"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
            <path d="M12.1 1.6a1.4 1.4 0 0 1 2 2l-.8.8-2-2 .8-.8zm-1.5 1.5 2 2-7.2 7.2a.5.5 0 0 1-.2.13l-2.6.87a.3.3 0 0 1-.38-.38l.87-2.6a.5.5 0 0 1 .12-.2l7.4-7.02z" />
          </svg>
        </button>
      </template>

      <ErrorMsg v-if="modalError" :msg="modalError" />

      <div v-if="!editingProfile" class="text-center py-2">
        <Avatar :src="me.photoUrl" :name="me.username" :size="150" />
        <h5 class="mt-3 mb-4 text-break">{{ me.username }}</h5>
        <button type="button" class="btn btn-outline-danger w-100" @click="logout">
          Esci
        </button>
      </div>

      <form v-else class="text-center py-2" @submit.prevent="saveProfile">
        <Avatar :src="form.photoUrl" :name="form.username" :size="150" />
        <div class="mt-3 text-start">
          <label class="form-label small">Immagine del profilo</label>
          <input
            type="file"
            class="form-control form-control-sm"
            accept="image/png,image/jpeg,image/gif,image/webp"
            :disabled="uploading || savingProfile"
            @change="onProfilePhotoSelected"
          >
          <label class="form-label small mt-3">Nome</label>
          <input
            v-model.trim="form.username"
            type="text"
            class="form-control"
            minlength="3"
            maxlength="16"
            required
          >
        </div>
        <div class="d-flex gap-2 mt-4">
          <button
            type="button"
            class="btn btn-outline-secondary flex-fill"
            :disabled="savingProfile"
            @click="cancelEditProfile"
          >
            Annulla
          </button>
          <button type="submit" class="btn btn-primary flex-fill" :disabled="uploading || savingProfile">
            {{ savingProfile ? 'Salvataggio...' : 'Salva' }}
          </button>
        </div>
      </form>
    </ModalShell>

    <ModalShell v-if="modal === 'newChat'" title="Nuova chat" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <UserSearch @pick="createDirectConversation" />
    </ModalShell>

    <ModalShell v-if="modal === 'newGroup'" title="Nuovo gruppo" @close="closeModal">
      <ErrorMsg v-if="modalError" :msg="modalError" />
      <label class="form-label">Nome del gruppo</label>
      <input v-model.trim="form.groupName" type="text" class="form-control mb-3" pattern="[a-zA-Z0-9_ ]{1,50}" required>
      <label class="form-label">Foto del gruppo</label>
      <div class="d-flex align-items-center gap-2 mb-3">
        <Avatar :src="form.groupPhoto" :name="form.groupName" :size="56" squared />
        <input
          type="file"
          class="form-control"
          accept="image/png,image/jpeg,image/gif,image/webp"
          :disabled="uploading"
          @change="onNewGroupPhotoSelected"
        >
      </div>
      <label class="form-label">Membri ({{ form.members.length }})</label>
      <UserSearch multiple :selected="form.members" @update:selected="form.members = $event" />
      <button
        type="button"
        class="btn btn-primary w-100 mt-3"
        :disabled="!form.groupName || uploading"
        @click="createGroup"
      >
        {{ uploading ? 'Caricamento foto...' : 'Crea gruppo' }}
      </button>
    </ModalShell>

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
      <div v-if="!addingMembers" class="mb-3">
        <div class="d-flex justify-content-between align-items-center mb-2">
          <label class="form-label m-0">Membri ({{ groupMembers.length }})</label>
          <button type="button" class="btn btn-sm btn-outline-primary d-inline-flex align-items-center gap-1" @click="startAddMembers">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
              <path d="M8 2.5a.75.75 0 0 1 .75.75v4h4a.75.75 0 0 1 0 1.5h-4v4a.75.75 0 0 1-1.5 0v-4h-4a.75.75 0 0 1 0-1.5h4v-4A.75.75 0 0 1 8 2.5z" />
            </svg>
            Aggiungi
          </button>
        </div>
        <ul class="list-group">
          <li
            v-for="m in groupMembers"
            :key="m.id"
            class="list-group-item py-1 d-flex align-items-center gap-2"
          >
            <Avatar :src="m.photoUrl" :name="m.username" :size="28" />
            <span class="text-truncate">{{ m.username }}</span>
          </li>
        </ul>
      </div>

      <div v-else class="mb-3">
        <div class="d-flex justify-content-between align-items-center mb-2">
          <label class="form-label m-0">Aggiungi membri</label>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="cancelAddMembers">
            Fatto
          </button>
        </div>
        <UserSearch
          multiple
          :selected="form.members"
          :exclude-ids="groupMemberIds"
          @update:selected="form.members = $event"
        />
        <button
          type="button"
          class="btn btn-outline-primary w-100 mt-2"
          :disabled="!form.members.length"
          @click="addMembers"
        >
          Aggiungi selezionati ({{ form.members.length }})
        </button>
      </div>

      <button type="button" class="btn btn-outline-danger w-100" @click="leaveGroup">Abbandona</button>
    </ModalShell>

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
      <p class="form-label mb-1">Oppure inoltra a un utente</p>
      <UserSearch @pick="forwardToUser" />
    </ModalShell>

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
      editingProfile: false,
      addingMembers: false,
      savingProfile: false,
      attachmentFile: null,
      attachmentPreview: '',
      attachmentName: '',
      threadError: '',
      modal: null,
      modalError: '',
      groupMembers: [],
      activeMessage: null,
      pollTimer: null,
      search: '',
      searchDropdownOpen: false,
      fabOpen: false,
      reactionChoices: ['👍', '❤️', '😂', '😮', '😢', '🔥', '👏', '🎉', '💩'],
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
    routeConvId() {
      return this.$route.params.conversationId || null
    },
    groupMemberIds() {
      return this.groupMembers.map((m) => m.id)
    },
    showSearchDropdown() {
      return this.searchDropdownOpen && !!this.search
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
      this.loadMe()
      this.loadConversations()
      if (this.selectedConv) this.loadMessages(true)
    }, 4000)
  },
  beforeUnmount() {
    if (this.pollTimer) clearInterval(this.pollTimer)
  },
  methods: {
    onSearchInput() {
      this.searchDropdownOpen = !!this.search
    },
    clearSearch() {
      this.search = ''
      this.searchDropdownOpen = false
    },
    openFromFab(name) {
      this.fabOpen = false
      this.openModal(name)
    },
    async loadMe() {
      try {
        const res = await api.getMyUserInfo()
        this.me = res.data
        localStorage.setItem('username', this.me.username)
      } catch {
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
    selectConversation(conv) {
      if (this.routeConvId === conv.id) return
      this.$router.push({ name: 'conversation', params: { conversationId: conv.id } })
    },
    goBackToList() {
      this.$router.push({ name: 'chat' })
    },
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
        this.clearSearch()
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
    async onNewGroupPhotoSelected(event) {
      const file = event.target.files && event.target.files[0]
      if (!file) return
      this.uploading = true
      this.modalError = ''
      try {
        this.form.groupPhoto = await this.uploadFile(file)
      } catch (err) {
        this.modalError = err.response?.data?.message || 'Errore nel caricamento della foto.'
      } finally {
        this.uploading = false
        event.target.value = ''
      }
    },
    async createGroup() {
      try {
        const res = await api.createGroup(
          this.form.groupName,
          this.form.members.map((m) => m.id)
        )
        if (this.form.groupPhoto) {
          await api.setGroupPhoto(res.data.id, this.form.groupPhoto)
        }
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
      this.addingMembers = false
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
    startAddMembers() {
      this.modalError = ''
      this.form.members = []
      this.addingMembers = true
    },
    cancelAddMembers() {
      this.form.members = []
      this.addingMembers = false
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
    startEditProfile() {
      this.modalError = ''
      this.form.username = this.me.username
      this.form.photoUrl = this.me.photoUrl || ''
      this.editingProfile = true
    },
    cancelEditProfile() {
      this.modalError = ''
      this.editingProfile = false
    },
    async saveProfile() {
      this.savingProfile = true
      this.modalError = ''
      try {
        if (this.form.username && this.form.username !== this.me.username) {
          await api.setMyUserName(this.form.username)
        }
        if ((this.form.photoUrl || '') !== (this.me.photoUrl || '')) {
          await api.setMyPhoto(this.form.photoUrl)
        }
        await this.loadMe()
        this.editingProfile = false
      } catch (err) {
        this.modalError = err.response?.status === 409
          ? 'Username già in uso.'
          : (err.response?.data?.message || "Errore nell'aggiornamento.")
      } finally {
        this.savingProfile = false
      }
    },
    async onProfilePhotoSelected(event) {
      const file = event.target.files && event.target.files[0]
      if (!file) return
      this.uploading = true
      this.modalError = ''
      try {
        this.form.photoUrl = await this.uploadFile(file)
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
        this.form.groupPhoto = ''
        this.form.members = []
      }
      this.modal = name
    },
    closeModal() {
      this.modal = null
      this.editingProfile = false
      this.addingMembers = false
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
.wt-shell {
  height: 100%;
  min-height: 0;
  overflow: hidden;
  position: relative;
}

.min-w-0 {
  min-width: 0;
}

.wt-search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #6c757d;
  pointer-events: none;
  z-index: 4;
}

.wt-search-clear {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  z-index: 4;
}

.wt-search-wrap {
  z-index: 1060;
}

.wt-search-dropdown {
  position: absolute;
  top: calc(100% + 0.35rem);
  left: 0;
  right: 0;
  z-index: 1060;
  max-height: 45vh;
  overflow-y: auto;
  border-radius: 0.75rem;
  background-color: #fff;
}

.wt-search-backdrop {
  position: absolute;
  inset: 0;
  z-index: 1055;
}

.wt-fab-wrap {
  position: absolute;
  right: 1.25rem;
  bottom: 1.25rem;
  z-index: 1050;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.5rem;
}

.wt-fab {
  width: 56px;
  height: 56px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.wt-fab-menu {
  min-width: 190px;
  border-radius: 0.75rem;
  overflow: hidden;
}

.wt-fab-backdrop {
  position: absolute;
  inset: 0;
  z-index: 1040;
}
</style>

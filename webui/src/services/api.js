import axios from './axios.js'

// Wrapper delle chiamate REST descritte in doc/api.yaml.
// Ogni funzione corrisponde a un operationId della spec.
export default {
  // Session
  doLogin: (name) => axios.post('/session', { name }),

  // Profilo
  getMyUserInfo: () => axios.get('/user'),
  setMyUserName: (username) => axios.put('/user/username', { username }),
  getMyPhoto: () => axios.get('/user/photo'),
  setMyPhoto: (photoUrl) => axios.put('/user/photo', { photoUrl }),
  searchUsers: (username) =>
    axios.get('/users', { params: username ? { username } : {} }),

  // Media: carica i byte grezzi del file, il tipo è riconosciuto dal backend
  uploadMedia: (file) =>
    axios.post('/media', file, {
      headers: { 'Content-Type': file.type || 'application/octet-stream' }
    }),

  // Conversazioni
  getMyConversations: () => axios.get('/conversations'),
  createConversation: (username) => axios.post('/conversations', { username }),
  getConversation: (conversationId) =>
    axios.get(`/conversations/${conversationId}`),

  // Messaggi
  sendMessage: (conversationId, content, replyToMessageId) =>
    axios.post(`/conversations/${conversationId}/messages`, {
      content,
      ...(replyToMessageId ? { replyToMessageId } : {})
    }),
  forwardMessage: (conversationId, messageId, targetConversationId) =>
    axios.post(
      `/conversations/${conversationId}/messages/${messageId}/forward`,
      { targetConversationId }
    ),
  deleteMessage: (conversationId, messageId) =>
    axios.delete(`/conversations/${conversationId}/messages/${messageId}`),

  // Reazioni
  commentMessage: (conversationId, messageId, reactionType) =>
    axios.post(
      `/conversations/${conversationId}/messages/${messageId}/reactions`,
      { reactionType }
    ),
  uncommentMessage: (conversationId, messageId, reactionId) =>
    axios.delete(
      `/conversations/${conversationId}/messages/${messageId}/reactions/${reactionId}`
    ),

  // Gruppi
  createGroup: (name, memberIds) =>
    axios.post('/groups', { name, ...(memberIds && memberIds.length ? { memberIds } : {}) }),
  getGroupDetails: (groupId) => axios.get(`/groups/${groupId}`),
  getGroupMembers: (groupId) => axios.get(`/groups/${groupId}/members`),
  addToGroup: (groupId, memberIds) =>
    axios.post(`/groups/${groupId}/members`, { memberIds }),
  setGroupName: (groupId, name) => axios.put(`/groups/${groupId}/name`, { name }),
  setGroupPhoto: (groupId, photoUrl) =>
    axios.put(`/groups/${groupId}/photo`, { photoUrl }),
  leaveGroup: (groupId) => axios.delete(`/groups/${groupId}/members/me`)
}

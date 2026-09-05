import axios from './axios.js'

export default {
  doLogin: (name) => axios.post('/session', { name }),

  getMyUserInfo: () => axios.get('/user'),
  setMyUserName: (username) => axios.put('/user/username', { username }),
  getMyPhoto: () => axios.get('/user/photo'),
  setMyPhoto: (photoUrl) => axios.put('/user/photo', { photoUrl }),
  searchUsers: (username) =>
    axios.get('/users', { params: username ? { username } : {} }),

  uploadMedia: (file) =>
    axios.post('/media', file, {
      headers: { 'Content-Type': file.type || 'application/octet-stream' }
    }),

  getMyConversations: () => axios.get('/conversations'),
  createConversation: (username) => axios.post('/conversations', { username }),
  getConversation: (conversationId) =>
    axios.get(`/conversations/${conversationId}`),

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

  commentMessage: (conversationId, messageId, reactionType) =>
    axios.post(
      `/conversations/${conversationId}/messages/${messageId}/reactions`,
      { reactionType }
    ),
  uncommentMessage: (conversationId, messageId, reactionId) =>
    axios.delete(
      `/conversations/${conversationId}/messages/${messageId}/reactions/${reactionId}`
    ),

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

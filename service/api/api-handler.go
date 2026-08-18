package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	// Session / Auth
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User
	rt.router.GET("/user", rt.wrap(rt.getMyUserInfo))
	rt.router.PUT("/user/username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users", rt.wrap(rt.searchUsers))

	// Conversations
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.POST("/conversations", rt.wrap(rt.createConversation))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation))

	// Messages
	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversations/:conversationId/messages/:messageId/forward", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.deleteMessage))

	// Reactions
	rt.router.POST("/conversations/:conversationId/messages/:messageId/reactions", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId/reactions/:reactionId", rt.wrap(rt.uncommentMessage))

	// Groups
	rt.router.POST("/groups", rt.wrap(rt.createGroup))
	rt.router.PUT("/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.DELETE("/groups/:groupId/members/me", rt.wrap(rt.leaveGroup))

	return rt.router
}

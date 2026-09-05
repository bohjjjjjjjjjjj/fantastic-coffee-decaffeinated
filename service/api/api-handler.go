package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.GET("/user", rt.wrap(rt.getMyUserInfo))
	rt.router.PUT("/user/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/user/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/user/photo", rt.wrap(rt.getMyPhoto))
	rt.router.GET("/users", rt.wrap(rt.searchUsers))

	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.POST("/conversations", rt.wrap(rt.createConversation))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation))

	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversations/:conversationId/messages/:messageId/forward", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.deleteMessage))

	rt.router.POST("/conversations/:conversationId/messages/:messageId/reactions", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId/reactions/:reactionId", rt.wrap(rt.uncommentMessage))

	rt.router.POST("/groups", rt.wrap(rt.createGroup))
	rt.router.GET("/groups/:groupId", rt.wrap(rt.getGroupDetails))
	rt.router.GET("/groups/:groupId/members", rt.wrap(rt.getGroupMembers))
	rt.router.POST("/groups/:groupId/members", rt.wrap(rt.addToGroup))
	rt.router.PUT("/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.DELETE("/groups/:groupId/members/me", rt.wrap(rt.leaveGroup))

	rt.router.POST("/media", rt.wrap(rt.uploadMedia))
	rt.router.GET("/media/:mediaId", rt.wrap(rt.getMedia))

	rt.router.GET("/liveness", rt.wrap(rt.liveness))

	return rt.router
}

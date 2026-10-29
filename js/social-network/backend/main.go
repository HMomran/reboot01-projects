package main

import (
	"log"
	"net/http"

	"social-network/pkg/auth"
	"social-network/pkg/db"
	"social-network/pkg/handlers"
	"social-network/pkg/middleware"
	ws "social-network/pkg/websocket"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	privateHub := ws.NewHub()
	groupHub := ws.NewGroupHub()
	feedHub := ws.NewFeedHub()

	// Wire the notification hub so createNotification pushes live events
	handlers.GlobalNotifHub = privateHub
	// Wire the feed hub so reactions/comments broadcast live updates
	handlers.GlobalFeedHub = feedHub
	// Wire notif hub into group chat so group messages trigger live notifications
	groupHub.NotifHub = privateHub

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Social Network API Running"))
	})

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	mux.HandleFunc("/register", auth.RegisterHandler(database))
	mux.HandleFunc("/login", auth.LoginHandler(database))
	mux.HandleFunc("/logout", auth.LogoutHandler(database))
	mux.HandleFunc("/me", handlers.MeHandler(database))
	mux.HandleFunc("/users", handlers.ListUsersHandler(database))
	mux.HandleFunc("/users/online", handlers.OnlineUsersHandler(database, privateHub))

	mux.HandleFunc("/follow", handlers.FollowUserHandler(database))
	mux.HandleFunc("/follow/accept", handlers.AcceptFollowRequestHandler(database))
	mux.HandleFunc("/follow/decline", handlers.DeclineFollowRequestHandler(database))
	mux.HandleFunc("/unfollow", handlers.UnfollowUserHandler(database))
	mux.HandleFunc("/follow/requests", handlers.ListPendingFollowRequestsHandler(database))
	mux.HandleFunc("/followers", handlers.ListFollowersHandler(database))
	mux.HandleFunc("/following", handlers.ListFollowingHandler(database))

	mux.HandleFunc("/profile", handlers.GetProfileHandler(database))
	mux.HandleFunc("/profile/posts", handlers.UserPostsHandler(database))
	mux.HandleFunc("/profile/privacy", handlers.UpdateProfilePrivacyHandler(database))
	mux.HandleFunc("/profile/password", handlers.ChangePasswordHandler(database))
	mux.HandleFunc("/profile/update", handlers.UpdateProfileHandler(database))

	mux.HandleFunc("/feed", handlers.FeedHandler(database))
	mux.HandleFunc("/explore", handlers.ExploreHandler(database))
	mux.HandleFunc("/posts", handlers.ListVisiblePostsHandler(database))
	mux.HandleFunc("/posts/create", handlers.CreatePostHandler(database))
	mux.HandleFunc("/posts/comments", handlers.ListPostCommentsHandler(database))
	mux.HandleFunc("/posts/comments/create", handlers.CreateCommentHandler(database))
	mux.HandleFunc("/posts/delete", handlers.DeletePostHandler(database))
	mux.HandleFunc("/posts/comments/delete", handlers.DeleteCommentHandler(database))
	mux.HandleFunc("/posts/react", handlers.ReactToPostHandler(database))
	mux.HandleFunc("/posts/comments/like", handlers.LikeCommentHandler(database))

	mux.HandleFunc("/notifications", handlers.ListNotificationsHandler(database))
	mux.HandleFunc("/notifications/read", handlers.MarkNotificationReadHandler(database))
	mux.HandleFunc("/notifications/read-all", handlers.MarkAllNotificationsReadHandler(database))
	mux.HandleFunc("/notifications/unread", handlers.UnreadNotificationsCountHandler(database))
	mux.HandleFunc("/notifications/clear-all", handlers.ClearAllNotificationsHandler(database))

	mux.HandleFunc("/groups", handlers.ListGroupsHandler(database))
	mux.HandleFunc("/groups/create", handlers.CreateGroupHandler(database))
	mux.HandleFunc("/groups/request", handlers.RequestToJoinGroupHandler(database))
	mux.HandleFunc("/groups/invite", handlers.InviteToGroupHandler(database))
	mux.HandleFunc("/groups/invites", handlers.ListMyGroupInvitesHandler(database))
	mux.HandleFunc("/groups/my", handlers.ListMyGroupsHandler(database))
	mux.HandleFunc("/groups/members", handlers.GroupMembersHandler(database))
	mux.HandleFunc("/groups/leave", handlers.LeaveGroupHandler(database))
	mux.HandleFunc("/groups/delete", handlers.DeleteGroupHandler(database))
	mux.HandleFunc("/groups/requests", handlers.ListGroupJoinRequestsHandler(database))
	mux.HandleFunc("/groups/invite/accept", handlers.AcceptGroupInviteHandler(database))
	mux.HandleFunc("/groups/invite/decline", handlers.DeclineGroupInviteHandler(database))
	mux.HandleFunc("/groups/request/accept", handlers.AcceptGroupJoinRequestHandler(database))
	mux.HandleFunc("/groups/request/decline", handlers.DeclineGroupJoinRequestHandler(database))

	mux.HandleFunc("/groups/posts", handlers.ListGroupPostsHandler(database))
	mux.HandleFunc("/groups/posts/create", handlers.CreateGroupPostHandler(database))
	mux.HandleFunc("/groups/posts/delete", handlers.DeleteGroupPostHandler(database))
	mux.HandleFunc("/groups/posts/comments", handlers.ListGroupCommentsHandler(database))
	mux.HandleFunc("/groups/posts/comments/create", handlers.CreateGroupCommentHandler(database))
	mux.HandleFunc("/groups/posts/comments/delete", handlers.DeleteGroupCommentHandler(database))

	mux.HandleFunc("/groups/events", handlers.ListGroupEventsHandler(database))
	mux.HandleFunc("/groups/events/create", handlers.CreateGroupEventHandler(database))
	mux.HandleFunc("/groups/events/respond", handlers.RespondToEventHandler(database))
	mux.HandleFunc("/groups/events/responses", handlers.ListEventResponsesHandler(database))

	mux.HandleFunc("/upload/avatar", handlers.UploadAvatarHandler(database))
	mux.HandleFunc("/upload/post-image", handlers.UploadPostImageHandler(database))
	mux.HandleFunc("/upload/comment-image", handlers.UploadCommentImageHandler(database))
	mux.HandleFunc("/upload/group-post-image", handlers.UploadGroupPostImageHandler(database))
	mux.HandleFunc("/upload/group-comment-image", handlers.UploadGroupCommentImageHandler(database))

	mux.HandleFunc("/ws/private", ws.PrivateChatHandler(database, privateHub))
	mux.HandleFunc("/ws/group", ws.GroupChatHandler(database, groupHub))
	mux.HandleFunc("/ws/presence", ws.PresenceHandler(database, privateHub))
	mux.HandleFunc("/ws/notifications", ws.NotificationsHandler(database, privateHub))
	mux.HandleFunc("/ws/feed", ws.FeedHandler(database, feedHub))

	mux.HandleFunc("/messages/private", handlers.ListPrivateMessagesHandler(database))
	mux.HandleFunc("/messages/group", handlers.ListGroupMessagesHandler(database))
	mux.HandleFunc("/messages/unread-count", handlers.GetUnreadCountHandler(database))
	mux.HandleFunc("/messages/mark-read", handlers.MarkReadHandler(database))
	mux.HandleFunc("/messages/recent", handlers.ListRecentConversationsHandler(database))

	log.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", middleware.CORSMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}

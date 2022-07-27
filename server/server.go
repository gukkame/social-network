package main

import (
	"fmt"
	"log"
	"net/http"
	db "real-time-forum/server/db"
	mid "real-time-forum/server/middleware"
	ath "real-time-forum/server/services/authentication"
	ds "real-time-forum/server/services/data"
	groups "real-time-forum/server/services/data/groups"
	dsm "real-time-forum/server/services/data/messages"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/users", mid.CORS(mid.Auth(ds.AllUsers)))
	//AUTHENTICATION
	http.HandleFunc("/signup", mid.CORS(ds.SignUp))
	http.HandleFunc("/login", mid.CORS(ds.LogIn))
	http.HandleFunc("/available/username", mid.CORS(ds.Username))
	http.HandleFunc("/available/email", mid.CORS(ds.Email))
	http.HandleFunc("/authedchat", mid.CORS(ath.AuthChat))
	//PROFILE
	http.HandleFunc("/profile", mid.CORS(ds.Profile))
	http.HandleFunc("/changeprofile", mid.CORS(ds.ChangeProfileStatus))
	http.HandleFunc("/activity", mid.CORS(ds.Activity))
	http.HandleFunc("/followers", mid.CORS(ds.Followers))

	//GROUPS
	http.HandleFunc("/groups", mid.CORS(ds.GetGroups))
	http.HandleFunc("/group", mid.CORS(groups.GroupInfoAndUserStatus))
	http.HandleFunc("/creategroup", mid.CORS(ds.NewGroup))
	// # Group Posts
	http.HandleFunc("/group/posts", mid.CORS(mid.Auth(groups.Posts)))
	http.HandleFunc("/group/post", mid.CORS(mid.Auth(groups.PostInfo)))
	http.HandleFunc("/group/post/new", mid.CORS(mid.Auth(groups.NewPost)))
	http.HandleFunc("/group/post/like", mid.CORS(mid.Auth(groups.LikePost)))
	http.HandleFunc("/group/post/dislike", mid.CORS(mid.Auth(groups.DislikePost)))
	// # # Group Post Comments
	http.HandleFunc("/group/post/comment/new", mid.CORS(mid.Auth(groups.NewComment)))
	http.HandleFunc("/group/post/comment/like", mid.CORS(mid.Auth(groups.LikeComment)))
	http.HandleFunc("/group/post/comment/dislike", mid.CORS(mid.Auth(groups.DislikeComment)))
	// # Group Events
	http.HandleFunc("/group/events", mid.CORS(mid.Auth(groups.Events)))
	http.HandleFunc("/group/event/new", mid.CORS(mid.Auth(groups.NewEvent)))
	http.HandleFunc("/group/event/going", mid.CORS(mid.Auth(groups.ReplyGoingEvent)))
	http.HandleFunc("/group/event/notgoing", mid.CORS(mid.Auth(groups.ReplyNotGoingEvent)))
	// # Invites to Group
	http.HandleFunc("/group/invite", mid.CORS(mid.Auth(groups.InviteToGroup)))
	http.HandleFunc("/group/invite/accept", mid.CORS(mid.Auth(groups.AcceptGroupInvite)))
	http.HandleFunc("/group/invite/deny", mid.CORS(mid.Auth(groups.DenyGroupInvite)))
	// # Join Request to Group
	http.HandleFunc("/group/join", mid.CORS(mid.Auth(groups.MakeJoinRequest)))
	http.HandleFunc("/group/join/cancel", mid.CORS(mid.Auth(groups.CancelJoinRequest)))
	http.HandleFunc("/group/leave", mid.CORS(mid.Auth(groups.LeaveGroup)))
	// # # Have to be owner to access these 3
	http.HandleFunc("/group/join/deny", mid.CORS(mid.Auth(groups.DenyJoinRequest)))
	http.HandleFunc("/group/join/accept", mid.CORS(mid.Auth(groups.AcceptJoinRequest)))
	http.HandleFunc("/group/join/requests", mid.CORS(mid.Auth(groups.ListJoinRequests)))

	//POSTS & COMMENTS
	http.HandleFunc("/createpost", mid.CORS(ds.CreatePost))
	http.HandleFunc("/editpost", mid.CORS(ds.EditPost))
	http.HandleFunc("/deletepost", mid.CORS(ds.DeletePost))
	http.HandleFunc("/onepost", mid.CORS(ds.GetOnePost))
	http.HandleFunc("/onecategory", mid.CORS(ds.GetOneCategory))
	http.HandleFunc("/allcategory", mid.CORS(ds.GetAllCategorys))
	http.HandleFunc("/createcomment", mid.CORS(ds.CreateComment))

	//LIKES
	http.HandleFunc("/likepost", mid.CORS(ds.LikePost))
	http.HandleFunc("/dislikepost", mid.CORS(ds.DislikePost))
	http.HandleFunc("/likecomment", mid.CORS(ds.LikeComment))
	http.HandleFunc("/dislikecomment", mid.CORS(ds.DislikeComment))

	//FOLLOWING
	http.HandleFunc("/follow", mid.CORS(ds.PerformFollow))
	http.HandleFunc("/checkfollow", mid.CORS(ds.CheckFollowRequest))
	http.HandleFunc("/cancelrequest", mid.CORS(ds.CancelFollowRequest))
	http.HandleFunc("/removefollower", mid.CORS(ds.RemoveFollower))
	http.HandleFunc("/acceptfollower", mid.CORS(ds.AcceptFollower))

	//CHAT SYSTEM
	http.HandleFunc("/WSconnect", dsm.InitiateChat)

	//OPENING AND CREATING DATABASE IF IT IS DELETED FOR SOME WIERD REASON (YOU SHOULD NEVER DELETE A DATABASE)
	db.Database()

	//CLOSING DATABASE CONNECTION WHEN MAIN FUNCTION GETS CLOSED AKA CTRL + C
	defer db.DBC.Close()

	//IMAGES -> ./resources
	fileServer := http.FileServer(http.Dir("./resources"))
	http.Handle("/resources/", http.StripPrefix("/resources", fileServer))

	//GOLANG SERVER
	fmt.Printf("API Server running at port http://localhost:8080/\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("(server.go) Golang server has stopped due to:")
		log.Fatal(err)
	}
}

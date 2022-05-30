package main

import (
	"fmt"
	"log"
	"net/http"
	db "real-time-forum/server/db"
	ds "real-time-forum/server/services/data"
	ath "real-time-forum/server/services/authentication"
	dsm "real-time-forum/server/services/data/messages"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	
	//GOLANG CRUD API HANDLERS (CRUD = CREATE/READ/UPDATE/DELETE)
	//AUTHENTICATION
	http.HandleFunc("/signup", ds.SignUp)
	http.HandleFunc("/login", ds.LogIn)
	http.HandleFunc("/available/username", ds.Username)
	http.HandleFunc("/available/email", ds.Email)
	http.HandleFunc("/authedchat", ath.AuthChat)
	
	//GENERAL DATA
	http.HandleFunc("/profile", ds.Profile)
	http.HandleFunc("/activity", ds.GetActivity)

	//POSTS & COMMENTS &
	http.HandleFunc("/createpost", ds.CreatePost)
	http.HandleFunc("/editpost", ds.EditPost)
	http.HandleFunc("/deletepost", ds.DeletePost)
	http.HandleFunc("/onepost", ds.GetOnePost)
	http.HandleFunc("/onecategory", ds.GetOneCategory)
	http.HandleFunc("/allcategory", ds.GetAllCategory)
	http.HandleFunc("/createcomment", ds.CreateComment)
	
	//LIKES
	http.HandleFunc("/likepost", ds.LikePost)
	http.HandleFunc("/dislikepost", ds.DislikePost)
	http.HandleFunc("/likecomment", ds.LikeComment)
	http.HandleFunc("/dislikecomment", ds.DislikeComment)
	
	//CHAT SYSTEM
	http.HandleFunc("/WSconnect", dsm.InitiateChat) 
		
	//OPENING AND CREATING DATABASE IF IT IS DELETED FOR SOME WIERD REASON (YOU SHOULD NEVER DELETE A DATABASE)
	db.Database()

	//CLOSING DATABASE CONNECTION WHEN MAIN FUNCTION GETS CLOSED AKA CTRL + C
	defer db.DBC.Close()

	//GOLANG SERVER
	fmt.Printf("API Server running at port http://localhost:8080/\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("(server.go) Golang server has stopped due to:")
		log.Fatal(err)
	}
}
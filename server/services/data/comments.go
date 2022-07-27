package data_services

import (
	"encoding/json"
	"fmt"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	"real-time-forum/server/services/data/groups"
	val "real-time-forum/server/services/validation"
	"strings"
)

type incDataS struct {
	Categoryname string
	Postid       string
	Content      string
	Image	string
}

type CommentsS struct {
	Id             int
	Content        string
	Created_at     string
	User_id        int
	Post_id        int
	Username       string
	Image	string
	Likes          []VotesS
	Category_title string
	User groups.User
	User_Image	string
}

var newComment incDataS

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)

	reqHeader := r.Header.Get("header1")
	splitToken := strings.Split(reqHeader, "Token=")
	token := strings.Join(splitToken, "")

	err := r.ParseMultipartForm(32 << 0) // maxMemory 32MB
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}
	
	newComment.Postid = r.Form["postid"][0]
	newComment.Content = r.Form["content"][0]
	newComment.Categoryname = r.Form["categoryname"][0]

	imagePath, err := groups.HandleImageUpload(r, "post/comment")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	newComment.Image = imagePath

	if !ath.AuthUser(token) {
		w.Write([]byte(`{"message": "User not authenticated"}`))
		return
	}

	if !val.ValidComment(newComment.Content) {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}

	row, err := db.DBC.Query("SELECT Users.id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
	if err != nil {
		w.Write([]byte(`{"message": "User not authenticated"}`))
		return
	}
	defer row.Close()

	var user ExistingId
	for row.Next() {
		err := row.Scan(
			&user.Id,
		)
		if err != nil {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
	}

	stmt, err := db.DBC.Prepare(`INSERT INTO Comments(content, created_at, user_id, post_id, image) VALUES(?, datetime("now"), ?, ?, ?)`)
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}

	stmt.Exec(newComment.Content, user.Id, newComment.Postid, newComment.Image)
	defer stmt.Close()

	rows2, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.Image, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC", newComment.Postid)
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}
	defer rows2.Close()

	var obj3 []CommentsS
	for rows2.Next() {
		var obj2 CommentsS
		err2 := rows2.Scan(&obj2.Id, &obj2.Content, &obj2.Image, &obj2.Created_at, &obj2.User_id, &obj2.Post_id, &obj2.Username)
		if err2 != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		rows4, err := db.DBC.Query("SELECT Comment_likes.type, Comment_likes.created_at, Comment_likes.user_id, Comment_likes.comment_id, Users.username FROM Comment_likes INNER JOIN Users ON Users.id = Comment_likes.user_id WHERE comment_id = ?", obj2.Id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows4.Close()

		var myCommentVotes []VotesS
		for rows4.Next() {
			var singleVotes VotesS
			err2 := rows4.Scan(&singleVotes.Type, &singleVotes.Created_at, &singleVotes.User_id, &singleVotes.Post_id, &singleVotes.Username)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
			myCommentVotes = append(myCommentVotes, singleVotes)
		}

		obj2.Likes = myCommentVotes

		obj3 = append(obj3, obj2)
	}

	var jsonData []byte
	jsonData, _ = json.Marshal((obj3))
	w.Write(jsonData)
}

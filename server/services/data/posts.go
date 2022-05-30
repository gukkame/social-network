package data_services

import (
	"encoding/json"
	"io"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	val "real-time-forum/server/services/validation"
	"strings"
)

type NewPostS struct {
	Title        string
	Description  string
	Categoryname string
}

type ExistingPostS struct {
	Title       string
	Description string
	Id          int
}

type ExistingId struct {
	Id int
}

type postResponseS struct {
	Message string
	Id      int64
}

type CookieS struct {
	Id       string
	Username string
}

type NeededData struct {
	Id            int
	Username      string
	Created_at    string
	Title         string
	Description   string
	CategoryTitle string
	Comments      []CommentsS
	Likes         []VotesS
	Cookie        CookieS
}

type OnePostS struct {
	Categoryname string
	Postid       string
}

var newPost NewPostS

func CreatePost(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &newPost)

		if val.ValidatePostData(newPost.Title, newPost.Description) {

			row, err := db.DBC.Query("SELECT id FROM Categories WHERE title = ?", newPost.Categoryname)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			defer row.Close()

			var existingCategory ExistingId
			for row.Next() {
				err := row.Scan(
					&existingCategory.Id,
				)
				if err != nil {
					w.Write([]byte(`{"message": "Malicious user detected"}`))
					return
				}
			}

			row1, err := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			defer row1.Close()

			var existingUser ExistingId
			for row1.Next() {
				err := row1.Scan(
					&existingUser.Id,
				)
				if err != nil {
					w.Write([]byte(`{"message": "Malicious user detected"}`))
					return
				}
			}

			stmt, err := db.DBC.Prepare(`INSERT INTO Posts(title, content, created_at, category_id, user_id) VALUES(?, ?, datetime("now"), ?, ?)`)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			defer stmt.Close()

			result, _ := stmt.Exec(newPost.Title, newPost.Description, existingCategory.Id, existingUser.Id)
			res, err := result.LastInsertId()
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}

			var response postResponseS
			response.Message = "Post created succesfully"
			response.Id = res
			var jsonData []byte
			jsonData, _ = json.Marshal((response))
			w.Write(jsonData)
		} else {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		var existingPost ExistingPostS
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &existingPost)

		if val.ValidatePostData(existingPost.Title, existingPost.Description) {

			preStmt, _ := db.DBC.Prepare("PRAGMA foreign_keys = ON")
			preStmt.Exec()
			preStmt.Close()

			stmt, err := db.DBC.Prepare(`UPDATE Posts SET title = ?, content = ? WHERE id = ?`)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}

			stmt.Exec(existingPost.Title, existingPost.Description, existingPost.Id)
			defer stmt.Close()

			w.Write([]byte(`{"message": "Post created succesfully"}`))
		} else {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		var existing ExistingId
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &existing)

		stmt, err := db.DBC.Prepare(`DELETE FROM Posts WHERE id = ?`)
		if err != nil {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}

		stmt.Exec(existing.Id)
		defer stmt.Close()

		w.Write([]byte(`{"message": "Post removed succesfully"}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

var onePost OnePostS

//SINGLE POST
func GetOnePost(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)

		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &onePost)

		rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE Posts.id = ? AND Categories.title = ?", onePost.Postid, onePost.Categoryname)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		var obj NeededData
		for rows.Next() {
			err := rows.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Username, &obj.CategoryTitle)
			if err != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
		}
		defer rows.Close()

		err = rows.Err()
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		rows2, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC", onePost.Postid)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows2.Close()

		var obj3 []CommentsS
		for rows2.Next() {
			var obj2 CommentsS
			err2 := rows2.Scan(&obj2.Id, &obj2.Content, &obj2.Created_at, &obj2.User_id, &obj2.Post_id, &obj2.Username)
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

			obj2.Likes = append(myCommentVotes)

			obj3 = append(obj3, obj2)

		}

		rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", onePost.Postid)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows3.Close()

		var myVotes []VotesS
		for rows3.Next() {
			var singleVotes VotesS
			err2 := rows3.Scan(&singleVotes.Type, &singleVotes.Created_at, &singleVotes.User_id, &singleVotes.Post_id, &singleVotes.Username)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
			myVotes = append(myVotes, singleVotes)
		}

		obj.Likes = append(myVotes)
		obj.Comments = append(obj3)

		if ath.AuthUser(token) {
			id, username := ath.CurrentCookieInfo(token)
			obj.Cookie.Id = id
			obj.Cookie.Username = username
		}

		var jsonData []byte
		jsonData, _ = json.Marshal((obj))
		w.Write(jsonData)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

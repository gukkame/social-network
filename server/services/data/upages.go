package data_services

import (
	"encoding/json"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	"strings"
)

type CompletePackageS3 struct {
	Package ActivityPackage
	Cookie  CookieS
}

type ActivityPackage struct {
	Posts      []NeededData
	VotedPosts []NeededData
	Comments   []CommentsS
}

func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

		row, err := db.DBC.Query("SELECT * FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		defer row.Close()

		var existingData NeverReleasedData
		var neededData ProfileS
		for row.Next() {
			err := row.Scan(
				&existingData.Id,
				&neededData.Username,
				&existingData.Password,
				&neededData.Email,
				&neededData.Age,
				&neededData.Gender,
				&neededData.FirstName,
				&neededData.LastName,
				&existingData.Date,
				&neededData.Token,
				&existingData.Expiry_date,
				&existingData.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		var jsonData []byte
		jsonData, _ = json.Marshal((neededData))
		w.Write(jsonData)

	}
}

func GetActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

		row1, err := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		defer row1.Close()

		var existingUser ExistingId
		for row1.Next() {
			err := row1.Scan(
				&existingUser.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		var aPackage ActivityPackage
		rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE user_id = ? ORDER BY Posts.created_at DESC", existingUser.Id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows.Close()

		existingData := []NeededData{}

		for rows.Next() {
			var obj NeededData
			err := rows.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Username, &obj.CategoryTitle)
			if err != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}

			rows2, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC", obj.Id)
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
				obj3 = append(obj3, obj2)
			}

			rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", obj.Id)
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
			existingData = append(existingData, obj)
			aPackage.Posts = append(existingData)
		}
		err = rows.Err()
		if err != nil {

			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		rows2, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Post_likes INNER JOIN Categories ON Categories.id = Posts.category_id INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Posts ON Posts.id = Post_likes.post_id WHERE Post_likes.user_id = ? ORDER BY Posts.created_at DESC", existingUser.Id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows2.Close()

		existingData2 := []NeededData{}

		for rows2.Next() {
			var obj NeededData
			err := rows2.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Username, &obj.CategoryTitle)
			if err != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}

			rows2, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC", obj.Id)
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
				obj3 = append(obj3, obj2)
			}

			rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", obj.Id)
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
			existingData2 = append(existingData2, obj)
			aPackage.VotedPosts = append(existingData2)
		}
		err = rows.Err()
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		rows3, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Comments.user_id INNER JOIN Categories ON Categories.id = Posts.category_id INNER JOIN Comments ON Comments.post_id = Posts.id WHERE Comments.user_id = ? ORDER BY Comments.created_at DESC", existingUser.Id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows3.Close()

		var obj3 []CommentsS
		for rows3.Next() {
			var obj2 CommentsS
			err2 := rows3.Scan(&obj2.Id, &obj2.Content, &obj2.Created_at, &obj2.User_id, &obj2.Post_id, &obj2.Username, &obj2.Category_title)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}

			rows4, err := db.DBC.Query("SELECT Comment_likes.type, Comment_likes.created_at, Comment_likes.user_id, Comment_likes.comment_id, Users.username FROM Comment_likes INNER JOIN Users ON Users.id = Comment_likes.user_id  WHERE comment_id = ?", obj2.Id)
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

		aPackage.Comments = append(obj3)

		var packages CompletePackageS3
		if ath.AuthUser(token) {
			id, username := ath.CurrentCookieInfo(token)
			packages.Cookie.Id = id
			packages.Cookie.Username = username
		}
		packages.Package = aPackage

		var jsonData []byte
		jsonData, _ = json.Marshal((packages))
		w.Write(jsonData)

	}
}

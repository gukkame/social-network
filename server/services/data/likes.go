package data_services

import (
	"encoding/json"
	"io"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type VotingS struct {
	Id   int
	Type string
}

type existingVoteS struct {
	Type       string
	Created_at string
	Post_id    int
	User_id    int
}

type VotesS struct {
	Type       string
	Created_at string
	Post_id    int
	User_id    int
	Username   string
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		var newLike VotingS
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &newLike)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
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
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		//IF LIKE EXISTS
		row2, err := db.DBC.Query("SELECT * FROM Post_likes WHERE post_id = ? AND user_id = ? AND type = ?", newLike.Id, user.Id, newLike.Type)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		defer row2.Close()

		var existingLike existingVoteS
		for row2.Next() {
			err := row2.Scan(
				&existingLike.Type,
				&existingLike.Created_at,
				&existingLike.Post_id,
				&existingLike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingLike.Post_id != 0 && existingLike.User_id != 0 {
			stmt, _ := db.DBC.Prepare(`DELETE FROM Post_likes WHERE post_id = ? AND user_id = ? AND type = ?`)
			stmt.Exec(newLike.Id, user.Id, newLike.Type)
			stmt.Close()
			rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", newLike.Id)
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

			if len(myVotes) == 0 {
				w.Write([]byte(`{"message": "No likes"}`))
				return
			}

			var jsonData []byte
			jsonData, _ = json.Marshal((myVotes))
			w.Write(jsonData)
			return
		}

		//IF DISLIKE EXISTS
		row3, err := db.DBC.Query("SELECT * FROM Post_likes WHERE type = ? AND post_id = ? AND user_id = ?", "dislike", newLike.Id, user.Id)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		defer row3.Close()

		var existingDislike existingVoteS
		for row3.Next() {
			err := row3.Scan(
				&existingDislike.Type,
				&existingDislike.Created_at,
				&existingDislike.Post_id,
				&existingDislike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingDislike.Post_id != 0 && existingDislike.User_id != 0 {
			stmt1, _ := db.DBC.Prepare(`DELETE FROM Post_likes WHERE post_id = ? AND user_id = ? AND type = ?`)
			stmt1.Exec(newLike.Id, user.Id, "dislike")
			defer stmt1.Close()
		}

		stmt2, err := db.DBC.Prepare(`INSERT INTO Post_likes(type, created_at, post_id, user_id) VALUES(?, datetime("now"), ?, ?)`)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		stmt2.Exec(newLike.Type, newLike.Id, user.Id)
		defer stmt2.Close()

		rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", newLike.Id)
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

		var jsonData []byte
		jsonData, _ = json.Marshal((myVotes))
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		var newDislike VotingS
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &newDislike)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
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
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		//IF DISLIKE EXISTS
		row3, err := db.DBC.Query("SELECT * FROM Post_likes WHERE type = ? AND post_id = ? AND user_id = ?", newDislike.Type, newDislike.Id, user.Id)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		defer row3.Close()

		var existingDislike existingVoteS
		for row3.Next() {
			err := row3.Scan(
				&existingDislike.Type,
				&existingDislike.Created_at,
				&existingDislike.Post_id,
				&existingDislike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingDislike.Post_id != 0 && existingDislike.User_id != 0 {
			stmt1, _ := db.DBC.Prepare(`DELETE FROM Post_likes WHERE post_id = ? AND user_id = ? AND type = ?`)
			stmt1.Exec(newDislike.Id, user.Id, newDislike.Type)
			defer stmt1.Close()

			rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", newDislike.Id)
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

			if len(myVotes) == 0 {
				w.Write([]byte(`{"message": "No likes"}`))
				return
			}

			var jsonData []byte
			jsonData, _ = json.Marshal((myVotes))
			w.Write(jsonData)
			return
		}

		//IF LIKE EXISTS
		row2, err := db.DBC.Query("SELECT * FROM Post_likes WHERE post_id = ? AND user_id = ? AND type = ?", newDislike.Id, user.Id, "like")
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		defer row2.Close()

		var existingLike existingVoteS
		for row2.Next() {
			err := row2.Scan(
				&existingLike.Type,
				&existingLike.Created_at,
				&existingLike.Post_id,
				&existingLike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingLike.Post_id != 0 && existingLike.User_id != 0 {
			stmt, _ := db.DBC.Prepare(`DELETE FROM Post_likes WHERE post_id = ? AND user_id = ? AND type = ?`)
			stmt.Exec(newDislike.Id, user.Id, "like")
			defer stmt.Close()
		}

		stmt2, err := db.DBC.Prepare(`INSERT INTO Post_likes(type, created_at, post_id, user_id) VALUES(?, datetime("now"), ?, ?)`)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		stmt2.Exec(newDislike.Type, newDislike.Id, user.Id)
		defer stmt2.Close()

		rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", newDislike.Id)
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

		var jsonData []byte
		jsonData, _ = json.Marshal((myVotes))
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		var newLike VotingS
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &newLike)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
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
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		//IF LIKE EXISTS
		row2, err := db.DBC.Query("SELECT * FROM Comment_likes WHERE comment_id = ? AND user_id = ? AND type = ?", newLike.Id, user.Id, newLike.Type)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		defer row2.Close()

		var existingLike existingVoteS
		for row2.Next() {
			err := row2.Scan(
				&existingLike.Type,
				&existingLike.Created_at,
				&existingLike.Post_id,
				&existingLike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingLike.Post_id != 0 && existingLike.User_id != 0 {
			stmt, _ := db.DBC.Prepare(`DELETE FROM Comment_likes WHERE comment_id = ? AND user_id = ? AND type = ?`)
			stmt.Exec(newLike.Id, user.Id, newLike.Type)
			defer stmt.Close()
			w.Write([]byte(`{"message": "Comment like removed"}`))
			return
		}

		//IF DISLIKE EXISTS
		row3, err := db.DBC.Query("SELECT * FROM Comment_likes WHERE type = ? AND comment_id = ? AND user_id = ?", "dislike", newLike.Id, user.Id)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		defer row3.Close()

		var existingDislike existingVoteS
		for row3.Next() {
			err := row3.Scan(
				&existingDislike.Type,
				&existingDislike.Created_at,
				&existingDislike.Post_id,
				&existingDislike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingDislike.Post_id != 0 && existingDislike.User_id != 0 {
			stmt1, _ := db.DBC.Prepare(`DELETE FROM Comment_likes WHERE comment_id = ? AND user_id = ? AND type = ?`)
			stmt1.Exec(newLike.Id, user.Id, "dislike")
			defer stmt1.Close()
		}

		stmt2, err := db.DBC.Prepare(`INSERT INTO Comment_likes(type, created_at, comment_id, user_id) VALUES(?, datetime("now"), ?, ?)`)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		stmt2.Exec(newLike.Type, newLike.Id, user.Id)
		defer stmt2.Close()
		w.Write([]byte(`{"message": "Comment like inserted"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		var newDislike VotingS
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		json.Unmarshal([]byte(reqBody), &newDislike)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
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
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		//IF DISLIKE EXISTS
		row3, err := db.DBC.Query("SELECT * FROM Comment_likes WHERE type = ? AND comment_id = ? AND user_id = ?", newDislike.Type, newDislike.Id, user.Id)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		defer row3.Close()

		var existingDislike existingVoteS
		for row3.Next() {
			err := row3.Scan(
				&existingDislike.Type,
				&existingDislike.Created_at,
				&existingDislike.Post_id,
				&existingDislike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingDislike.Post_id != 0 && existingDislike.User_id != 0 {
			stmt1, _ := db.DBC.Prepare(`DELETE FROM Comment_likes WHERE comment_id = ? AND user_id = ? AND type = ?`)
			stmt1.Exec(newDislike.Id, user.Id, newDislike.Type)
			defer stmt1.Close()
			w.Write([]byte(`{"message": "Comment dislike removed"}`))
			return
		}

		//IF LIKE EXISTS
		row2, err := db.DBC.Query("SELECT * FROM Comment_likes WHERE comment_id = ? AND user_id = ? AND type = ?", newDislike.Id, user.Id, "like")
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		defer row2.Close()

		var existingLike existingVoteS
		for row2.Next() {
			err := row2.Scan(
				&existingLike.Type,
				&existingLike.Created_at,
				&existingLike.Post_id,
				&existingLike.User_id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}
		}

		if existingLike.Post_id != 0 && existingLike.User_id != 0 {
			stmt, _ := db.DBC.Prepare(`DELETE FROM Comment_likes WHERE comment_id = ? AND user_id = ? AND type = ?`)
			stmt.Exec(newDislike.Id, user.Id, "like")
			defer stmt.Close()
		}

		stmt2, err := db.DBC.Prepare(`INSERT INTO Comment_likes(type, created_at, comment_id, user_id) VALUES(?, datetime("now"), ?, ?)`)
		if err != nil {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		stmt2.Exec(newDislike.Type, newDislike.Id, user.Id)
		defer stmt2.Close()

		w.Write([]byte(`{"message": "Comment dislike inserted"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

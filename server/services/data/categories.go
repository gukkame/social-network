package data_services

import (
	"encoding/json"
	"io"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	"strings"
)

type CategoryS struct {
	Categoryname string
}

type CategoriesPosts struct {
	Go         []NeededData
	Html       []NeededData
	Css        []NeededData
	JavaScript []NeededData
	Vuejs      []NeededData
}

type CompletePackageS1 struct {
	Posts  []NeededData
	Cookie CookieS
}

type CompletePackageS2 struct {
	Posts  CategoriesPosts
	Cookie CookieS
}

var CateName CategoryS

//GETTING ALL EXISTING POSTS FOR 1 CATEGORY FE: Go
func GetOneCategory(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")
		json.Unmarshal([]byte(reqBody), &CateName)

		row, err := db.DBC.Query("SELECT id FROM Categories WHERE title = ?", CateName.Categoryname)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		var existingCategory ExistingId
		for row.Next() {
			err := row.Scan(
				&existingCategory.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
		}
		row.Close()

		rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE category_id = ? ORDER BY Posts.created_at DESC", existingCategory.Id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		var existingData []NeededData

		var obj NeededData

		for rows.Next() {
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
		}

		rows.Close()
		err = rows.Err()
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		var packages CompletePackageS1
		if ath.AuthUser(token) {
			id, username := ath.CurrentCookieInfo(token)
			packages.Cookie.Id = id
			packages.Cookie.Username = username
		}
		packages.Posts = append(existingData)

		var jsonData []byte
		jsonData, _ = json.Marshal((packages))
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

//HOME PAGE LATEST POSTS 3 LATEST FROM EACH CATEGORY
func GetAllCategory(w http.ResponseWriter, r *http.Request) {
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

		var allcategory CategoriesPosts
		existingData := []NeededData{}
		for i := 0; i <= 5; i++ {
			rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE category_id = ? ORDER BY Posts.created_at DESC LIMIT 3", i)
			if err != nil {
				w.Write([]byte(`{"message": "Get request failed"}`))
				return
			}
			defer rows.Close()

			for rows.Next() {
				var obj NeededData
				err := rows.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Username, &obj.CategoryTitle)
				if err != nil {
					w.Write([]byte(`{"message": "Get request failed"}`))
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
			}

			switch i {
			case 1:
				allcategory.Go = append(existingData)
			case 2:
				allcategory.Html = append(existingData)
			case 3:
				allcategory.Css = append(existingData)
			case 4:
				allcategory.JavaScript = append(existingData)
			case 5:
				allcategory.Vuejs = append(existingData)
			}
			existingData = []NeededData{}

			err = rows.Err()
			if err != nil {
				w.Write([]byte(`{"message": "Get request failed"}`))
				return
			}
		}

		var packages CompletePackageS2
		if ath.AuthUser(token) {
			id, username := ath.CurrentCookieInfo(token)
			packages.Cookie.Id = id
			packages.Cookie.Username = username
		}
		packages.Posts = allcategory

		var jsonData []byte
		jsonData, _ = json.Marshal((packages))
		w.Write(jsonData)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

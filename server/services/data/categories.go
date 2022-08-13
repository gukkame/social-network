package data_services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	db "real-time-forum/server/db"
	mw "real-time-forum/server/middleware"
	"real-time-forum/server/services/data/groups"
)

type CategoryS struct {
	Categoryname string
}

type CategoriesPosts struct {
	Go         []NeededData `json:"Go"`
	Html       []NeededData `json:"Html"`
	Css        []NeededData `json:"Css"`
	JavaScript []NeededData `json:"JavaScript"`
	Vuejs      []NeededData `json:"Vuejs"`
}

type CompletePackageS1 struct {
	Posts  []NeededData
	Cookie CookieS
}

type CompletePackageS2 struct {
	Posts  CategoriesPosts
	Cookie CookieS
}

var Categories = []string{"Go", "HTML5", "CSS", "JavaScript", "Vue.js"}

var CateName CategoryS

func GetOneCategory(w http.ResponseWriter, r *http.Request) {
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	userId := groups.UserId(w, r)

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

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

		rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Posts.privacy, Users.username, Users.id, Users.avatar_image, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE category_id = ? ORDER BY Posts.created_at DESC", existingCategory.Id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		var existingData []NeededData
		var obj NeededData

		for rows.Next() {
			err := rows.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Privacy, &obj.Username, &obj.User_id, &obj.User_image, &obj.CategoryTitle)
			if err != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}

			if obj.Privacy == "Private" {
				var matches int
				query := "SELECT COUNT(user_id) FROM Posts_allowed_users WHERE post_id = ? AND user_id = ?"
				row := db.DBC.QueryRow(query, obj.Id, userId)
				err := row.Scan(&matches)
				if err != nil {
					w.Write([]byte(`{"message": "Post request failed"}`))
					return
				}

				if matches == 0 {
					continue
				}
			}
			// if privacy "followers" and request maker is not post owner
			if obj.Privacy == "Followers" && userId != obj.User_id {
				var matches int
				query := "SELECT COUNT(follower_id) FROM Followers WHERE status = ? AND follower_id = ? AND recipient_id = ?"
				row := db.DBC.QueryRow(query, "following", userId, obj.User_id)
				err := row.Scan(&matches)
				if err != nil {
					w.Write([]byte(`{"message": "Post request failed"}`))
					return
				}

				if matches == 0 {
					continue
				}
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

			obj.Likes = myVotes
			obj.Comments = obj3

			existingData = append(existingData, obj)
		}

		rows.Close()
		err = rows.Err()
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		var jsonData []byte
		jsonData, _ = json.Marshal((existingData))
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

//HOME PAGE LATEST POSTS 3 LATEST FROM EACH CATEGORY
func GetAllCategorys(w http.ResponseWriter, r *http.Request) {
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	userId := groups.UserId(w, r)

	switch r.Method {
	case "POST":

		var allcategory CategoriesPosts
		existingData := []NeededData{}
		for i := 0; i <= 5; i++ {
			rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Posts.image, Posts.privacy, Users.username, Users.id, Users.avatar_image, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE category_id = ? ORDER BY Posts.created_at DESC LIMIT 3", i)
			if err != nil {
				w.Write([]byte(`{"message": "Get request failed"}`))
				return
			}
			defer rows.Close()

			for rows.Next() {
				var obj NeededData
				err := rows.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Image, &obj.Privacy, &obj.Username, &obj.User_id, &obj.User_image, &obj.CategoryTitle)
				if err != nil {
					w.Write([]byte(`{"message": "Get request failed"}`))
					return
				}
				if obj.Privacy == "Private" {
					var matches int
					query := "SELECT COUNT(user_id) FROM Posts_allowed_users WHERE post_id = ? AND user_id = ?"
					row := db.DBC.QueryRow(query, obj.Id, userId)
					err := row.Scan(&matches)
					if err != nil {
						w.Write([]byte(`{"message": "Post request failed"}`))
						return
					}

					if matches == 0 {
						continue
					}
				}
				// if privacy "followers" and request maker is not post owner
				if obj.Privacy == "Followers" && userId != obj.User_id {
					var matches int
					query := "SELECT COUNT(follower_id) FROM Followers WHERE status = ? AND follower_id = ? AND recipient_id = ?"
					row := db.DBC.QueryRow(query, "following", userId, obj.User_id)
					err := row.Scan(&matches)
					if err != nil {
						w.Write([]byte(`{"message": "Post request failed"}`))
						return
					}

					if matches == 0 {
						continue
					}
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
				obj.Likes = myVotes
				obj.Comments = obj3
				existingData = append(existingData, obj)
			}

			switch i {
			case 1:
				allcategory.Go = existingData
			case 2:
				allcategory.Html = existingData
			case 3:
				allcategory.Css = existingData
			case 4:
				allcategory.JavaScript = existingData
			case 5:
				allcategory.Vuejs = existingData
			}
			existingData = []NeededData{}

			err = rows.Err()
			if err != nil {
				w.Write([]byte(`{"message": "Get request failed"}`))
				return
			}
		}

		var jsonData []byte
		jsonData, _ = json.Marshal((allcategory))
		w.Write(jsonData)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func (category ExistingId) DbGetAllPosts(limit int) (posts Posts, err error) {
	query := fmt.Sprintf("SELECT Posts.user_id, Posts.id, Posts.title, Posts.content, Posts.Image, Posts.created_at, Users.username, Users.avatar_image, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE category_id = ? ORDER BY Posts.created_at DESC LIMIT %d", limit)
	rows, err := db.DBC.Query(query, category.Id)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		post := NeededData{}
		err := rows.Scan(&post.User_id, &post.Id, &post.Title, &post.Description, &post.Image, &post.Created_at, &post.Username, &post.User_image, &post.CategoryTitle)
		if err != nil {
			fmt.Println(err)
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

type Posts []NeededData

func (posts Posts) FilterByPrivacy(userId int) (allowedPosts []NeededData, err error) {
	for _, post := range posts {
		if post.Privacy == "Followers" {
			isFollower, err := CheckFollowing(userId, post.User_id)
			if err != nil {
				return allowedPosts, err
			}
			if isFollower {
				allowedPosts = append(allowedPosts, post)
			}
		}
		if post.Privacy == "Private" {
			isAllowed, err := post.IsUserAllowed(userId)
			if err != nil {
				return allowedPosts, err
			}
			if isAllowed {
				allowedPosts = append(allowedPosts, post)
			}
		}
	}
	return allowedPosts, nil
}

func DbGetCategoryId(categoryTitle string) (category ExistingId, err error) {
	query := "SELECT id FROM Categories WHERE title = ?"
	row := db.DBC.QueryRow(query, categoryTitle)
	err = row.Scan(&category.Id)
	if err != nil {
		return category, err
	}
	return category, nil
}

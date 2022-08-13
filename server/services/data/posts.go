package data_services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	groups "real-time-forum/server/services/data/groups"
	val "real-time-forum/server/services/validation"
	"strings"
)

type NewPostS struct {
	Title        string
	Description  string
	Categoryname string
	Image        string
	Privacy      string
	AllowedUsers []string
}

type ExistingPostS struct {
	Title       string
	Description string
	Image       string
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
	User_id       int
	Image         string
	Comments      []CommentsS
	Likes         []VotesS
	// Cookie        CookieS
	Privacy      string
	User         groups.User
	User_image   string
	AllowedUsers []string
}

type OnePostS struct {
	Categoryname string
	Postid       string
}

var newPost NewPostS

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
		return
	}
	userId := groups.UserId(w, r)
	reqHeader := r.Header.Get("header1")
	splitToken := strings.Split(reqHeader, "Token=")
	token := strings.Join(splitToken, "")

	if !ath.AuthUser(token) {
		w.Write([]byte(`{"message": "User not authenticated"}`))
		return
	}
	err := r.ParseMultipartForm(32 << 0) // maxMemory 32MB
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}
	newPost.AllowedUsers = strings.Split(r.Form["allowedUsers"][0], ",")
	newPost.Title = r.Form["title"][0]
	newPost.Description = r.Form["description"][0]
	newPost.Categoryname = r.Form["categoryname"][0]
	newPost.Privacy = r.Form["privacy"][0]

	imagePath, err := groups.HandleImageUpload(r, "post")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	newPost.Image = imagePath

	if !val.ValidatePostData(newPost.Title, newPost.Description) {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}

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

	stmt, err := db.DBC.Prepare(`INSERT INTO Posts(title, content, image, privacy, created_at, category_id, user_id) VALUES(?, ?, ?, ?, datetime("now"), ?, ?)`)
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(newPost.Title, newPost.Description, newPost.Image, newPost.Privacy, existingCategory.Id, existingUser.Id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	res, err := result.LastInsertId()
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}

	stmt, err = db.DBC.Prepare(`INSERT INTO Posts_allowed_users(user_id, post_id) VALUES(?, ?)`)
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}
	stmt.Exec(userId, res)
	if newPost.AllowedUsers[0] != "null" {
		for _, username := range newPost.AllowedUsers {
			var userId int
			row1 := db.DBC.QueryRow("SELECT id FROM Users WHERE username = ?", username)
			err := row1.Scan(&userId)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}

			stmt, err := db.DBC.Prepare(`INSERT INTO Posts_allowed_users(user_id, post_id) VALUES(?, ?)`)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			stmt.Exec(userId, res)
		}
	}

	var response postResponseS
	response.Message = "Post created succesfully"
	response.Id = res
	var jsonData []byte
	jsonData, _ = json.Marshal((response))
	w.Write(jsonData)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
		return
	}
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
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
		return
	}

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
}

var onePost OnePostS

//SINGLE POST
func GetOnePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)


	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}
	json.Unmarshal([]byte(reqBody), &onePost)

	allPostData, err := onePost.DbGetPost()
	if err != nil {
		var emptyPost NeededData
		groups.SendResponse(w, emptyPost)
		return
	}
	allPostData.User, err = groups.GetUserById(allPostData.User_id)
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}

	allPostData.Comments, err = allPostData.DbGetAllComments()
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}

	for index, comment := range allPostData.Comments {
		allPostData.Comments[index].Likes, err = allPostData.Comments[index].DbGetAllVotes()
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}

		allPostData.Comments[index].User, err = groups.GetUserById(comment.User_id)
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
	}

	allPostData.Likes, err = allPostData.DbGetAllVotes()
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}

	var jsonData []byte
	jsonData, _ = json.Marshal(allPostData)
	w.Write(jsonData)
}

func (data OnePostS) DbGetPost() (post NeededData, err error) {
	row := db.DBC.QueryRow("SELECT Posts.id, Posts.title, Posts.image, Posts.content, Posts.created_at, Posts.privacy, Users.username, Users.id, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE Posts.id = ? AND Categories.title = ?", data.Postid, data.Categoryname)
	err = row.Scan(&post.Id, &post.Title, &post.Image, &post.Description, &post.Created_at, &post.Privacy, &post.Username, &post.User_id, &post.CategoryTitle)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (post NeededData) DbGetAllComments() (comments []CommentsS, err error) {
	query := "SELECT Comments.id, Comments.content, Comments.image, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC"
	rows, err := db.DBC.Query(query, onePost.Postid)
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		comment := CommentsS{}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Image, &comment.Created_at, &comment.User_id, &comment.Post_id, &comment.Username)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (comment CommentsS) DbGetAllVotes() (votes []VotesS, err error) {
	query := "SELECT Comment_likes.type, Comment_likes.created_at, Comment_likes.user_id, Comment_likes.comment_id, Users.username FROM Comment_likes INNER JOIN Users ON Users.id = Comment_likes.user_id WHERE comment_id = ?"
	rows, err := db.DBC.Query(query, comment.Id)
	if err != nil {
		return votes, err
	}

	for rows.Next() {
		vote := VotesS{}
		err := rows.Scan(&vote.Type, &vote.Created_at, &vote.User_id, &vote.Post_id, &vote.Username)
		if err != nil {
			return votes, err
		}

		votes = append(votes, vote)
	}
	return votes, nil
}

func (post NeededData) DbGetAllVotes() (votes []VotesS, err error) {
	query := "SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?"
	rows, err := db.DBC.Query(query, post.Id)
	if err != nil {
		return votes, err
	}

	for rows.Next() {
		vote := VotesS{}
		err := rows.Scan(&vote.Type, &vote.Created_at, &vote.User_id, &vote.Post_id, &vote.Username)
		if err != nil {
			return votes, err
		}
		votes = append(votes, vote)

	}
	return votes, nil
}

func (post NeededData) IsUserAllowed(userId int) (isAllowed bool, err error) {
	var matches int
	query := "SELECT COUNT(user_id) FROM Posts_allowed_users WHERE post_id = ? and user_id = ?"
	row := db.DBC.QueryRow(query, post.Id, userId)
	err = row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches == 1, nil
}

package groups

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type AllPostData struct {
	Post Post
	User User
	Likes []Like
	Comments []Comment 
}

type AllCommentData struct {
	Comment Comment
	Likes []Like
}

type ReplyPost struct {
	Id	int64
}
// Return all groups posts and other info about the posts
func Posts(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)

	posts, err := group.GetAllPosts()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	var responseData []AllPostData
	
	for index, post := range posts {
		var allPostData AllPostData
		responseData = append(responseData, allPostData)
		responseData[index].Post = post

		// Get Comments Amount
		responseData[index].Post.Comments_amount, err = post.CountComments()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		// Get All Likes
		responseData[index].Post.Likes, err = post.GetAllLikes()
		if err == sql.ErrNoRows {
			err = nil
			responseData[index].Post.Likes = []Like{}
		}
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}

		for likeIndex, like := range responseData[index].Post.Likes {
			user, err := GetUserById(like.User_id)
			if err != nil {
				w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
				return
			}
			responseData[index].Post.Likes[likeIndex].Username = user.Username
		}

		responseData[index].User, err = GetUserById(post.User_id)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
	}
	SendResponse(w, responseData)
}
// Creates a new post in a group
func NewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 0) // maxMemory 32MB
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}

	var post Post
	post.Title = r.Form["title"][0]
	post.Content = r.Form["description"][0]
	post.Group_id, err = strconv.Atoi(r.Form["GroupID"][0])
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}
	userId := UserId(w, r)

	post.User_id = userId
	imagePath, err := HandleImageUpload(r, "group/posts")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	post.Image = imagePath

	var replyPost ReplyPost
	replyPost.Id, err = post.Insert()
	if err != nil {
		w.Write([]byte(`{"message":"Post request failed"}`))
		return
	}

	SendResponse(w, replyPost)
}
// Gets all info about one post including user, comments (and it's likes and users) and likes
func PostInfo(w http.ResponseWriter, r *http.Request) {
	post := UnmarshalPost(w, r)

	var allPostData AllPostData
	err := post.GetSingle()
	if err != nil {
		var emptyPost Post
		SendResponse(w, emptyPost)
		return
	}
	allPostData.Post = post

	// Get All Post Likes
	allPostData.Likes, err = post.GetAllLikes()
	if err == sql.ErrNoRows {
		err = nil
		allPostData.Likes = []Like{}
	}
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	// Get All Post Like Usernames
	for index, like := range allPostData.Likes {
		user, err := GetUserById(like.User_id)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		allPostData.Likes[index].Username = user.Username
	}

	// Get All Post Comments
	allPostData.Comments, err = post.GetAllComments()
	if err == sql.ErrNoRows {
		err = nil
		allPostData.Comments = []Comment{}
	}
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	// Get All Post Comment Likes And Comment Users
	for index, comment := range allPostData.Comments {
		likes, err := comment.GetAllLikes()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		allPostData.Comments[index].Likes = likes

		user, err := GetUserById(comment.User_id)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		allPostData.Comments[index].User = user

		// Get All Post Comment Like Usernames	
		for likeIndex, like := range allPostData.Comments[index].Likes {
			user, err := GetUserById(like.User_id)
			if err != nil {
				w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
				return
			}
			allPostData.Comments[index].Likes[likeIndex].Username = user.Username
		}
	}	

	allPostData.User, err = GetUserById(post.User_id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	allPostData.Post.Comments_amount = len(allPostData.Comments)

	SendResponse(w, allPostData)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	like := UnmarshalLike(w, r)
	like.User_id = UserId(w, r)
	like.Type = "Like"

	NewPostVote(w, r, like)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	like := UnmarshalLike(w, r)
	like.User_id = UserId(w, r)
	like.Type = "Dislike"	
	
	NewPostVote(w, r, like)	
}
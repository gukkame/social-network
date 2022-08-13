package groups

import (
	"fmt"
	"net/http"
	"strconv"
)
func NewComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 0) // maxMemory 32MB
	if err != nil {
		w.Write([]byte(`{"message": "Malicious user detected"}`))
		return
	}

	var comment Comment 
	comment.Content = r.Form["Content"][0]
	comment.Post_id, err = strconv.Atoi(r.Form["Post_id"][0])
	if err != nil {
		w.Write([]byte(`{"message": "Post request failed"}`))
		return
	}
	userId := UserId(w, r)

	comment.User_id = userId
	imagePath, err := HandleImageUpload(r, "group/posts/comments")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	comment.Image = imagePath

	err = comment.Insert()
	if err != nil {
		w.Write([]byte(`{"message":"Post request failed"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	like := UnmarshalLike(w, r)
	like.User_id = UserId(w, r)
	like.Type = "Like"

	NewCommentVote(w, r, like)
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	like := UnmarshalLike(w, r)
	like.User_id = UserId(w, r)
	like.Type = "Dislike"	
	
	NewCommentVote(w, r, like)	
}
package groups

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Invite struct {
	Username string `json:"username"`
	GroupId  int    `json:"group_id"`
}

func SendResponse(w http.ResponseWriter, data interface{}) {
	responseData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(responseData)
}

func GetGroupInfo(group Group) (Group, error) {
	group, err := group.GetSingle()
	if err != nil {
		return group, err
	}

	group, err = group.GetOwner()
	if err != nil {
		return group, err
	}

	group, err = group.GetMembersCount()
	if err != nil {
		return group, err
	}
	return group, nil
}
// Returns "Owner", "Member", "invited", "Requested" or "Not Requested"
func (group Group) GetUserStatus(r *http.Request) (string, error) {
	reqHeader := r.Header.Get("header1")
	if len(reqHeader) == 0 {
		return "Not Requested", nil
	}
	splitToken := strings.Split(reqHeader, "Token=")
	token := strings.Join(splitToken, "")

	userId, err := GetUserID(token)
	if err != nil {
		return "", err
	}

	userStatus, err := group.UserRequestStatus(userId)
	if err != nil {
		return "", err
	}

	return userStatus, nil
}

func UnmarshalGroup(w http.ResponseWriter, r *http.Request) (group Group) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &group)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return group
}

func UnmarshalMember(w http.ResponseWriter, r *http.Request) (member Member) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &member)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return member
}

func UnmarshalEvent(w http.ResponseWriter, r *http.Request) (event Event) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &event)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return event
}

func UserId(w http.ResponseWriter, r *http.Request) (userId int) {
	reqHeader := r.Header.Get("header1")
	splitToken := strings.Split(reqHeader, "Token=")
	token := strings.Join(splitToken, "")

	userId, err := GetUserID(token)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	return userId
}

func UnmarshalInvite(w http.ResponseWriter, r *http.Request) (invite Invite) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &invite)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return invite
}

func UnmarshalPost(w http.ResponseWriter, r *http.Request) (post Post) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &post)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return post
}

func UnmarshalReply(w http.ResponseWriter, r *http.Request) (reply EventReply) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &reply)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return reply
}

func UnmarshalLike(w http.ResponseWriter, r *http.Request) (like Like) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	err = json.Unmarshal(body, &like)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	return like
}

func HandleImageUpload(r *http.Request, fileSavePath string) (imagePath string, err error) {
	image, header, err := r.FormFile("img")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return imagePath, err
	}

	if header.Size > 1048576 {
		return imagePath, errors.New("image is too big")
	}

	defer image.Close()
	fileNameParts := strings.Split(header.Filename, ".")

	fileType := fileNameParts[len(fileNameParts)-1]
	fileName := uuid.New().String()
	imagePath = fmt.Sprintf("./resources/%s/%s.%s", fileSavePath, fileName, fileType) 

	out, err := os.OpenFile(imagePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return imagePath, err
	}

	defer out.Close()
	io.Copy(out, image)

	return imagePath[1:], nil
}

func (like Like) DoesPostVoteExist() (exists bool, err error) {
	err = like.GetSinglePostVote()	
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}	
	return true, nil
}

func (like Like) DoesCommentVoteExist() (exists bool, err error) {
	err = like.GetSingleCommentVote()	
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	return true, nil
} 

func NewPostVote(w http.ResponseWriter, r *http.Request, like Like) {
	exists, err := like.DoesPostVoteExist()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	if exists {
		var existingVote Like
		existingVote.Post_id = like.Post_id
		existingVote.User_id = like.User_id

		err = existingVote.GetSinglePostVote()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}

		if existingVote.Type == like.Type {
			err := like.DeletePostVote()
			if err != nil {
				w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		err := like.UpdatePostVote()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	err = like.InsertPostVote()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewCommentVote(w http.ResponseWriter, r *http.Request, like Like) {
	exists, err := like.DoesCommentVoteExist()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	if exists {
		var existingVote Like
		existingVote.Comment_id = like.Comment_id
		existingVote.User_id = like.User_id

		err = existingVote.GetSingleCommentVote()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}

		if existingVote.Type == like.Type {
			err := like.DeleteCommentVote()
			if err != nil {
				w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		err := like.UpdateCommentVote()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	err = like.InsertCommentVote()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewEventReply(w http.ResponseWriter, r *http.Request, reply EventReply) {
	exists, err := reply.DoesEventReplyExist()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	if exists {
		var existingReply EventReply
		existingReply.Event_id = reply.Event_id
		existingReply.User_id = reply.User_id

		err = existingReply.GetSingle()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}

		if existingReply.Status == reply.Status {
			err := reply.Delete()
			if err != nil {
				w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		err := reply.Update()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	err = reply.Insert()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (reply EventReply) DoesEventReplyExist() (exists bool, err error) {
	err = reply.GetSingle()	
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}	
	return true, nil
}
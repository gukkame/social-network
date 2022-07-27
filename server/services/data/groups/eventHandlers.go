package groups

import (
	"fmt"
	"net/http"
)
type AllEventData struct {
	Event	Event	`json:"Event"`
	Replies	[]EventReply	`json:"Replies"`
}

func Events(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)	

	events, err := group.GetAllEvents()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	var responseData []AllEventData

	for index, event := range events {
		responseData = append(responseData, AllEventData{})
		responseData[index].Event = event
		responseData[index].Replies, err = event.GetAllReplies()	
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		user, err := GetUserById(event.User_id)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		responseData[index].Event.User = user

		for replyIndex, reply := range responseData[index].Replies {
			user, err := GetUserById(reply.User_id)
			if err != nil {
				w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
				return
			}
			responseData[index].Replies[replyIndex].Username = user.Username
		}
	}
	SendResponse(w, responseData)
}

func NewEvent(w http.ResponseWriter, r *http.Request) {
	event := UnmarshalEvent(w, r)
	event.User_id = UserId(w, r)

	err := event.Insert()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}
func ReplyGoingEvent(w http.ResponseWriter, r *http.Request) {
	reply := UnmarshalReply(w, r)
	reply.User_id = UserId(w, r)
	reply.Status = "Going"	
	
	NewEventReply(w, r, reply)	
}
func ReplyNotGoingEvent(w http.ResponseWriter, r *http.Request) {
	reply := UnmarshalReply(w, r)
	reply.User_id = UserId(w, r)
	reply.Status = "Not Going"	
	
	NewEventReply(w, r, reply)	
}
package data_services

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type ContentS struct {
	Message  string
	Sender   string
	Receiver string
	Datetime string
	IsGroup  int
	Image    string
}

type ReceivedDataS struct {
	Type     string
	Content  ContentS
	User1    string
	User2    string
	MsgCount int
	Group    bool
}

//WEBSOCKET CONNECTION UPGRADER
var upgrader = websocket.Upgrader{
	//DONT KNOW WHAT NEXT 2 LINES MEAN :)
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//ALLOWS WS CONNECTION TO BE ESTABLISHED FROM EVERY SOURCE (NORMALLY NOT ALLOWED FROM EVERY SOURCE)
	CheckOrigin: func(r *http.Request) bool { return true },
}

//EVENT LOOP
func ListenToMsgs(conn *websocket.Conn) {
	go func() {
		for {
			var newMessage ReceivedDataS
			err := conn.ReadJSON(&newMessage)
			if err != nil {
				RemoveChatUser(conn)
				break
			}

			//Notifications
			if newMessage.Type == "following" {
				FollowNotif(conn, newMessage.User1, newMessage.User2)
			}
			if newMessage.Type == "GroupInvNotif" {
				GroupInvNotif(conn, newMessage)
			}
			if newMessage.Type == "GroupJoinReqNotif" {
				GroupJoinReqNotif(conn, newMessage)
			}
			if newMessage.Type == "NewEventNotif" {
				EventCreatedNotif(conn, newMessage)
			}
			if newMessage.Type == "allNotifications" {
				err := ListAllNotif(conn)
				if err != nil {
					fmt.Println(err)
				}
			}

			//Chat system
			if newMessage.Type == "allClients" {
				err := ListAllClients(conn)
				if err != nil {
					fmt.Println(err)
				}
			}
			if newMessage.Type == "allGroups" {
				err := ListAllGroups(conn)
				if err != nil {
					fmt.Println(err)
				}
			}

			if newMessage.Type == "closeClient" {
				RemoveChatUser(conn)
				conn.Close()
				break
			}

			if newMessage.Type == "privateMSG" {
				SendMessage(conn, newMessage)
			}
			if newMessage.Type == "groupMSG" {
				SendGroupMessage(conn, newMessage)
			}

			if newMessage.Type == "messageHistory" {
				GetPreviousMessages(conn, newMessage.MsgCount, newMessage.User1, newMessage.User2, newMessage.Type)
			}

			if newMessage.Type == "messageHistoryGroup" {
				GetPrevMsgGroup(conn, newMessage.MsgCount, newMessage.User1, newMessage.User2, newMessage.Type)
			}

			if newMessage.Type == "typing" {
				for _, user := range allClients.Clients {
					if user.Username == newMessage.User2 && user.Status == "1" {
						TypingInProgress(newMessage)
					}
				}

			}
		}
	}()

}

//INITIATES A WEBSOCKET CONNECTION FROM NORMAL HTTP REQUEST
func InitiateChat(w http.ResponseWriter, r *http.Request) {
	reqHeader := r.Header.Get("Cookie")
	splitToken := strings.Split(reqHeader, "Token=")
	token := strings.Join(splitToken, "")
	realToken := strings.Split(token, ":")

	//REMOVES USER FROM ACTIVE USERS SLICE, IF IT HAD PREVIOUSLY CRASHED
	CheckForCrashedUser(realToken[1])
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//CREATES A CHAT USER
	CreateChatUser(realToken[0], realToken[1], conn)

	//PUTS THE USER INTO AN INF EVENT LOOP
	ListenToMsgs(conn)
}

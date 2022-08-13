package data_services

import (
	"fmt"
	"log"
	db "real-time-forum/server/db"
	"strconv"

	"github.com/gorilla/websocket"
)

type Notif struct {
	Type       string
	Content    string
	Sender     string
	Recipient  string
	Group_id   string
	Created_at string
}
type LstAllNotif struct {
	Type  string
	Notif []Notif
}

func FollowNotif(currentConn *websocket.Conn, user1 string, user2 string) {
	var response ResponseData
	var user2Conn *websocket.Conn
	response.Type = "followNotif"
	response.Content.Sender = user1
	response.Content.Receiver = user2
	response.Content.Message = user1 + " has requested to follow you!"
	status := false

	user1ID, err := GetUserId(user1)
	if err != nil {
		fmt.Println(err)
	}
	user2ID, err := GetUserId(user2)
	if err != nil {
		fmt.Println(err)
	}

	//Checks if user is online and gets its conn
	for _, user := range allClients.Clients {
		if user.Username == user2 && user.Status == "1" {
			user2Conn = user.Conn
			status = true
		}
	}

	//Check if this notification already exists in db, if yes then updates time
	var one string
	err = db.DBC.QueryRow("SELECT type FROM Notifications WHERE type = ? AND user_id = ? AND recipient_id = ?", response.Type, user1ID, user2ID).Scan(&one)
	if err != nil {
		one = "0"
	}
	if one == "0" {
		//Save notification
		stmt, err := db.DBC.Prepare(`INSERT INTO Notifications(type, content, user_id, recipient_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, datetime("now"))`)
		if err != nil {
			fmt.Println("Didnt managed to save notification in database", err)
		}
		stmt.Exec(response.Type, response.Content.Message, user1ID, user2ID, 0)
		defer stmt.Close()

	} else {

		stmt, _ := db.DBC.Prepare(`DELETE FROM Notifications WHERE type = ? AND user_id = ? AND recipient_id = ?`)
		stmt.Exec(response.Type, user1ID, user2ID)
		defer stmt.Close()

		stmt, err := db.DBC.Prepare(`INSERT INTO Notifications(type, content, user_id, recipient_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, datetime("now"))`)
		if err != nil {
			fmt.Println("Didnt managed to save notification in database", err)
		}
		stmt.Exec(response.Type, response.Content.Message, user1ID, user2ID, 0)
		defer stmt.Close()
	}
	if status {
		err1 := user2Conn.WriteJSON(response)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}
func GroupInvNotif(currentConn *websocket.Conn, data ReceivedDataS) {
	var response ReceivedDataS
	var user2Conn *websocket.Conn
	response.Type = "GroupInvNotif"
	response.Content.Sender = data.Content.Sender
	response.Content.Receiver = data.Content.Receiver
	response.Content.IsGroup = data.Content.IsGroup
	response.Content.Message = data.Content.Sender + " has invited you to join the group!"
	status := false

	isGroup := 0
	SaveNotif(response, isGroup, data.Content.IsGroup)

	for _, user := range allClients.Clients {
		if user.Username == data.Content.Receiver && user.Status == "1" {
			user2Conn = user.Conn
			status = true
		}
	}
	if status {
		err1 := user2Conn.WriteJSON(response)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}
func GroupJoinReqNotif(currentConn *websocket.Conn, data ReceivedDataS) {
	var response ReceivedDataS
	var user2Conn *websocket.Conn
	response.Type = "GroupJoinReqNotif"
	response.Content.Sender = data.Content.Sender
	response.Content.IsGroup = data.Content.IsGroup

	row, err := db.DBC.Query("SELECT Users.username FROM Users INNER JOIN Group_users ON Group_users.user_id = Users.id WHERE Group_users.group_id = ? AND Group_users.status = ?", data.Content.IsGroup, "Owner")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var groupOwner string
	for row.Next() {
		err2 := row.Scan(&groupOwner)
		if err2 != nil {
			log.Fatal(err)
			return
		}

	}

	response.Content.Receiver = groupOwner
	response.Content.Message = data.Content.Sender + " requested to join the group!"
	status := false

	isGroup := 0
	SaveNotif(response, isGroup, data.Content.IsGroup)

	for _, user := range allClients.Clients {
		if user.Username == groupOwner && user.Status == "1" {
			user2Conn = user.Conn
			status = true
		}
	}

	if status {

		err1 := user2Conn.WriteJSON(response)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}
func EventCreatedNotif(currentConn *websocket.Conn, data ReceivedDataS) {
	var response ReceivedDataS
	var user2Conn []ClientS

	response.Type = "NewEventNotif"
	response.Content.Sender = data.Content.Sender
	response.Content.IsGroup = data.Content.IsGroup
	response.Content.Receiver = strconv.Itoa(data.Content.IsGroup)
	status := false
	name, err := GetGroupName(data.Content.IsGroup)
	if err != nil {
		fmt.Println(err)
		return
	}
	response.Content.Message = "New event at " + name + " group!"

	isGroup := 1
	SaveNotif(response, isGroup, data.Content.IsGroup)

	row, err1 := db.DBC.Query("SELECT Users.username FROM Users INNER JOIN Group_users ON Group_users.user_id = Users.id WHERE Group_users.group_id = ? ", data.Content.IsGroup)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	defer row.Close()

	for row.Next() {
		var user2 ClientS
		err := row.Scan(&user2.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
		if user2.Username != data.Content.Sender {
			user2Conn = append(user2Conn, user2)
		}
	}

	for _, onlineUser := range allClients.Clients {
		for i, groupMember := range user2Conn {
			if onlineUser.Username == groupMember.Username && onlineUser.Status == "1" {
				user2Conn[i].Conn = onlineUser.Conn
				status = true
				user2Conn[i].IsGroup = 1
			}
		}
	}

	if status {
		for _, user := range user2Conn {
			if user.Conn != nil {
				err1 := user.Conn.WriteJSON(response)
				if err1 != nil {
					fmt.Println(err1)
				}
			}

		}
		
	}
}

func SaveNotif(data ReceivedDataS, group int, group_id int) {

	var user2ID int
	user1ID, err := GetUserId(data.Content.Sender)
	if err != nil {
		fmt.Println(err)
	}
	if group == 0 {
		user2ID, err = GetUserId(data.Content.Receiver)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		user2ID = 0
	}

	//Check if this notification already exists in db, if yes then updates time

	//Save notification
	stmt, err := db.DBC.Prepare(`INSERT INTO Notifications(type, content, user_id, recipient_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, datetime("now"))`)
	if err != nil {
		fmt.Println("Didnt managed to save notification in database", err)
	}

	stmt.Exec(data.Type, data.Content.Message, user1ID, user2ID, group_id)
	defer stmt.Close()

}

func ListAllNotif(conn *websocket.Conn) error {
	var userID int
	var err error
	var listAllNotif LstAllNotif

	var allNotif []Notif

	for _, item := range allClients.Clients {
		if item.Conn == conn {
			userID, err = GetUserId(item.Username)
			if err != nil {
				return err
			}
		}
	}

	row, err := db.DBC.Query("SELECT type, content, user_id, recipient_id, group_id, created_at FROM Notifications WHERE recipient_id = ? OR recipient_id = ? ORDER BY created_at DESC ", userID, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer row.Close()

	for row.Next() {
		var notification Notif
		err := row.Scan(&notification.Type, &notification.Content, &notification.Sender, &notification.Recipient, &notification.Group_id, &notification.Created_at)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if notification.Recipient == "0" {
			var one string
			err = db.DBC.QueryRow("SELECT user_id FROM Group_users WHERE group_id = ? AND user_id = ? AND status = ? OR group_id = ? AND user_id = ? AND status = ?", notification.Group_id, userID, "Owner", notification.Group_id, userID, "Member").Scan(&one)
			if err != nil {
				one = "0"
			}
			defer row.Close()
			if one != "0" {
				allNotif = append(allNotif, notification)
			}
		} else {
			allNotif = append(allNotif, notification)
		}
	}

	if len(allNotif) <= 10 {
		listAllNotif.Notif = allNotif
	} else {
		listAllNotif.Notif = allNotif[:10]
	}
	listAllNotif.Type = "ListAllNotif"

	err = conn.WriteJSON(listAllNotif)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

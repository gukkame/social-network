package data_services

import (
	"fmt"
	db "real-time-forum/server/db"

	"github.com/gorilla/websocket"
)

type ResponseData struct {
	Type    string
	Content ContentS
}

type ExistingId struct {
	Id int
}

type MsgHistory struct {
	Type    string
	Content []ContentS
}

func SendMessage(currentConn *websocket.Conn, incData ReceivedDataS) {
	var response ResponseData
	var user2Conn *websocket.Conn
	response.Type = incData.Type
	response.Content = incData.Content
	for _, item := range allClients.Clients {
		if item.Username == incData.User2 {
			user2Conn = item.Conn
		}

	}

	date := SaveMessage(incData.Content.Message, incData.User1, incData.User2)
	response.Content.Datetime = date
	response.Content.Receiver = incData.User2
	currentConn.WriteJSON(response)
	user2Conn.WriteJSON(response)
}

func SaveMessage(Content string, User1 string, User2 string) string {
	user1ID, err := GetUserId(User1)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	user2ID, err := GetUserId(User2)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	stmt, err := db.DBC.Prepare(`INSERT INTO Messages(content, created_at, sender_id, recipient_id) VALUES(?, datetime("now"), ?, ?)`)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	result, _ := stmt.Exec(Content, user1ID, user2ID)
	res, err := result.LastInsertId()
	if err != nil {
		return ""
	}

	row, err := db.DBC.Query("SELECT Messages.created_at FROM Messages WHERE Messages.id = ?", res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer row.Close()

	var oneMessage ContentS
	for row.Next() {
		err := row.Scan(
			&oneMessage.Datetime,
		)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	defer stmt.Close()
	return oneMessage.Datetime
}

func GetPreviousMessages(conn *websocket.Conn, Count int, User1 string, User2 string, Type string) {
	user1ID, err := GetUserId(User1)
	if err != nil {
		fmt.Println(err)
		return
	}

	user2ID, err := GetUserId(User2)
	if err != nil {
		fmt.Println(err)
		return
	}

	Count = Count + 10

	var finalResponse MsgHistory
	row, err := db.DBC.Query("SELECT Messages.content, Messages.created_at, Users.username FROM Messages INNER JOIN Users ON Users.id = Messages.sender_id WHERE Messages.sender_id = ? AND Messages.recipient_id = ? OR Messages.sender_id = ? AND Messages.recipient_id = ? ORDER BY Messages.created_at DESC LIMIT  ?", user1ID, user2ID, user2ID, user1ID, Count)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer row.Close()

	var msges []ContentS
	for row.Next() {
		var oneMessage ContentS
		err := row.Scan(
			&oneMessage.Message,
			&oneMessage.Datetime,
			&oneMessage.Sender,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		msges = append(msges, oneMessage)
	}

	finalResponse.Type = Type
	finalResponse.Content = msges
	conn.WriteJSON(finalResponse)
}

func TypingInProgress(incData ReceivedDataS) {
	var response ResponseData
	var user2Conn *websocket.Conn
	response.Type = incData.Type
	response.Content = incData.Content
	for _, item := range allClients.Clients {
		if item.Username == incData.User2 {
			user2Conn = item.Conn
		}

	}

	user2Conn.WriteJSON(response)
}

//GET USERID FROM DATABASE USING USERNAME
func GetUserId(User string) (int, error) {
	row, err := db.DBC.Query("SELECT Users.id FROM Users WHERE username = ?", User)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer row.Close()

	var user ExistingId
	for row.Next() {
		err := row.Scan(
			&user.Id,
		)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}

	return user.Id, err
}

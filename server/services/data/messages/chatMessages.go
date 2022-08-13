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
type ExistingName struct {
	Name string
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
	status := false
	for _, user := range allClients.Clients {
		if user.Username == incData.User2 && user.Status == "1" {
			user2Conn = user.Conn
			status = true
		}
	}
	date := SaveMessage(incData.Content.Message, incData.User1, incData.User2)
	response.Content.Datetime = date
	response.Content.Receiver = incData.User2
	err := currentConn.WriteJSON(response)
	if err != nil {
		fmt.Println(err)
	}
	if status {
		err1 := user2Conn.WriteJSON(response)
		if err1 != nil {
			fmt.Println(err1)
		}
	}

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

func GetPrevMsgGroup(conn *websocket.Conn, Count int, User1 string, GroupName string, Type string) {

	groupID, err := GetGroupId(GroupName)
	if err != nil {
		fmt.Println(err)
		return
	}

	Count = Count + 10

	var finalResponse MsgHistory
	row, err := db.DBC.Query("SELECT Group_messages.content, Group_messages.created_at, Users.username FROM Group_messages INNER JOIN Users ON Users.id = Group_messages.sender_id WHERE Group_messages.group_id = ? ORDER BY Group_messages.created_at DESC LIMIT  ?", groupID, Count)
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
		oneMessage.IsGroup = 1
		
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
	err := user2Conn.WriteJSON(response)
	if err != nil {
		fmt.Println(err)
		return
	}
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

//GET USERNAME FROM DATABASE USING USERID
func GetUsername(Userid string) (string, error) {
	row, err := db.DBC.Query("SELECT Users.username FROM Users WHERE id = ?", Userid)
	if err != nil {
		fmt.Println(err)
		return "0", err
	}
	defer row.Close()

	var user ExistingName
	for row.Next() {
		err := row.Scan(
			&user.Name,
		)
		if err != nil {
			fmt.Println(err)
			return "0", err
		}
	}

	return user.Name, err
}

//GET GROUPID FROM DB USING GROUP NAME
func GetGroupId(groupName string) (int, error) {
	row, err := db.DBC.Query("SELECT User_groups.id FROM User_groups WHERE name = ?", groupName)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer row.Close()

	var id ExistingId
	for row.Next() {
		err := row.Scan(
			&id.Id,
		)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}

	return id.Id, err
}
//GET GROUP NAME FROM DB USING GROUP NAME
func GetGroupName(groupid int) (string, error) {
	row, err := db.DBC.Query("SELECT User_groups.name FROM User_groups WHERE id = ?", groupid)
	if err != nil {
		fmt.Println(err)
		return "0", err
	}
	defer row.Close()

	var name ExistingName
	for row.Next() {
		err := row.Scan(
			&name.Name,
		)
		if err != nil {
			fmt.Println(err)
			return "0", err
		}
	}

	return name.Name, err
}

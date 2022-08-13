package data_services

import (
	"fmt"
	db "real-time-forum/server/db"

	"github.com/gorilla/websocket"
)

type GroupInfo struct {
	Name        string
	Image       string
	Id          int
	LastMessage string
}

type AllGroupsN struct {
	Groups []GroupInfo
}

type LstAllGroups struct {
	Type   string
	Groups AllGroupsN
}

func SendGroupMessage(currentConn *websocket.Conn, incData ReceivedDataS) {
	var response ResponseData
	var user2Conn []ClientS

	response.Type = incData.Type
	response.Content = incData.Content
	status := false

	groupID, err := GetGroupId(incData.User2)
	if err != nil {
		fmt.Println(err)
		return
	}
	row, err1 := db.DBC.Query("SELECT Users.username FROM Users INNER JOIN Group_users ON Group_users.user_id = Users.id WHERE Group_users.group_id = ? ", groupID)
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
		if user2.Username != incData.User1 {
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

	date := SaveGroupMessage(incData.Content.Message, incData.User1, incData.User2)
	response.Content.Datetime = date
	response.Content.Receiver = incData.User2
	response.Content.IsGroup = 1
	err = currentConn.WriteJSON(response)
	if err != nil {
		fmt.Println(err)
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

func SaveGroupMessage(Content string, User1 string, User2 string) string {
	user1ID, err := GetUserId(User1)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	groupID, err := GetGroupId(User2)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	stmt, err := db.DBC.Prepare(`INSERT INTO Group_messages(content, created_at, sender_id, group_id) VALUES(?, datetime("now"), ?, ?)`)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	result, _ := stmt.Exec(Content, user1ID, groupID)
	res, err := result.LastInsertId()
	if err != nil {
		return ""
	}

	row, err := db.DBC.Query("SELECT Group_messages.created_at FROM Group_messages WHERE Group_messages.id = ?", res)
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

//Sends back list of all groups names, last message and all the members!
func ListAllGroups(conn *websocket.Conn) error {
	var userID int
	var err error
	var groupsAll LstAllGroups
	groupsAll.Type = "allGroups"

	//All online users
	for _, item := range allClients.Clients {
		if item.Conn == conn {
			userID, err = GetUserId(item.Username)
			if err != nil {
				return err
			}
		}

	}

	//All group names! user is following to
	// var allgroups []string
	var group GroupInfo
	var allgroups AllGroupsN
	row2, err := db.DBC.Query("SELECT id, name, group_image FROM User_groups INNER JOIN Group_users ON Group_users.group_id = User_groups.id WHERE Group_users.user_id = ? AND status = ? OR Group_users.user_id = ? AND status = ?", userID, "Member", userID, "Owner")

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer row2.Close()

	for row2.Next() {

		err := row2.Scan(&group.Id, &group.Name, &group.Image)
		if err != nil {
			fmt.Println(err)
			return err
		}
		allgroups.Groups = append(allgroups.Groups, group)

	}
	groupsAll.Groups = allgroups

	err = conn.WriteJSON(groupsAll)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

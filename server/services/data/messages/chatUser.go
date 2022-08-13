package data_services

import (
	"fmt"
	db "real-time-forum/server/db"
	"strconv"

	"github.com/gorilla/websocket"
)

type ClientS struct {
	Username    string
	Avatar_image	string
	Conn        *websocket.Conn
	LastMessage string
	Status      string
	IsGroup int
}

type AllClientsS struct {
	Clients []ClientS
}

type ListAllClientsS struct {
	Type    string
	Clients AllClientsS
}
type Follower struct {
	follower_id  string
	recipient_id string
}

var allClients AllClientsS

func CreateChatUser(token string, username string, conn *websocket.Conn) {
	var newWSclient ClientS
	newWSclient.Username = username
	newWSclient.Conn = conn
	newWSclient.Status = "1"
	allClients.Clients = append(allClients.Clients, newWSclient)
}

func RemoveChatUser(conn *websocket.Conn) {
	for index, item := range allClients.Clients {
		if item.Conn == conn {
			allClients.Clients = append(allClients.Clients[0:index], allClients.Clients[index+1:]...)
		}

	}
}

func CheckForCrashedUser(username string) {
	for index, item := range allClients.Clients {
		if item.Username == username {
			allClients.Clients = append(allClients.Clients[0:index], allClients.Clients[index+1:]...)
		}
	}
}

func ListAllClients(conn *websocket.Conn) error {
	var userID int
	var err error
	var clientsAll ListAllClientsS
	clientsAll.Type = "allClients"

	//All online users
	for _, item := range allClients.Clients {
		if item.Conn == conn {
			userID, err = GetUserId(item.Username)
			if err != nil {
				return err
			}
		}

	}

	//All users main user is following to, and they follow back
	var followers1 Follower
	var allusrs []Follower
	row2, err := db.DBC.Query("SELECT Followers.follower_id,Followers.recipient_id FROM Followers WHERE Followers.follower_id = ? OR Followers.recipient_id = ? ", userID, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer row2.Close()

	for row2.Next() {
		err := row2.Scan(&followers1.follower_id, &followers1.recipient_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		allusrs = append(allusrs, followers1)

	}

	for _, user := range allusrs {
		for _, user2 := range allusrs {
			if user.follower_id == user2.recipient_id && user2.follower_id == user.recipient_id {
				if user.follower_id == strconv.Itoa(userID) {
					username, err := GetUsername(user2.follower_id)
					if err != nil {
						return err
					}
					double := false
					for _, onlineUser := range allClients.Clients {
						if onlineUser.Username == username {
							double = true
						}
					}
					if !double{
						allClients.Clients = append(allClients.Clients, ClientS{Username: username, Status: "0"})
						break
					}

				}
			}
		}
	}

	var Allclients AllClientsS
	var clients []ClientS
	for _, item := range allClients.Clients {
		var allfollowers []Follower
		if item.Conn != conn {
			var client ClientS
			client.Username = item.Username
			userID2, err2 := GetUserId(item.Username)
			if err2 != nil {
				return err
			}

			var followers Follower

			row0, err := db.DBC.Query("SELECT Followers.follower_id,Followers.recipient_id FROM Followers WHERE Followers.follower_id = ? AND Followers.recipient_id = ? OR Followers.follower_id = ? AND Followers.recipient_id = ? ", userID, userID2, userID2, userID)
			if err != nil {
				fmt.Println(err)
				return err
			}
			defer row0.Close()

			for row0.Next() {
				err := row0.Scan(&followers.follower_id, &followers.recipient_id)
				if err != nil {
					fmt.Println(err)
					return err
				}
				allfollowers = append(allfollowers, followers)

			}
			if len(allfollowers) == 2 {
				row, err := db.DBC.Query("SELECT Messages.created_at FROM Messages WHERE Messages.sender_id = ? AND Messages.recipient_id = ? OR Messages.sender_id = ? AND Messages.recipient_id = ? ORDER BY Messages.created_at DESC LIMIT 1", userID, userID2, userID2, userID)
				if err != nil {
					fmt.Println(err)
					return err
				}
				defer row.Close()

				for row.Next() {
					err := row.Scan(
						&client.LastMessage,
					)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}
				//Add img to client
				row1, err := db.DBC.Query("SELECT avatar_image FROM Users WHERE id = ?", userID2)
				if err != nil {
					fmt.Println(err)
					return err
				}
				defer row1.Close()
				for row1.Next() {
					err := row1.Scan(
						&client.Avatar_image,
					)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}

				if len(client.LastMessage) == 0 {
					client.LastMessage = "no date"
				}
				client.Status = item.Status
				clients = append(clients, client)
			}
		}
	}
	Allclients.Clients = clients
	clientsAll.Clients = Allclients

	err3 := conn.WriteJSON(clientsAll)
	if err3 != nil {
		fmt.Println(err)
		return err
	}
	return err
}

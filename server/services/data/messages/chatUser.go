package data_services

import (
	"fmt"
	db "real-time-forum/server/db"

	"github.com/gorilla/websocket"
)

type ClientS struct {
	Username    string
	Conn        *websocket.Conn
	LastMessage string
}

type AllClientsS struct {
	Clients []ClientS
}

type ListAllClientsS struct {
	Type    string
	Clients AllClientsS
}

var allClients AllClientsS

func CreateChatUser(token string, username string, conn *websocket.Conn) {
	var newWSclient ClientS
	newWSclient.Username = username
	newWSclient.Conn = conn
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
	for _, item := range allClients.Clients {
		if item.Conn == conn {
			userID, err = GetUserId(item.Username)
			if err != nil {
				return err
			}
		}

	}

	var Allclients AllClientsS
	var clients []ClientS
	for _, item := range allClients.Clients {
		if item.Conn != conn {
			var client ClientS
			client.Username = item.Username
			userID2, err2 := GetUserId(item.Username)
			if err2 != nil {
				return err
			}

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
			if len(client.LastMessage) == 0 {
				client.LastMessage = "no date"
			}
			clients = append(clients, client)
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

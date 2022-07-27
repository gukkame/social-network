package data_services 

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	val "real-time-forum/server/services/validation"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type GroupUserIds struct {
	userID  int
	groupID int
	userStatus string
}

type Group struct {
	Group_id     int    `json:"GroupID"`
	Title        string `json:"Title"`
	Content      string `json:"Content"`
	Image        string `json:"Image"`
	Created_at   string `json:"Created_at"`
	Creator_id   int    `json:"Creator_id"`
	Creator_name string `json:"Creator_name"`
	Member_count int    `json:"Member_count"`
	Group_member int    `json:"Group_member"`
	User_status  string `json:"User_status"`
}

//Get all groups
func GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)

		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		row, err := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()

		var user_id int
		for row.Next() {
			err2 := row.Scan(&user_id)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}

		}

		//get all users and group ids that they are part of
		rows, err := db.DBC.Query("SELECT * FROM Group_users")
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows.Close()

		var groupUserIds []GroupUserIds
		for rows.Next() {
			var ids GroupUserIds
			err2 := rows.Scan(&ids.userID, &ids.groupID, &ids.userStatus)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
			groupUserIds = append(groupUserIds, ids)
		}

		var groupids []int
		for _, bothIds := range groupUserIds {
			groupids = append(groupids, bothIds.groupID)
		}
		//sum up all members in groups
		memberCount := printUniqueValue(groupids)

		//Get all info about all groups
		rows2, err := db.DBC.Query("SELECT User_groups.id,User_groups.name,	User_groups.content,User_groups.group_image,User_groups.created_at,	User_groups.creator_id,	Users.username FROM User_groups INNER JOIN Users ON Users.id = User_groups.creator_id")
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows2.Close()

		var groupsInfo []Group
		for rows2.Next() {
			var groupInfo Group
			err2 := rows2.Scan(&groupInfo.Group_id, &groupInfo.Title, &groupInfo.Content, &groupInfo.Image, &groupInfo.Created_at, &groupInfo.Creator_id, &groupInfo.Creator_name)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
			for group_id, membercount := range memberCount {
				if groupInfo.Group_id == group_id {
					groupInfo.Member_count = membercount
				}
			}
			for _, v := range groupUserIds {
				if v.userID == user_id && v.groupID == groupInfo.Group_id {
					groupInfo.Group_member = 1
				}
				if groupInfo.Group_member != 1 {
					groupInfo.Group_member = 0
				}
			}

			groupsInfo = append(groupsInfo, groupInfo)
		}

		//Send data to front
		var jsonData []byte
		jsonData, _ = json.Marshal((groupsInfo))
		w.Write(jsonData)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func printUniqueValue(arr []int) map[int]int {
	dict := make(map[int]int)
	for _, num := range arr {
		dict[num] = dict[num] + 1
	}
	return dict
}

func NewGroup(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":

		err := req.ParseMultipartForm(32 << 0) // maxMemory 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		title := (req.Form["title"][0])
		content := (req.Form["content"][0])
		splitToken := strings.Split((req.Form["token"][0]), "Token=")
		token := strings.Join(splitToken, "")

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		if val.ValidPostTitle(title) && val.ValidPostDescription(content) {

			row, err := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer row.Close()

			var user_id int
			for row.Next() {
				err2 := row.Scan(&user_id)
				if err2 != nil {
					w.Write([]byte(`{"message": "Post request failed"}`))
					return
				}

			}
			stmt, err := db.DBC.Prepare(`INSERT INTO User_groups(name, content, group_image, created_at, creator_id) VALUES(?, ?, ?, datetime("now"), ?)`)
			if err != nil {
				fmt.Println(err)
				return
			}

			in, header, err := req.FormFile("image")

			if in == nil {
				stmt.Exec(title, content, "", user_id)
				defer stmt.Close()

			} else {
				if header.Header.Get("Content-Type") == "image/gif" {
					w.Write([]byte(`{"message": "Cant have gif as profile picture"}`))
					return
				}
				if header.Size <= 1048576 {

					if err != nil {
						fmt.Println(err)
						return
					}
					defer in.Close()
					id := uuid.New()
					img_id := id.String()
					s := strings.Split(header.Filename, ".")

					out, err := os.OpenFile("./resources/group/"+img_id+"."+s[len(s)-1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("ERORR#", err)
						return
					}

					defer out.Close()
					io.Copy(out, in)

					stmt.Exec(title, content, "/resources/group/"+img_id+"."+s[len(s)-1], user_id)
					defer stmt.Close()
				} else {
					w.Write([]byte(`{"message": "Image is too big!"}`))
					return
				}
			}
			row2, err := db.DBC.Query("SELECT id FROM User_groups WHERE creator_id = ?", user_id)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer row2.Close()

			var group_id int
			for row2.Next() {
				err2 := row2.Scan(&group_id)
				if err2 != nil {
					w.Write([]byte(`{"message": "Post request failed"}`))
					return
				}
			}
			stmt1, err := db.DBC.Prepare(`INSERT INTO Group_users(user_id, group_id, status) VALUES(?, ?, ?)`)
			if err != nil {
				fmt.Println(err)
				return
			}
			stmt1.Exec(user_id, group_id, "Owner")

			//Send data to front
			var jsonData []byte
			jsonData, _ = json.Marshal((group_id))
			w.Write(jsonData)
		} else {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"msg": "Can't find method requested"}`))
	}
}

func (group Group) IsUserOwner(userId string) (bool, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return false, err
	}
	if userIdInt == group.Creator_id {
		return true, nil
	}
	return false, nil
}
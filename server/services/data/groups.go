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
	"strings"
)

type GroupUserIds struct {
	userID  int
	groupID int
}
type GroupInfo struct {
	groupID    int
	title      string
	content    string
	image      string
	created_at string
	creator_id int
	creator_name string
}

type SendInfo struct {
	memberCount  map[int]int
	myGroups     []int
	otherGroups  []int
	allGroupInfo []GroupInfo
}

//Get all groups
func GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)

		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")
		///
		row, err := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()
		///
		var user_id int
		err1 := row.Scan(&user_id)
		fmt.Println(err1)
		///
		rows, err := db.DBC.Query("SELECT * FROM Group_users")
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows.Close()

		var groupUserIds []GroupUserIds
		for rows.Next() {
			var ids GroupUserIds
			err2 := rows.Scan(&ids.userID, &ids.groupID)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
			groupUserIds = append(groupUserIds, ids)
		}
		fmt.Println("groupUserIds(userID, groupID): ", groupUserIds)
		///
		var groupids []int
		var myGroups []int
		var otherGroups []int
		for _, bothIds := range groupUserIds {
			groupids = append(groupids, bothIds.groupID)
			if bothIds.userID == user_id {
				myGroups = append(myGroups, bothIds.groupID)
			} else {
				otherGroups = append(otherGroups, bothIds.groupID)
			}
		}
		memberCount := printUniqueValue(groupids)
		groups := printUniqueValue(otherGroups)
		fmt.Println("membercount in groupID(map): ", memberCount)
		fmt.Println("GroupIDs that I am part of(int[]): ", myGroups)

		// To get Other groups!!!
		allOtherGroups := make([]int, 0, len(groups))
		for k := range groups {
			allOtherGroups = append(allOtherGroups, k)
		}
		fmt.Println("Other groupsID (int[]): ", allOtherGroups)
		// To get group_id and membercount!!!
		// group_id := make([]int, 0, len(memberCount))
		// membercount := make([]int, 0, len(memberCount))
		// for k, v := range memberCount {
		// 	group_id = append(group_id, k)
		// 	membercount = append(membercount, v)
		// }

		// rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", onePost.Postid)

		rows2, err := db.DBC.Query("SELECT User_groups.id,User_groups.name,	User_groups.content,User_groups.group_image,User_groups.created_at,	User_groups.creator_id,	Users.username FROM User_groups INNER JOIN Users ON Users.id = User_groups.creator_id")
		if err != nil {
			w.Write([]byte(`{"message": "Post request failed"}`))
			return
		}
		defer rows2.Close()

		var groupsInfo []GroupInfo
		for rows2.Next() {
			var groupInfo GroupInfo
			err2 := rows2.Scan(&groupInfo.groupID, &groupInfo.title, &groupInfo.content, &groupInfo.image, &groupInfo.created_at, &groupInfo.creator_id, &groupInfo.creator_name)
			if err2 != nil {
				w.Write([]byte(`{"message": "Post request failed"}`))
				return
			}
			groupsInfo = append(groupsInfo, groupInfo)
		}
		fmt.Println("All group info([]GroupInfo): ", groupsInfo)
		// var sendinfo SendInfo
		sendinfo := SendInfo{memberCount: memberCount, myGroups: myGroups, otherGroups: allOtherGroups, allGroupInfo: groupsInfo}
		///Send data to front
		var jsonData []byte
		jsonData, _ = json.Marshal((sendinfo))
		w.Write(jsonData)

		// var jsonData []byte
		// jsonData, _ = json.Marshal(sendinfo)
		// w.Write(jsonData)

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
	fmt.Println(dict)
	return dict
}

func NewGroup(w http.ResponseWriter, req *http.Request) {

	SetupCORS(&w, req)
	if (*req).Method == "OPTIONS" {
		return
	}

	switch req.Method {
	case "POST":

		err := req.ParseMultipartForm(32 << 0) // maxMemory 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		title := (req.Form["title"][0])
		content := (req.Form["content"][0])
		token := (req.Form["token"][0])

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
			err1 := row.Scan(&user_id)
			fmt.Println(err1)

			stmt, err := db.DBC.Prepare(`INSERT INTO User_groups(name, content, group_image, created_at, creator_id) VALUES(?, ?, datetime("now"), ?, ?)`)
			if err != nil {
				fmt.Println(err)
				return
			}

			in, header, err := req.FormFile("image")
			if in == nil {
				stmt.Exec(title, content, "", user_id)
				defer stmt.Close()
			} else {
				// fmt.Println("Image size ", header.Size)
				if header.Size <= 1048576 {

					if err != nil {
						fmt.Println(err)
						return
					}
					defer in.Close()

					s := strings.Split(header.Filename, ".")
					out, err := os.OpenFile("./resources/group/group_img_"+title+"."+s[len(s)-1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					fmt.Println("Path to image", "./resources/group/"+title+"."+s[len(s)-1])
					if err != nil {
						fmt.Println("ERORR#", err)
						return
					}

					defer out.Close()
					io.Copy(out, in)

					stmt.Exec(title, content, "/server/resources/group/"+title+"."+s[len(s)-1], user_id)
					defer stmt.Close()
				} else {
					w.Write([]byte(`{"message": "Image is too big!"}`))
					return
				}
			}
			w.Write([]byte(`{"message": "Data inserted"}`))
		} else {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"msg": "Can't find method requested"}`))
	}

	msg := "heyy!:))"

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		fmt.Println("Last")
		panic(err)
	}

}

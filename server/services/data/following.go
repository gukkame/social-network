package data_services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	mw "real-time-forum/server/middleware"
	"strings"
)

type Follower struct {
	Username     string
	Avatar_image string
}

type FollowerPackage struct {
	Following []Follower
	Followers []Follower
	FollowReq []Follower
	Status    string
}

func Followers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var browsed LogInS
		json.Unmarshal([]byte(reqBody), &browsed)

		//CHECK IF PROFILE EXISTS
		row2, _ := db.DBC.Query("SELECT username FROM Users WHERE username = ?", browsed.Username)
		defer row2.Close()

		//IF IT DOESNT EXISTS
		if !row2.Next() {
			w.Write([]byte(`{"message": "Profile does not exist"}`))
			return
		}

		//CHECK IF OWNER
		var owner LogInS
		row3, _ := db.DBC.Query("SELECT username, profile_status FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		defer row3.Close()

		for row3.Next() { //I was lazy to make more structs
			err := row3.Scan(
				&owner.Username, //username
				&owner.Password, //profile_status just in case it is an owner to send it back for UI
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		if owner.Username == browsed.Username {

			if !ath.AuthUser(token) {
				w.Write([]byte(`{"message": "User not authenticated"}`))
				return
			}

			neededData, err := getFollowers(browsed.Username, "owner")
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
			neededData.Status = "owner"

			var jsonData []byte
			jsonData, _ = json.Marshal((neededData))
			w.Write(jsonData)
			return
		}

		//CHECK IF PROFILE PUBLIC OR PRIVATE
		var profile_status LogInS
		row4, _ := db.DBC.Query("SELECT profile_status FROM Users WHERE username = ?", browsed.Username)
		defer row4.Close()

		for row4.Next() {
			err := row4.Scan(
				&profile_status.Username,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		if profile_status.Username == "public" {
			neededData, err := getFollowers(browsed.Username, "notowner")
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
			neededData.Status = "notowner"

			var jsonData []byte
			jsonData, _ = json.Marshal((neededData))
			w.Write(jsonData)
			return

		}
		if profile_status.Username == "private" {
			if CheckFollower(token, browsed.Username) {
				neededData, err := getFollowers(browsed.Username, "notowner")
				if err != nil {
					w.Write([]byte(`{"message": "Profile does not exist"}`))
					return
				}
				neededData.Status = "notowner"

				var jsonData []byte
				jsonData, _ = json.Marshal((neededData))
				w.Write(jsonData)
				return
			} else {
				w.Write([]byte(`{"message": "Not following"}`))
				return
			}
		}

	}
}

func getFollowers(browsed string, status string) (FollowerPackage, error) {
	var FinalPackage FollowerPackage

	var browsedId ExistingUserS
	row, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed)
	defer row.Close()

	for row.Next() {
		err := row.Scan(
			&browsedId.Id,
		)
		if err != nil {
			return FollowerPackage{}, err
		}
	}

	row2, _ := db.DBC.Query("SELECT recipient_id FROM Followers WHERE follower_id = ? AND status = ?", browsedId.Id, "following")
	defer row2.Close()
	var Following []Follower
	for row2.Next() {
		var reciID ExistingUserS
		err := row2.Scan(
			&reciID.Id,
		)
		if err != nil {
			return FollowerPackage{}, err
		}

		row4, _ := db.DBC.Query("SELECT Users.username, Users.avatar_image FROM Users WHERE id = ?", reciID.Id)
		defer row4.Close()
		for row4.Next() {
			var oneFollowing Follower
			err := row4.Scan(
				&oneFollowing.Username,
				&oneFollowing.Avatar_image,
			)
			if err != nil {
				return FollowerPackage{}, err
			}

			Following = append(Following, oneFollowing)
		}

	}
	FinalPackage.Following = Following
	row6, _ := db.DBC.Query("SELECT follower_id FROM Followers WHERE recipient_id = ? AND status = ?", browsedId.Id, "following")
	defer row6.Close()
	var Followers []Follower
	for row6.Next() {
		var reciID ExistingUserS
		err := row6.Scan(
			&reciID.Id,
		)
		if err != nil {
			return FollowerPackage{}, err
		}

		row3, _ := db.DBC.Query("SELECT Users.username, Users.avatar_image FROM Users WHERE id = ?", reciID.Id)
		defer row3.Close()
		for row3.Next() {
			var oneFollower Follower
			err := row3.Scan(
				&oneFollower.Username,
				&oneFollower.Avatar_image,
			)
			if err != nil {
				return FollowerPackage{}, err
			}
			Followers = append(Followers, oneFollower)
		}
		FinalPackage.Followers = Followers

	}

	if status != "owner" {
		return FinalPackage, nil
	}

	row7, _ := db.DBC.Query("SELECT follower_id FROM Followers WHERE recipient_id = ? AND status = ?", browsedId.Id, "requested")
	defer row7.Close()
	var FollowReqs []Follower
	for row7.Next() {
		fmt.Println("1")
		var reciID ExistingUserS
		err := row7.Scan(
			&reciID.Id,
		)
		if err != nil {
			fmt.Println(err)
			return FollowerPackage{}, err
		}

		row4, _ := db.DBC.Query("SELECT Users.username, Users.avatar_image FROM Users WHERE id = ?", reciID.Id)
		defer row4.Close()
		for row4.Next() {
			var oneFollowReq Follower
			err := row4.Scan(
				&oneFollowReq.Username,
				&oneFollowReq.Avatar_image,
			)
			if err != nil {
				return FollowerPackage{}, err
			}
			FollowReqs = append(FollowReqs, oneFollowReq)
		}
		FinalPackage.FollowReq = FollowReqs

	}

	return FinalPackage, nil

}

//CHECKING IF PERSON IS FOLLOWER OR NOT
func CheckFollower(token string, browsed string) bool {
	var browsedId ExistingUserS
	row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed)
	defer row7.Close()

	for row7.Next() {
		err := row7.Scan(
			&browsedId.Id,
		)
		if err != nil {
			return false
		}
	}

	var followerId ExistingUserS
	row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
	defer row5.Close()
	for row5.Next() {
		err := row5.Scan(
			&followerId.Id,
		)
		if err != nil {
			return false
		}
	}

	row6, _ := db.DBC.Query("SELECT status FROM Followers WHERE status = ? AND follower_id = ? AND recipient_id = ?", "following", followerId.Id, browsedId.Id)
	defer row6.Close()

	return row6.Next()
}

func PerformFollow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var browsed LogInS
		json.Unmarshal([]byte(reqBody), &browsed)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		//CHECK IF PROFILE EXISTS
		row2, _ := db.DBC.Query("SELECT username FROM Users WHERE username = ?", browsed.Username)

		//IF IT DOESNT EXISTS
		if !row2.Next() {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}
		row2.Close() // Not defering but instantly closing due to it otherwise locking the database for the next action.

		//CHECK IF PROFILE PUBLIC OR PRIVATE
		var profile_status LogInS
		row4, _ := db.DBC.Query("SELECT profile_status FROM Users WHERE username = ?", browsed.Username)
		defer row4.Close()

		for row4.Next() {
			err := row4.Scan(
				&profile_status.Username,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		if profile_status.Username == "public" {
			if CheckFollower(token, browsed.Username) {
				err := deleteFollower(token, browsed.Username)
				if err != nil {
					w.Write([]byte(`{"message": "Profile does not exist"}`))
					return
				}
			} else {
				err := addFollower(token, browsed.Username, profile_status.Username)
				if err != nil {
					w.Write([]byte(`{"message": "Profile does not exist"}`))
					return
				}
			}
		}

		if profile_status.Username == "private" {
			if CheckFollower(token, browsed.Username) {
				err := deleteFollower(token, browsed.Username)
				if err != nil {
					w.Write([]byte(`{"message": "Profile does not exist"}`))
					return
				}
			} else {
				err := addFollower(token, browsed.Username, profile_status.Username)
				if err != nil {
					w.Write([]byte(`{"message": "Profile does not exist"}`))
					return
				}
			}
		}

	}
}

//FOR OWNER TO ACCEPT A FOLLOW REQUEST AT PROFILE
func AcceptFollower(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var browsed LogInS
		json.Unmarshal([]byte(reqBody), &browsed)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		var browsedId ExistingUserS
		row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed.Username)
		defer row7.Close()

		for row7.Next() {
			err := row7.Scan(
				&browsedId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		var followerId ExistingUserS
		row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		defer row5.Close()
		for row5.Next() {
			err := row5.Scan(
				&followerId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		
		stmt3, err := db.DBC.Prepare(`UPDATE Followers SET status = ? WHERE status = ? AND follower_id = ? AND recipient_id = ?`)
		if err != nil {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
		stmt3.Exec("following","requested", browsedId.Id, followerId.Id)
		defer stmt3.Close()

	}
}

//FOR OWNER TO REMOVE FOLLOWERS
func RemoveFollower(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var browsed LogInS
		json.Unmarshal([]byte(reqBody), &browsed)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		var browsedId ExistingUserS
		row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed.Username)
		defer row7.Close()

		for row7.Next() {
			err := row7.Scan(
				&browsedId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		var followerId ExistingUserS
		row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		defer row5.Close()
		for row5.Next() {
			err := row5.Scan(
				&followerId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}
		stmt, err := db.DBC.Prepare(`DELETE FROM Followers WHERE follower_id = ? AND recipient_id = ?`)
		if err != nil {
			w.Write([]byte(`{"message": "Profile does not exist"}`))
			return
		}

		defer stmt.Close()
		stmt.Exec(browsedId.Id, followerId.Id)
	}
}

//TO DELETE USERS FOLLOWING TO A PERSON IF THEY ARE ALREADY FOLLOWING
func deleteFollower(token string, browsed string) error {
	var browsedId ExistingUserS
	row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed)
	defer row7.Close()

	for row7.Next() {
		err := row7.Scan(
			&browsedId.Id,
		)
		if err != nil {
			return err
		}
	}

	var followerId ExistingUserS
	row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
	defer row5.Close()
	for row5.Next() {
		err := row5.Scan(
			&followerId.Id,
		)
		if err != nil {
			return err
		}
	}
	stmt, err := db.DBC.Prepare(`DELETE FROM Followers WHERE follower_id = ? AND recipient_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	stmt.Exec(followerId.Id, browsedId.Id)
	return err
}

func addFollower(token string, browsed string, profile_status string) error {
	var browsedId ExistingUserS
	row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed)
	defer row7.Close()

	for row7.Next() {
		err := row7.Scan(
			&browsedId.Id,
		)
		if err != nil {
			return err
		}
	}

	var followerId ExistingUserS
	row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
	defer row5.Close()
	for row5.Next() {
		err := row5.Scan(
			&followerId.Id,
		)
		if err != nil {
			return err
		}
	}
	stmt, err := db.DBC.Prepare(`INSERT INTO Followers VALUES(?, datetime('now'), ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	if profile_status == "public" {
		stmt.Exec("following", followerId.Id, browsedId.Id)
	} else {
		stmt.Exec("requested", followerId.Id, browsedId.Id)
	}

	return err
}

//TO CHECK IF USER1 HAS REQUESTED TO FOLLOW USER2
func CheckFollowRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var browsed LogInS
		json.Unmarshal([]byte(reqBody), &browsed)

		//CHECK IF PROFILE EXISTS
		row2, _ := db.DBC.Query("SELECT username FROM Users WHERE username = ?", browsed.Username)

		//IF IT DOESNT EXISTS
		if !row2.Next() {
			w.Write([]byte(`{"message": "Profile does not exist"}`))
			return
		}
		row2.Close()

		var browsedId ExistingUserS
		row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed.Username)
		defer row7.Close()

		for row7.Next() {
			err := row7.Scan(
				&browsedId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		var followerId ExistingUserS
		row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		defer row5.Close()
		for row5.Next() {
			err := row5.Scan(
				&followerId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		row6, _ := db.DBC.Query("SELECT status FROM Followers WHERE status = ? AND follower_id = ? AND recipient_id = ?", "requested", followerId.Id, browsedId.Id)
		defer row6.Close()
		if !row6.Next() {
			w.Write([]byte(`{"message": "Not requested"}`))
			return
		} else {
			w.Write([]byte(`{"message": "Requested"}`))
			return
		}

	}
}

//CANCEL FOLLOW REQUEST 
func CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqHeader := r.Header.Get("header1")
		splitToken := strings.Split(reqHeader, "Token=")
		token := strings.Join(splitToken, "")

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var browsed LogInS
		json.Unmarshal([]byte(reqBody), &browsed)

		if !ath.AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		//CHECK IF PROFILE EXISTS
		row2, _ := db.DBC.Query("SELECT username FROM Users WHERE username = ?", browsed.Username)

		//IF IT DOESNT EXISTS
		if !row2.Next() {
			w.Write([]byte(`{"message": "Profile does not exist"}`))
			return
		}
		row2.Close()

		var browsedId ExistingUserS
		row7, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed.Username)
		defer row7.Close()

		for row7.Next() {
			err := row7.Scan(
				&browsedId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		var followerId ExistingUserS
		row5, _ := db.DBC.Query("SELECT id FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		defer row5.Close()
		for row5.Next() {
			err := row5.Scan(
				&followerId.Id,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		stmt, err := db.DBC.Prepare(`DELETE FROM Followers WHERE status = ? AND follower_id = ? AND recipient_id = ?`)
		if err != nil {
			w.Write([]byte(`{"message": "Profile does not exist"}`))
			return
		}
		defer stmt.Close()
		stmt.Exec("requested", followerId.Id, browsedId.Id)

	}
}

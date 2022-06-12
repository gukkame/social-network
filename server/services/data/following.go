package data_services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	"strings"
)

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
	if !row6.Next() {
		return false
	}
	return true
}

func PerformFollow(w http.ResponseWriter, r *http.Request) {
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

func CheckFollowRequest(w http.ResponseWriter, r *http.Request) {
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

func CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
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

package authentication

import (
	"encoding/json"
	"net/http"
	db "real-time-forum/server/db"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

//ALSO PASTED HERE DUE TO GOLANG RESTRICTON OF NOT ALLOWED TO HAVE PACKAGE IMPORT LOOPS &FROM middleware.go
func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, header1")
}

type ExistingSessionS struct {
	Token       string
	Expiry_date string
	User_id     int
}

type Cookie struct {
	Id       string
	Username string
}

type UsernameS struct {
	Username string
	Token    string
}

func DeleteSession(token string) {
	stmt, _ := db.DBC.Prepare(`DELETE FROM Sessions WHERE token = ?`)
	stmt.Exec(token)
	defer stmt.Close()
}

func GenerateCookieInfo(username string) Cookie {
	newId := uuid.NewV4().String()
	return Cookie{newId, username}
}

func GenerateSession(token_id string, user_id int) error {
	stmt, err := db.DBC.Prepare("INSERT INTO Sessions(token, user_id, expiry_date) VALUES(?, ?, datetime('now','+1 years'))")
	if err != nil {
		return err
	}
	stmt.Exec(token_id, user_id)
	defer stmt.Close()
	return err
}

func AuthUser(token string) bool {
	row, err := db.DBC.Query("SELECT * FROM Sessions WHERE token = ?", token)
	if err != nil {
		DeleteSession(token)
		return false
	}
	defer row.Close()

	var existingSession ExistingSessionS
	for row.Next() {
		err := row.Scan(
			&existingSession.Token,
			&existingSession.Expiry_date,
			&existingSession.User_id,
		)
		if err != nil {
			DeleteSession(token)
			return false
		}

	}

	dateNow := time.Now()

	//GETTING THE REQUIRED DATE FORMAT
	dateSession := existingSession.Expiry_date
	cut1 := strings.Split(dateSession, "T")
	join1 := strings.Join(cut1, " ")
	cut2 := strings.Split(join1, "Z")
	join2 := strings.Join(cut2, "")
	DateTimeSession, err := time.Parse("2006-01-02 15:04:05", join2)
	if err != nil {
		DeleteSession(token)
		return false
	}

	//COMPARING DATETIME NOW TO SESSION EXIPRY_DATE
	if !dateNow.Before(DateTimeSession) {
		DeleteSession(token)
		return false
	}
	return true
}

func DeleteOldSessions() error {
	stmt, err := db.DBC.Prepare(`DELETE FROM Sessions WHERE expiry_date < datetime("now")`)
	if err != nil {
		return err
	}

	stmt.Exec()
	defer stmt.Close()
	return err
}

func CurrentCookieInfo(token string) (string, string) {
	row, _ := db.DBC.Query("SELECT Users.Username, Sessions.token FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
	defer row.Close()

	var existingUser UsernameS
	for row.Next() {
		err := row.Scan(
			&existingUser.Username,
			&existingUser.Token,
		)
		if err != nil {
			return existingUser.Token, existingUser.Username
		}
	}
	return existingUser.Token, existingUser.Username
}

func AuthChat(w http.ResponseWriter, r *http.Request) {
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

		if !AuthUser(token) {
			w.Write([]byte(`{"message": "User not authenticated"}`))
			return
		}

		var payload UsernameS
		username, id := CurrentCookieInfo(token)
		payload.Username = username
		payload.Token = id

		var jsonData []byte
		jsonData, _ = json.Marshal((payload))
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

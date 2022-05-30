package data_services

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	db "real-time-forum/server/db"
	ath "real-time-forum/server/services/authentication"
	val "real-time-forum/server/services/validation"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type NewUserS struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       string
	Gender    string
	Password  string
}

type ExistingUserS struct {
	Id        int
	Username  string
	Password  string
	Email     string
	Age       string
	Gender    string
	FirstName string
	LastName  string
	Date      string
}

type ProfileS struct {
	Username  string
	Email     string
	Age       string
	Gender    string
	FirstName string
	LastName  string
	Token     string
}

type NeverReleasedData struct {
	Id          int
	Password    string
	Date        string
	Expiry_date string
	User_id     int
}

type LogInS struct {
	Username string
	Password string
}

type UsernameS struct {
	Username string `json: "username"`
}

type EmailS struct {
	Email string `json: "email"`
}

var newUser NewUserS

func SignUp(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(reqBody), &newUser)
		if val.ValidateUserData(newUser.Username, newUser.Email, newUser.FirstName, newUser.LastName, newUser.Age, newUser.Gender, newUser.Password) {

			newUser.Password, err = HashPassword(newUser.Password)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}

			stmt, err := db.DBC.Prepare(`INSERT INTO Users(username, password, email, age, gender, first_name, last_name, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, datetime("now"))`)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}

			stmt.Exec(newUser.Username, newUser.Password, newUser.Email, newUser.Age, newUser.Gender, newUser.FirstName, newUser.LastName)
			defer stmt.Close()
			w.Write([]byte(`{"message": "Data inserted"}`))
		} else {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

var existingUser LogInS

func LogIn(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(reqBody), &existingUser)

		if val.ValidateLoginData(existingUser.Username, existingUser.Password) {

			var row *sql.Rows
			if strings.Contains(existingUser.Username, "@") {
				row, err = db.DBC.Query("SELECT * FROM Users WHERE email = ?", existingUser.Username)
				if err != nil {
					w.Write([]byte(`{"message": "Invalid credentials"}`))
					return
				}

			} else {
				row, err = db.DBC.Query("SELECT * FROM Users WHERE username = ?", existingUser.Username)
				if err != nil {
					w.Write([]byte(`{"message": "Invaild credentials"}`))
					return
				}
			}
			defer row.Close()

			var existingData ExistingUserS
			for row.Next() {
				err := row.Scan(
					&existingData.Id,
					&existingData.Username,
					&existingData.Password,
					&existingData.Email,
					&existingData.Age,
					&existingData.Gender,
					&existingData.FirstName,
					&existingData.LastName,
					&existingData.Date,
				)
				if err != nil {
					w.Write([]byte(`{"message": "Invaild credentials"}`))
					return
				}
			}

			if existingData.Id == 0 {
				w.Write([]byte(`{"message": "Invalid credentials"}`))
				return
			}

			err = row.Err()
			if err != nil {
				w.Write([]byte(`{"message": "Invaild credentials"}`))
				return
			}

			if !CheckPasswordHash(existingUser.Password, existingData.Password) {
				w.Write([]byte(`{"message": "Invalid credentials"}`))
				return
			}
			newCookie := ath.GenerateCookieInfo(existingData.Username)
			error := ath.GenerateSession(newCookie.Id, existingData.Id)
			if error != nil {
				w.Write([]byte(`{"message": "Invalid credentials"}`))
				return
			}

			var cookieData []byte
			cookieData, _ = json.Marshal(newCookie)
			w.Write(cookieData)

		} else {
			w.Write([]byte(`{"message": "Invalid credentials"}`))
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

var oneUsername UsernameS

func Username(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal([]byte(reqBody), &oneUsername)

		row, err := db.DBC.Query("SELECT username FROM Users WHERE username = ?", oneUsername.Username)

		defer row.Close()
		if row.Next() {
			w.Write([]byte(`{"value": "false"}`))
		} else {
			w.Write([]byte(`{"value": "true"}`))
		}
	}
}

var oneEmail EmailS

func Email(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal([]byte(reqBody), &oneEmail)

		row, err := db.DBC.Query("SELECT email FROM Users WHERE email = ?", oneEmail.Email)

		defer row.Close()
		//IF IT EXISTS
		if row.Next() {
			w.Write([]byte(`{"value": "false"}`))
			//IT DOESNT
		} else {
			w.Write([]byte(`{"value": "true"}`))
		}
	}
}

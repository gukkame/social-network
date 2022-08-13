package data_services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	db "real-time-forum/server/db"
	mw "real-time-forum/server/middleware"
	ath "real-time-forum/server/services/authentication"
	"real-time-forum/server/services/data/groups"
	val "real-time-forum/server/services/validation"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type NewUserS struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	NickName  string
	AboutMe   string
	Age       string
	Gender    string
	Password  string
}

type ExistingUserS struct {
	Id             int
	Username       string
	Password       string
	Email          string
	Age            string
	Gender         string
	FirstName      string
	LastName       string
	NickName       string
	AboutMe        string
	Avatar_image   string
	Profile_status string
	Date           string
}

type ProfileS struct {
	Username       string
	Email          string
	Age            string
	Gender         string
	FirstName      string
	LastName       string
	NickName       string
	AboutMe        string
	Avatar_image   string
	Date           string
	Status         string
	Profile_status string
}

type NeverReleasedData struct {
	Id             int
	Password       string
	Date           string
	Expiry_date    string
	User_id        int
	Profile_status string
}

type LogInS struct {
	Username string
	Password string
	Status   string
}

type UsernameS struct {
	Username string `json:"username"`
}

type EmailS struct {
	Email string `json:"email"`
}

var newUser NewUserS

func SignUp(w http.ResponseWriter, r *http.Request) {
	mw.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)

		err := r.ParseMultipartForm(32 << 0) // maxMemory 32MB
		if err != nil {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
		
		newUser.Username = (r.Form["username"][0])
		newUser.Email = (r.Form["email"][0])
		newUser.FirstName = (r.Form["firstname"][0])
		newUser.LastName = (r.Form["lastname"][0])
		newUser.NickName = (r.Form["nickname"][0])
		newUser.AboutMe = (r.Form["aboutme"][0])
		newUser.Age = (r.Form["age"][0])
		newUser.Gender = (r.Form["gender"][0])
		newUser.Password = (r.Form["password"][0])

		if val.ValidateUserData(newUser.Username, newUser.Email, newUser.FirstName, newUser.LastName, newUser.Age, newUser.Gender, newUser.Password, newUser.NickName, newUser.AboutMe) {
			newUser.Password, err = mw.HashPassword(newUser.Password)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			stmt, err := db.DBC.Prepare(`INSERT INTO Users(username, password, email, age, gender, first_name, last_name, nickname, about_me, avatar_image, profile_status, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime("now"))`)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			// Saving img at
			in, header, err := r.FormFile("img")
			if in == nil {
				stmt.Exec(newUser.Username, newUser.Password, newUser.Email, newUser.Age, newUser.Gender, newUser.FirstName, newUser.LastName, newUser.NickName, newUser.AboutMe, "", "private")
				defer stmt.Close()
			} else {
				if header.Header.Get("Content-Type") == "image/gif" {
					w.Write([]byte(`{"message": "Cant have gif as profile picture"}`))
					return
				}
				if header.Size <= 1048576 {

					if err != nil {
						w.Write([]byte(`{"message": "Malicious user detected"}`))
						return
					}
					defer in.Close()
					id := uuid.New()
					img_id := id.String()
					s := strings.Split(header.Filename, ".")

					out, err := os.OpenFile("./resources/profile/"+img_id+"."+s[len(s)-1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						w.Write([]byte(`{"message": "Malicious user detected"}`))
						return
					}

					defer out.Close()
					io.Copy(out, in)

					stmt.Exec(newUser.Username, newUser.Password, newUser.Email, newUser.Age, newUser.Gender, newUser.FirstName, newUser.LastName, newUser.NickName, newUser.AboutMe, "/resources/profile/"+img_id+"."+s[len(s)-1], "private")
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
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

var existingUser LogInS

func LogIn(w http.ResponseWriter, r *http.Request) {
	mw.SetupCORS(&w, r)
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
					&existingData.NickName,
					&existingData.AboutMe,
					&existingData.Avatar_image,
					&existingData.Profile_status,
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

			if !mw.CheckPasswordHash(existingUser.Password, existingData.Password) {
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
	mw.SetupCORS(&w, r)
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
		if err != nil {
			log.Fatal(err)
		}

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
	mw.SetupCORS(&w, r)
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
		if err != nil {
			log.Fatal(err)
		}
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

func AllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := DbGetAllUsers()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	groups.SendResponse(w, users)
}

func DbGetAllUsers() (users []groups.User, err error) {
	query := "SELECT Id, Username FROM Users"
	rows, err := db.DBC.Query(query)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := groups.User{}
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CheckFollowing(userIdA int, userIdB int) (isAllowed bool, err error) {
	var matches int
	query := "SELECT COUNT(status) FROM Followers WHERE follower_id = ? AND recipient_id = ?"
	row := db.DBC.QueryRow(query, userIdA, userIdB)

	err = row.Scan(&matches)
	if err != nil {
		return false, err
	}

	return matches == 1, nil
}

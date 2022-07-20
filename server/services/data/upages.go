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

type CompletePackageS3 struct {
	Package ActivityPackage
	Status  string
}

type ActivityPackage struct {
	Posts      []NeededData
	VotedPosts []NeededData
	Comments   []CommentsS
}

func Profile(w http.ResponseWriter, r *http.Request) {
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

			neededData, err := getProfile(browsed.Username)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
			neededData.Status = "owner"
			neededData.Profile_status = owner.Password

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

			neededData, err := getProfile(browsed.Username)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}

			if CheckFollower(token, browsed.Username) {
				neededData.Status = "following"
			} else {
				neededData.Status = "notfollowing"
			}

			var jsonData []byte
			jsonData, _ = json.Marshal((neededData))
			w.Write(jsonData)
			return
		}

		if profile_status.Username == "private" {

			//CHECK IF PERSON IS FOLLOWER
			if !CheckFollower(token, browsed.Username) {
				w.Write([]byte(`{"message": "Not following"}`))
				return
			}

			neededData, err := getProfile(browsed.Username)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
			neededData.Status = "following"

			var jsonData []byte
			jsonData, _ = json.Marshal((neededData))
			w.Write(jsonData)
			return
		}

	}
}

func ChangeProfileStatus(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
		row2.Close() // Not defering but instantly closing due to it otherwise locking the database for the next action.

		//CHECK IF OWNER
		var owner LogInS
		row3, _ := db.DBC.Query("SELECT username, profile_status FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)

		for row3.Next() { //I was lazy to make more structs
			err := row3.Scan(
				&owner.Username, //username
				&owner.Password, //profile_status just in case it is an owner to change the status
			)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
		}
		defer row3.Close()

		if owner.Username == browsed.Username {

			stmt, err := db.DBC.Prepare(`UPDATE Users SET profile_status = ? WHERE username = ?`)
			if err != nil {
				w.Write([]byte(`{"message": "Malicious user detected"}`))
				return
			}
			if owner.Password == "private" {
				stmt.Exec("public", browsed.Username)
				var browsedId ExistingUserS
				row8, _ := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed.Username)
				for row8.Next() {
					err := row8.Scan(
						&browsedId.Id,
					)
					if err != nil {
						w.Write([]byte(`{"message": "Malicious user detected"}`))
						return
					}
				}
				row8.Close()

				stmt3, err := db.DBC.Prepare(`UPDATE Followers SET status = ? WHERE status = ? AND recipient_id = ?`)
				if err != nil {
					w.Write([]byte(`{"message": "Malicious user detected"}`))
					return
				}
				stmt3.Exec("following", "requested", browsedId.Id)
				defer stmt3.Close()

			} else {
				stmt.Exec("private", browsed.Username)
			}
			defer stmt.Close()

			w.Write([]byte(`{"message": "Succesfully changed"}`))
			return
		} else {
			w.Write([]byte(`{"message": "Malicious user detected"}`))
			return
		}
	}

}

func getProfile(browsed string) (ProfileS, error) {
	row, err := db.DBC.Query("SELECT * FROM Users WHERE username = ?", browsed)
	if err != nil {
		return ProfileS{}, err
	}
	defer row.Close()

	var existingData NeverReleasedData
	var neededData ProfileS
	for row.Next() {
		err := row.Scan(
			&existingData.Id,
			&neededData.Username,
			&existingData.Password,
			&neededData.Email,
			&neededData.Age,
			&neededData.Gender,
			&neededData.FirstName,
			&neededData.LastName,
			&neededData.NickName,
			&neededData.AboutMe,
			&neededData.Avatar_image,
			&existingData.Profile_status,
			&neededData.Date,
		)
		if err != nil {
			return ProfileS{}, err
		}
	}
	return neededData, err
}

func Activity(w http.ResponseWriter, r *http.Request) {
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
		defer row2.Close()

		//IF IT DOESNT EXISTS
		if !row2.Next() {
			w.Write([]byte(`{"message": "Profile does not exist"}`))
			return
		}

		//CHECK IF OWNER
		var owner LogInS
		row3, _ := db.DBC.Query("SELECT username FROM Users INNER JOIN Sessions ON Sessions.user_id = Users.id WHERE token = ?", token)
		defer row3.Close()

		for row3.Next() {
			err := row3.Scan(
				&owner.Username,
			)
			if err != nil {
				w.Write([]byte(`{"message": "Profile does not exist"}`))
				return
			}
		}

		if owner.Username == browsed.Username {
			neededData, err := getActivity(browsed.Username)
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
			neededData, err := getActivity(browsed.Username)
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

			//CHECK IF PERSON IS FOLLOWER
			if !CheckFollower(token, browsed.Username) {
				w.Write([]byte(`{"message": "Not following"}`))
				return
			}

			neededData, err := getActivity(browsed.Username)
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
	}
}

func getActivity(browsed string) (CompletePackageS3, error) {
	row1, err := db.DBC.Query("SELECT id FROM Users WHERE username = ?", browsed)
	if err != nil {
		return CompletePackageS3{}, err
	}
	defer row1.Close()

	var existingUser ExistingId
	for row1.Next() {
		err := row1.Scan(
			&existingUser.Id,
		)
		if err != nil {
			return CompletePackageS3{}, err
		}
	}

	var aPackage ActivityPackage
	rows, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Categories ON Categories.id = Posts.category_id WHERE user_id = ? ORDER BY Posts.created_at DESC", existingUser.Id)
	if err != nil {
		return CompletePackageS3{}, err
	}
	defer rows.Close()

	existingData := []NeededData{}

	for rows.Next() {
		var obj NeededData
		err := rows.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Username, &obj.CategoryTitle)
		if err != nil {
			return CompletePackageS3{}, err
		}

		rows2, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC", obj.Id)
		if err != nil {
			return CompletePackageS3{}, err
		}
		defer rows2.Close()

		var obj3 []CommentsS
		for rows2.Next() {
			var obj2 CommentsS
			err2 := rows2.Scan(&obj2.Id, &obj2.Content, &obj2.Created_at, &obj2.User_id, &obj2.Post_id, &obj2.Username)
			if err2 != nil {
				return CompletePackageS3{}, err
			}
			obj3 = append(obj3, obj2)
		}

		rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", obj.Id)
		if err != nil {
			return CompletePackageS3{}, err
		}
		defer rows3.Close()

		var myVotes []VotesS
		for rows3.Next() {
			var singleVotes VotesS
			err2 := rows3.Scan(&singleVotes.Type, &singleVotes.Created_at, &singleVotes.User_id, &singleVotes.Post_id, &singleVotes.Username)
			if err2 != nil {
				return CompletePackageS3{}, err
			}
			myVotes = append(myVotes, singleVotes)
		}

		obj.Likes = append(myVotes)
		obj.Comments = append(obj3)
		existingData = append(existingData, obj)
		aPackage.Posts = append(existingData)
	}
	err = rows.Err()
	if err != nil {
		return CompletePackageS3{}, err
	}

	rows2, err := db.DBC.Query("SELECT Posts.id, Posts.title, Posts.content, Posts.created_at, Users.username, Categories.title AS category_title FROM Post_likes INNER JOIN Categories ON Categories.id = Posts.category_id INNER JOIN Users ON Users.id = Posts.user_id INNER JOIN Posts ON Posts.id = Post_likes.post_id WHERE Post_likes.user_id = ? ORDER BY Posts.created_at DESC", existingUser.Id)
	if err != nil {
		return CompletePackageS3{}, err
	}
	defer rows2.Close()

	existingData2 := []NeededData{}

	for rows2.Next() {
		var obj NeededData
		err := rows2.Scan(&obj.Id, &obj.Title, &obj.Description, &obj.Created_at, &obj.Username, &obj.CategoryTitle)
		if err != nil {
			return CompletePackageS3{}, err
		}

		rows2, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username FROM Comments INNER JOIN Users ON Users.id = Comments.user_id WHERE post_id = ? ORDER BY Comments.created_at DESC", obj.Id)
		if err != nil {
			return CompletePackageS3{}, err
		}
		defer rows2.Close()

		var obj3 []CommentsS
		for rows2.Next() {
			var obj2 CommentsS
			err2 := rows2.Scan(&obj2.Id, &obj2.Content, &obj2.Created_at, &obj2.User_id, &obj2.Post_id, &obj2.Username)
			if err2 != nil {
				return CompletePackageS3{}, err
			}
			obj3 = append(obj3, obj2)
		}

		rows3, err := db.DBC.Query("SELECT Post_likes.type, Post_likes.created_at, Post_likes.user_id, Post_likes.post_id, Users.username FROM Post_likes INNER JOIN Users ON Users.id = Post_likes.user_id WHERE post_id = ?", obj.Id)
		if err != nil {
			return CompletePackageS3{}, err
		}
		defer rows3.Close()

		var myVotes []VotesS
		for rows3.Next() {
			var singleVotes VotesS
			err2 := rows3.Scan(&singleVotes.Type, &singleVotes.Created_at, &singleVotes.User_id, &singleVotes.Post_id, &singleVotes.Username)
			if err2 != nil {
				return CompletePackageS3{}, err
			}
			myVotes = append(myVotes, singleVotes)
		}

		obj.Likes = append(myVotes)
		obj.Comments = append(obj3)
		existingData2 = append(existingData2, obj)
		aPackage.VotedPosts = append(existingData2)
	}
	err = rows.Err()
	if err != nil {
		return CompletePackageS3{}, err
	}

	rows3, err := db.DBC.Query("SELECT Comments.id, Comments.content, Comments.created_at, Comments.user_id, Comments.post_id, Users.username, Categories.title AS category_title FROM Posts INNER JOIN Users ON Users.id = Comments.user_id INNER JOIN Categories ON Categories.id = Posts.category_id INNER JOIN Comments ON Comments.post_id = Posts.id WHERE Comments.user_id = ? ORDER BY Comments.created_at DESC", existingUser.Id)
	if err != nil {
		return CompletePackageS3{}, err
	}
	defer rows3.Close()

	var obj3 []CommentsS
	for rows3.Next() {
		var obj2 CommentsS
		err2 := rows3.Scan(&obj2.Id, &obj2.Content, &obj2.Created_at, &obj2.User_id, &obj2.Post_id, &obj2.Username, &obj2.Category_title)
		if err2 != nil {
			return CompletePackageS3{}, err
		}

		rows4, err := db.DBC.Query("SELECT Comment_likes.type, Comment_likes.created_at, Comment_likes.user_id, Comment_likes.comment_id, Users.username FROM Comment_likes INNER JOIN Users ON Users.id = Comment_likes.user_id  WHERE comment_id = ?", obj2.Id)
		if err != nil {
			return CompletePackageS3{}, err
		}
		defer rows4.Close()

		var myCommentVotes []VotesS
		for rows4.Next() {
			var singleVotes VotesS
			err2 := rows4.Scan(&singleVotes.Type, &singleVotes.Created_at, &singleVotes.User_id, &singleVotes.Post_id, &singleVotes.Username)
			if err2 != nil {
				return CompletePackageS3{}, err
			}
			myCommentVotes = append(myCommentVotes, singleVotes)
		}

		obj2.Likes = append(myCommentVotes)

		obj3 = append(obj3, obj2)

	}

	aPackage.Comments = append(obj3)

	var packages CompletePackageS3

	packages.Package = aPackage
	return packages, err
}

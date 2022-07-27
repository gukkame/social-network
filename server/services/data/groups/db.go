package groups

import (
	"database/sql"
	"fmt"
	db "real-time-forum/server/db"
)

func (group Group) GetSingle() (Group, error) {
	query := "SELECT name, content, group_image, created_at, creator_id FROM User_groups WHERE id = ?"
	row := db.DBC.QueryRow(query, group.Group_id)

	err := row.Scan(&group.Title, &group.Content, &group.Image, &group.Created_at, &group.Creator_id)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (group Group) GetOwner() (Group, error) {
	query := "SELECT username FROM Users WHERE id = ?"
	row := db.DBC.QueryRow(query, group.Creator_id)

	err := row.Scan(&group.Creator_name)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (group Group) GetMembersCount() (Group, error) {
	query := "SELECT COUNT(User_id) FROM Group_users WHERE group_id = ? AND (status = ? OR status = ?)"
	row := db.DBC.QueryRow(query, group.Group_id, "Owner", "Member")

	err := row.Scan(&group.Member_count)
	if err != nil {
		return group, err
	}
	return group, nil
}

func GetUserID(token string) (userId int, err error) {
	query := "SELECT user_id FROM Sessions where token = ?"
	row := db.DBC.QueryRow(query, token)

	err = row.Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (group Group) UserRequestStatus(userId int) (status string, err error) {
	query := "SELECT status FROM Group_users WHERE group_id = ? AND user_id = ?"
	row := db.DBC.QueryRow(query, group.Group_id, userId)

	err = row.Scan(&status)
	if err == sql.ErrNoRows {
		return "Not Requested", nil
	}
	if err != nil {
		return status, err
	}
	return status, nil
}

func (member Member) Insert() (err error) {
	stmt, err := db.DBC.Prepare("INSERT INTO Group_users (user_id, group_id, status) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(member.User_id, member.Group_id, member.User_status)
	return nil
}

func (member Member) Delete() (err error) {
	stmt, err := db.DBC.Prepare("DELETE FROM Group_users WHERE user_id = ? AND group_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(member.User_id, member.Group_id)
	return nil
}

func (group Group) GetAllJoinRequests() (members []Member, err error) {
	query := "SELECT user_id, group_id, status FROM Group_users WHERE group_id = ? AND status = ?"
	rows, err := db.DBC.Query(query, group.Group_id, "Requested")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		member := Member{}
		err = rows.Scan(&member.User_id, &member.Group_id, &member.User_status)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, err
}

func GetUserById(userId int) (user User, err error) {
	query := "SELECT id, username, avatar_image FROM Users WHERE id = ?"
	row := db.DBC.QueryRow(query, userId)

	err = row.Scan(&user.Id, &user.Username, &user.Avatar_image)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (member Member) UpdateStatus(status string) (err error) {
	stmt, err := db.DBC.Prepare("UPDATE Group_users SET status = ? WHERE user_id = ? AND group_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(status, member.User_id, member.Group_id)
	return nil
}

func (notification Notification) Insert() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Notifications (type, content_id, recipient_id, user_id, created_at) VALUES (?, ?, ?, ?, datetime("now"))`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(notification.Type, notification.Content_id, notification.Recipient_id, notification.User_id, notification.Created_at)
	return nil
}

func GetUserByUsername(username string) (user User, err error) {
	query := "SELECT id, username, avatar_image FROM Users WHERE username = ?"
	row := db.DBC.QueryRow(query, username)

	err = row.Scan(&user.Id, &user.Username, &user.Avatar_image)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (group Group) GetAllPosts() (posts []Post, err error) {
	query := "SELECT id, title, content, image, created_at, group_id, user_id FROM Group_posts WHERE group_id = ? ORDER BY created_at DESC"
	rows, err := db.DBC.Query(query, group.Group_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Image, &post.Created_at, &post.Group_id, &post.User_id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, err
}

func (post Post) Insert() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Group_posts (title, content, image, created_at, group_id, user_id) VALUES (?, ?, ?, datetime("now"), ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(post.Title, post.Content, post.Image, post.Group_id, post.User_id)
	return nil
}

func (post Post) CountComments() (amount int, err error) {
	query := "SELECT COUNT(id) FROM Group_post_comments WHERE group_post_id = ?"
	row := db.DBC.QueryRow(query, post.Id)

	err = row.Scan(&amount)
	if err != nil {
		return amount, err
	}
	return amount, nil
}

func (post Post) GetAllLikes() (likes []Like, err error) {
	query := "SELECT type, created_at, user_id, group_post_id FROM Group_post_likes WHERE group_post_id = ?"
	rows, err := db.DBC.Query(query, post.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		like := Like{}
		err = rows.Scan(&like.Type, &like.Created_at, &like.User_id, &like.Post_id)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, err
}

func (post Post) GetAllComments() (comments []Comment, err error) {
	query := "SELECT id, content, image, created_at, user_id, group_post_id FROM Group_post_comments WHERE group_post_id = ? ORDER BY created_at DESC"
	rows, err := db.DBC.Query(query, post.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Image, &comment.Created_at, &comment.User_id, &comment.Post_id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, err
}

func (post *Post) GetSingle() (err error) {
	query := "SELECT id, title, content, image, created_at, group_id, user_id FROM Group_posts WHERE group_id = ? AND id = ?"
	row := db.DBC.QueryRow(query, post.Group_id, post.Id)

	err = row.Scan(&post.Id, &post.Title, &post.Content, &post.Image, &post.Created_at, &post.Group_id, &post.User_id)
	if err != nil {
		return fmt.Errorf("post.GetSingle:%s", err)
	}

	return nil
}

func (comment Comment) GetAllLikes() (likes []Like, err error) {
	query := "SELECT type, created_at, user_id, group_post_comment_id FROM Group_post_comment_likes WHERE group_post_comment_id = ?"
	rows, err := db.DBC.Query(query, comment.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		like := Like{}
		err = rows.Scan(&like.Type, &like.Created_at, &like.User_id, &like.Comment_id)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, err
}

func (like Like) InsertPostVote() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Group_post_likes (type, created_at, user_id, group_post_id) VALUES (?, datetime("now"), ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(like.Type, like.User_id, like.Post_id)
	return nil
}

func (like Like) UpdatePostVote() (err error) {
	stmt, err := db.DBC.Prepare(`UPDATE Group_post_likes SET type = ?, created_at = datetime("now") WHERE user_id = ? AND group_post_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(like.Type, like.User_id, like.Post_id)
	return nil

}

func (like *Like) GetSinglePostVote() (err error) {
	query := "SELECT type, created_at, user_id, group_post_id FROM Group_post_likes WHERE group_post_id = ? AND user_id = ?"
	row := db.DBC.QueryRow(query, like.Post_id, like.User_id)

	err = row.Scan(&like.Type, &like.Created_at, &like.User_id, &like.Post_id)
	if err != nil {
		return err
	}

	return nil
}


func (like *Like) DeletePostVote() (err error) {
	stmt, err := db.DBC.Prepare("DELETE FROM Group_post_likes WHERE user_id = ? AND group_post_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(like.User_id, like.Post_id)
	return nil
}

func (comment Comment) Insert() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Group_post_comments (content, created_at, user_id, group_post_id, image) VALUES (?, datetime("now"), ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(comment.Content, comment.User_id, comment.Post_id, comment.Image)
	return nil
}

func (like *Like) GetSingleCommentVote() (err error) {
	query := "SELECT type, created_at, user_id, group_post_comment_id FROM Group_post_comment_likes WHERE group_post_comment_id = ? AND user_id = ?"
	row := db.DBC.QueryRow(query, like.Comment_id, like.User_id)

	err = row.Scan(&like.Type, &like.Created_at, &like.User_id, &like.Comment_id)
	if err != nil {
		return err
	}

	return nil
}

func GetMostRecentPost() (post Post, err error) {
	query := "SELECT last_insert_rowid() FROM Group_posts"
	row := db.DBC.QueryRow(query)
	err = row.Scan(&post.Id)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (like Like) InsertCommentVote() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Group_post_comment_likes (type, created_at, user_id, group_post_comment_id) VALUES (?, datetime("now"), ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(like.Type, like.User_id, like.Comment_id)
	return nil
}

func (like Like) UpdateCommentVote() (err error) {
	stmt, err := db.DBC.Prepare(`UPDATE Group_post_comment_likes SET type = ?, created_at = datetime("now") WHERE user_id = ? AND group_post_comment_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(like.Type, like.User_id, like.Comment_id)
	return nil

}

func (like *Like) DeleteCommentVote() (err error) {
	stmt, err := db.DBC.Prepare("DELETE FROM Group_post_comment_likes WHERE user_id = ? AND group_post_comment_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(like.User_id, like.Comment_id)
	return nil
}

func (group Group) GetAllEvents() (events []Event, err error) {
	query := "SELECT id, title, content, happening_at, created_at, group_id, user_id FROM Group_events WHERE group_id = ?"
	rows, err := db.DBC.Query(query, group.Group_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.Id, &event.Title, &event.Content, &event.Happening_at, &event.Created_at, &event.Group_id, &event.User_id)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, err
}

func (event Event) GetAllReplies() (replies []EventReply, err error) {
	query := "SELECT status, user_id, group_events_id FROM Group_events_users WHERE group_events_id = ?"
	rows, err := db.DBC.Query(query, event.Id) 
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		reply := EventReply{}
		err = rows.Scan(&reply.Status, &reply.User_id, &reply.Event_id)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	return replies, err
}

func (event Event) Insert() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Group_events (title, content, happening_at, group_id, user_id, created_at) VALUES (?, ?, ?, ?, ?, datetime("now"))`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(event.Title, event.Content, event.Happening_at, event.Group_id, event.User_id)
	return nil
}

func (reply *EventReply) GetSingle() (err error) {
	query := "SELECT status, user_id, group_events_id FROM Group_events_users WHERE group_events_id = ? AND user_id = ?"
	row := db.DBC.QueryRow(query, reply.Event_id, reply.User_id)

	err = row.Scan(&reply.Status, &reply.User_id, &reply.Event_id)
	if err != nil {
		return err
	}

	return nil
}

func (reply EventReply) Insert() (err error) {
	stmt, err := db.DBC.Prepare(`INSERT INTO Group_events_users (status, user_id, group_events_id) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(reply.Status, reply.User_id, reply.Event_id)
	return nil
}

func (reply EventReply) Update() (err error) {
	stmt, err := db.DBC.Prepare(`UPDATE Group_events_users SET status = ? WHERE user_id = ? AND group_events_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(reply.Status, reply.User_id, reply.Event_id)
	return nil

}
func (reply EventReply) Delete() (err error) {
	stmt, err := db.DBC.Prepare("DELETE FROM Group_events_users WHERE user_id = ? AND group_events_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(reply.User_id, reply.Event_id)
	return nil
}
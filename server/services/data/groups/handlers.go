package groups

import (
	"database/sql"
	"fmt"
	"net/http"
)
// Gets Group info and the Current Users Status in the group
func GroupInfoAndUserStatus(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)

	group, err := GetGroupInfo(group)
	if err == sql.ErrNoRows {
		w.Write([]byte(`{"message":"Group does not exist"}`))
		return
	}
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	group.User_status, err = group.GetUserStatus(r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	SendResponse(w, group)
}

func MakeJoinRequest(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)
	userId := UserId(w, r)

	var member Member
	member.User_id = userId
	member.Group_id = group.Group_id
	member.User_status = "Requested"

	err := member.Insert()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CancelJoinRequest(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)
	userId := UserId(w, r)

	var member Member
	member.Group_id = group.Group_id
	member.User_id = userId

	err := member.Delete()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}
// Returns all current pending join request to a group
func ListJoinRequests(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)	
	userId := UserId(w, r)

	status, err := group.UserRequestStatus(userId)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	if status != "Owner" {
		w.Write([]byte(`{"message":"Not Owner"}`))
		return
	}

	requests, err := group.GetAllJoinRequests()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	users := []User{}
	for _, request := range requests {
		user, err := GetUserById(request.User_id)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
		users = append(users, user)
	}
	SendResponse(w, users)
}

func AcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
	member := UnmarshalMember(w, r)

	err := member.UpdateStatus("Member")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DenyJoinRequest(w http.ResponseWriter, r *http.Request) {
	member := UnmarshalMember(w, r)

	err := member.Delete()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LeaveGroup(w http.ResponseWriter, r *http.Request) {
	group := UnmarshalGroup(w, r)
	userId := UserId(w, r)

	var member Member
	member.Group_id = group.Group_id
	member.User_id = userId

	err := member.Delete()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func InviteToGroup(w http.ResponseWriter, r *http.Request) {
	invite := UnmarshalInvite(w, r)

	user, err := GetUserByUsername(invite.Username)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	var group Group
	group.Group_id = invite.GroupId
	status, err := group.UserRequestStatus(user.Id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	var member Member
	member.User_id = user.Id
	member.Group_id = invite.GroupId

	switch status {
	case "Invited":
		w.Write([]byte(`{"message":"Already Invited"}`))
		return
	case "Member":
	case "Owner":
		w.Write([]byte(`{"message":"Already Member"}`))
		return
	case "Requested":
		err = member.UpdateStatus("Invited")
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
	case "Not Requested":
		member.User_status = "Invited"
		err = member.Insert()
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func DenyGroupInvite(w http.ResponseWriter, r *http.Request) {
	member := UnmarshalMember(w, r)
	userId := UserId(w, r)

	member.User_id = userId
	err := member.Delete()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AcceptGroupInvite(w http.ResponseWriter, r *http.Request) {
	member := UnmarshalMember(w, r)
	userId := UserId(w, r)

	member.User_id = userId
	err := member.UpdateStatus("Member")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
}



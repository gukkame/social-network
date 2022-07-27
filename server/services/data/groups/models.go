package groups

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

type Post struct {
	Id              int    `json:"Id"`
	Title           string `json:"title"`
	Content         string `json:"description"`
	Image           string `json:"Image"`
	Created_at      string `json:"Created_at"`
	Group_id        int    `json:"GroupID"`
	User_id         int    `json:"User_id"`
	Comments_amount int    `json:"Comments_amount"`
	Likes           []Like `json:"Likes"`
}

type Member struct {
	User_id     int    `json:"User_id"`
	Group_id    int    `json:"Group_id"`
	User_status string `json:"User_status"`
}

type User struct {
	Id           int    `json:"Id"`
	Username     string `json:"Username"`
	Avatar_image string `json:"Avatar_image"`
}

type Notification struct {
	Type         string `json:"type"`
	Content_id   int    `json:"content_id"`
	Recipient_id int    `json:"recepient_id"`
	User_id      int    `json:"user_id"`
	Created_at   string `json:"created_at"`
}

type Like struct {
	Type       string `json:"Type"`
	Created_at string `json:"Created_at"`
	Post_id    int    `json:"Post_id"`
	Comment_id int    `json:"Comment_id"`
	User_id    int    `json:"User_id"`
	Username   string `json:"Username"`
}

type Comment struct {
	Id         int    `json:"Id"`
	Content    string `json:"Content"`
	Image      string `json:"Image"`
	Created_at string `json:"Created_at"`
	User_id    int    `json:"User_id"`
	Post_id    int    `json:"Group_id"`
	Likes      []Like `json:"Likes"`
	User       User   `json:"User"`
}

type Event struct {
	Id	int	`json:"Id"`
	Title           string `json:"Title"`
	Content    string `json:"Content"`
	Happening_at string	`json:"Happening_at"`
	Created_at string `json:"Created_at"`
	Group_id int `json:"Group_id"`
	User_id int	`json:"User_id"`
	User User	`json:"User"`
}

type EventReply struct {
	Status	string	`json:"Status"`
	User_id	int	`json:"User_id"`
	Event_id	int	`json:"Event_id"`
	Username	string `json:"Username"`
}
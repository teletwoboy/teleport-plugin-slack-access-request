package slack

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Deleted  bool   `json:"deleted"`
}

type TeamInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ReviewersChannel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsMember bool   `json:"is_member"`
}

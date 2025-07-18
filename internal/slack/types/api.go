package types

type TeamInfo struct {
	ID   string
	Name string
}

type ReviewersChannel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsMember bool   `json:"is_member"`
}

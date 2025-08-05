package types

type TeamInfo struct {
	ID   string
	Name string
}

type ReviewersChannel struct {
	ID       string
	Name     string
	IsMember bool
}

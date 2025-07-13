package teleport

type User struct {
	Username string `json:"username"`
}

type UserAccessInfo struct {
	Roles         []string `json:"roles"`
	RequireReason bool     `json:"require_reason"`
}

package util

import "strings"

func MatchesIdentifier(username, email string) bool {
	// slack email always has "@"
	id, _, _ := strings.Cut(email, "@")
	if username == email || username == id {
		return true
	}
	return false
}

package utils

import (
	"regexp"
	"strings"
)

var eventIDPattern = regexp.MustCompile(`^EVT\d{3}$`)
var userIDPattern = regexp.MustCompile(`^USR\d{3}$`)

func IsValidEventID(id string) bool {
	return eventIDPattern.MatchString(id)
}

func IsValidUserID(id string) bool {
	return userIDPattern.MatchString(id)
}

func IsValidUserPayload(name, email string) bool {
	return name != "" && strings.Contains(email, "@")
}

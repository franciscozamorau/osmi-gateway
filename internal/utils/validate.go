package utils

import "regexp"

var eventIDPattern = regexp.MustCompile(`^EVT\d{3}$`)
var userIDPattern = regexp.MustCompile(`^USR\d{3}$`)

func IsValidEventID(id string) bool {
	return eventIDPattern.MatchString(id)
}

func IsValidUserID(id string) bool {
	return userIDPattern.MatchString(id)
}

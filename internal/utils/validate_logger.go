package utils

import "log"

// Validación con logging para debugging en producción
func IsValidEventIDWithLog(eventID string) bool {
	isValid := IsValidEventID(eventID) // ✅ Usa TU función real
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid event ID format: %s", eventID)
	}
	return isValid
}

// Validación con logging para user IDs
func IsValidUserIDWithLog(userID string) bool {
	isValid := IsValidUserID(userID) // ✅ Usa TU función real
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid user ID format: %s", userID)
	}
	return isValid
}

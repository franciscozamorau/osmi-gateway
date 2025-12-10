package utils

import (
	"log"
	"strings"
)

// Validación con logging para debugging en producción
func IsValidEventIDWithLog(eventID string) bool {
	eventID = strings.TrimSpace(eventID)
	isValid := IsValidEventID(eventID)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid event ID format: %s", eventID)
	}
	return isValid
}

// Validación con logging para user IDs
func IsValidUserIDWithLog(userID string) bool {
	userID = strings.TrimSpace(userID)
	isValid := IsValidUserID(userID)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid user ID format: %s", userID)
	}
	return isValid
}

// Validación con logging para category IDs
func IsValidCategoryIDWithLog(categoryID string) bool {
	categoryID = strings.TrimSpace(categoryID)
	isValid := IsValidCategoryID(categoryID)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid category ID format: %s", categoryID)
	}
	return isValid
}

// Validación con logging para UUIDs genéricos
func IsValidUUIDWithLog(uuid string) bool {
	uuid = strings.TrimSpace(uuid)
	isValid := IsValidUUID(uuid)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid UUID format: %s", uuid)
	}
	return isValid
}

// Validación con logging para customer IDs
func IsValidCustomerIDWithLog(customerID string) bool {
	customerID = strings.TrimSpace(customerID)
	isValid := IsValidCustomerID(customerID)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid customer ID format: %s", customerID)
	}
	return isValid
}

// Validación con logging para emails
func ValidateEmailWithLog(email string) bool {
	email = strings.TrimSpace(email)
	isValid := ValidateEmail(email)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid email format: %s", email)
	}
	return isValid
}

// Validación con logging para teléfonos
func ValidatePhoneWithLog(phone string) bool {
	phone = strings.TrimSpace(phone)
	isValid := ValidatePhone(phone)
	if !isValid {
		log.Printf("VALIDATION_WARN: Invalid phone format: %s", phone)
	}
	return isValid
}

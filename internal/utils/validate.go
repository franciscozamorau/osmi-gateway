package utils

import (
	"regexp"
	"strings"
)

var (
	// Expresiones regulares compiladas
	uuidRegex  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^[\+]?[0-9\s\-\(\)]{10,}$`)
)

// IsValidUUID valida si un string es un UUID válido
func IsValidUUID(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	return uuidRegex.MatchString(strings.ToLower(s))
}

// IsValidEmail valida formato de email (alias para ValidateEmail)
func IsValidEmail(email string) bool {
	return ValidateEmail(email)
}

// IsValidPhone valida formato de teléfono (alias para ValidatePhone)
func IsValidPhone(phone string) bool {
	return ValidatePhone(phone)
}

// IsValidEventID valida si un event ID es UUID válido
func IsValidEventID(eventID string) bool {
	return IsValidUUID(eventID)
}

// IsValidUserID valida si un user ID es UUID válido
func IsValidUserID(userID string) bool {
	return IsValidUUID(userID)
}

// IsValidCategoryID valida si un category ID es UUID válido
func IsValidCategoryID(categoryID string) bool {
	return IsValidUUID(categoryID)
}

// IsValidCustomerID valida si un customer ID es UUID válido
func IsValidCustomerID(customerID string) bool {
	return IsValidUUID(customerID)
}

// ValidatePhoneE164 valida formato de teléfono E.164
func ValidatePhoneE164(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return true // Teléfono opcional
	}
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return e164Regex.MatchString(phone)
}

// ValidatePhone valida formato de teléfono más flexible
func ValidatePhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return true // Teléfono opcional
	}
	return phoneRegex.MatchString(phone)
}

// ValidateEmail valida formato de email
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false // Email requerido
	}
	return emailRegex.MatchString(email)
}

// ValidateEmailOptional valida email cuando es opcional
func ValidateEmailOptional(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return true // Email opcional
	}
	return emailRegex.MatchString(email)
}

// IsValidUserPayload valida si el payload de usuario tiene nombre y email válidos
func IsValidUserPayload(name, email string) bool {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return false
	}

	// Email puede ser opcional para usuarios
	if email != "" && !ValidateEmail(email) {
		return false
	}

	return true
}

// IsValidEventPayload valida los campos mínimos para un evento
func IsValidEventPayload(name, location string) bool {
	name = strings.TrimSpace(name)
	location = strings.TrimSpace(location)

	return name != "" && location != ""
}

// IsValidCustomerPayload valida los campos mínimos para un cliente
func IsValidCustomerPayload(name, email string) bool {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	return name != "" && ValidateEmail(email)
}

// IsValidTicketPayload valida los campos para crear un ticket
func IsValidTicketPayload(eventID, userID, categoryID string) bool {
	return IsValidUUID(eventID) && IsValidUUID(userID) && IsValidUUID(categoryID)
}

// CleanPhone limpia y formatea un número de teléfono
func CleanPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}

	// Remover todos los caracteres no numéricos excepto el +
	hasPlus := strings.HasPrefix(phone, "+")
	digits := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	if hasPlus {
		return "+" + digits
	}
	return digits
}

// NormalizeEmail normaliza un email (minúsculas, trim)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizePhone normaliza un teléfono (solo dígitos)
func NormalizePhone(phone string) string {
	return CleanPhone(phone)
}

// NormalizeName normaliza un nombre (trim, capitalize)
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Capitalizar primera letra de cada palabra
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

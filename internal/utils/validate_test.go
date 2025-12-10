package utils

import (
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", true},
		{"d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", true},
		{"invalid-uuid", false},
		{"", false},
		{"EVT123", false},
		{"USR456", false},
		{"  1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7  ", true}, // Con espacios
	}

	for _, test := range tests {
		result := IsValidUUID(test.input)
		if result != test.expected {
			t.Errorf("IsValidUUID(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidEventID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", true},
		{"d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", true},
		{"EVT123", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidEventID(test.input)
		if result != test.expected {
			t.Errorf("IsValidEventID(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidUserID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", true},
		{"d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", true},
		{"USR456", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidUserID(test.input)
		if result != test.expected {
			t.Errorf("IsValidUserID(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidCategoryID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", true},
		{"d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", true},
		{"CAT123", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidCategoryID(test.input)
		if result != test.expected {
			t.Errorf("IsValidCategoryID(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidCustomerID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1", true},
		{"123", true},
		{"0", false},
		{"-1", false},
		{"abc", false},
		{"", false},
		{"  123  ", true}, // Con espacios
	}

	for _, test := range tests {
		result := IsValidCustomerID(test.input)
		if result != test.expected {
			t.Errorf("IsValidCustomerID(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestValidatePhoneE164(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+1234567890", true},
		{"+529876543210", true},
		{"1234567890", false},
		{"+123", false},
		{"", true},                // Vacío es válido (opcional)
		{"  +1234567890  ", true}, // Con espacios
	}

	for _, test := range tests {
		result := ValidatePhoneE164(test.input)
		if result != test.expected {
			t.Errorf("ValidatePhoneE164(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+1234567890", true},
		{"+529876543210", true},
		{"1234567890", true},
		{"(555) 123-4567", true},
		{"+123", false}, // Muy corto
		{"abc", false},  // Sin dígitos
		{"", true},      // Vacío es válido (opcional)
	}

	for _, test := range tests {
		result := ValidatePhone(test.input)
		if result != test.expected {
			t.Errorf("ValidatePhone(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"", false},                    // Email requerido
		{"  test@example.com  ", true}, // Con espacios
	}

	for _, test := range tests {
		result := ValidateEmail(test.input)
		if result != test.expected {
			t.Errorf("ValidateEmail(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestValidateEmailOptional(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"", true}, // Email opcional - vacío es válido
	}

	for _, test := range tests {
		result := ValidateEmailOptional(test.input)
		if result != test.expected {
			t.Errorf("ValidateEmailOptional(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidUserPayload(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Francisco", "fran@osmi.com", true},
		{"Francisco", "invalid-email", false},
		{"Francisco", "", true}, // Email opcional
		{"", "fran@osmi.com", false},
		{"", "", false},
		{"  ", "fran@osmi.com", false},
		{"  Francisco  ", "  fran@osmi.com  ", true}, // Con espacios
	}

	for _, test := range tests {
		result := IsValidUserPayload(test.name, test.email)
		if result != test.expected {
			t.Errorf("IsValidUserPayload(%q, %q) = %v, expected %v",
				test.name, test.email, result, test.expected)
		}
	}
}

func TestIsValidEventPayload(t *testing.T) {
	tests := []struct {
		name     string
		location string
		expected bool
	}{
		{"Concierto", "Auditorio Nacional", true},
		{"Concierto", "", false},
		{"", "Auditorio Nacional", false},
		{"", "", false},
		{"  Concierto  ", "  Auditorio Nacional  ", true}, // Con espacios
	}

	for _, test := range tests {
		result := IsValidEventPayload(test.name, test.location)
		if result != test.expected {
			t.Errorf("IsValidEventPayload(%q, %q) = %v, expected %v",
				test.name, test.location, result, test.expected)
		}
	}
}

func TestIsValidCustomerPayload(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Francisco", "fran@osmi.com", true},
		{"Francisco", "invalid-email", false},
		{"Francisco", "", false}, // Email requerido para customer
		{"", "fran@osmi.com", false},
		{"", "", false},
	}

	for _, test := range tests {
		result := IsValidCustomerPayload(test.name, test.email)
		if result != test.expected {
			t.Errorf("IsValidCustomerPayload(%q, %q) = %v, expected %v",
				test.name, test.email, result, test.expected)
		}
	}
}

func TestIsValidTicketPayload(t *testing.T) {
	tests := []struct {
		eventID    string
		userID     string
		categoryID string
		expected   bool
	}{
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", "d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", "a1b2c3d4-e5f6-7890-abcd-ef1234567890", true},
		{"invalid", "d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", "a1b2c3d4-e5f6-7890-abcd-ef1234567890", false},
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", "invalid", "a1b2c3d4-e5f6-7890-abcd-ef1234567890", false},
		{"1c248f4a-f1a8-4556-9f1a-1b6a21bfadb7", "d2e8c5b1-9f4a-4c8d-b3e2-7a5f1c9b4d8a", "invalid", false},
		{"", "", "", false},
	}

	for _, test := range tests {
		result := IsValidTicketPayload(test.eventID, test.userID, test.categoryID)
		if result != test.expected {
			t.Errorf("IsValidTicketPayload(%q, %q, %q) = %v, expected %v",
				test.eventID, test.userID, test.categoryID, result, test.expected)
		}
	}
}

func TestCleanPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+1 (555) 123-4567", "+15551234567"},
		{"555-123-4567", "+5551234567"},
		{"+52 55 1234 5678", "+525512345678"},
		{"", ""},
		{"  +1 555 123 4567  ", "+15551234567"},
	}

	for _, test := range tests {
		result := CleanPhone(test.input)
		if result != test.expected {
			t.Errorf("CleanPhone(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Test@Example.COM", "test@example.com"},
		{"  USER@DOMAIN.COM  ", "user@domain.com"},
		{"", ""},
	}

	for _, test := range tests {
		result := NormalizeEmail(test.input)
		if result != test.expected {
			t.Errorf("NormalizeEmail(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"francisco zamora", "Francisco Zamora"},
		{"  john doe  ", "John Doe"},
		{"", ""},
		{"mC donald", "Mc Donald"}, // Nota: Esto podría necesitar ajustes específicos
	}

	for _, test := range tests {
		result := NormalizeName(test.input)
		if result != test.expected {
			t.Errorf("NormalizeName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

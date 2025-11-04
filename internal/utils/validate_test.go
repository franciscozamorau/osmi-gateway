package utils

import (
	"testing"
)

func TestIsValidEventID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"EVT123", true},
		{"EVT999", true},
		{"EVT12", false},
		{"EVT1234", false},
		{"ABC123", false},
		{"", false},
		{"evt123", false}, // case sensitive
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
		{"USR123", true},
		{"USR999", true},
		{"USR12", false},
		{"USR1234", false},
		{"ABC123", false},
		{"", false},
		{"usr123", false}, // case sensitive
	}

	for _, test := range tests {
		result := IsValidUserID(test.input)
		if result != test.expected {
			t.Errorf("IsValidUserID(%q) = %v, expected %v", test.input, result, test.expected)
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
		{"", "fran@osmi.com", false},
		{"", "", false},
		{"Francisco", "", false},
	}

	for _, test := range tests {
		result := IsValidUserPayload(test.name, test.email)
		if result != test.expected {
			t.Errorf("IsValidUserPayload(%q, %q) = %v, expected %v",
				test.name, test.email, result, test.expected)
		}
	}
}

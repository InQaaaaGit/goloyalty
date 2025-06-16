package service

import (
	"testing"
)

func TestIsValidOrderNumber(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "Valid order number",
			number:   "12345678903",
			expected: true,
		},
		{
			name:     "Invalid order number",
			number:   "1234567890",
			expected: false,
		},
		{
			name:     "Empty string",
			number:   "",
			expected: false,
		},
		{
			name:     "Non-numeric characters",
			number:   "1234567890a",
			expected: false,
		},
		{
			name:     "Another valid number",
			number:   "9278923470",
			expected: true,
		},
	}

	svc := &loyaltyService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.isValidOrderNumber(tt.number)
			if result != tt.expected {
				t.Errorf("isValidOrderNumber(%s) = %v, expected %v", tt.number, result, tt.expected)
			}
		})
	}
}

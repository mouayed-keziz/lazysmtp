package main

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	id := generateID()
	if len(id) != 8 {
		t.Errorf("Expected ID length 8, got %d", len(id))
	}

	id2 := generateID()
	if id == id2 {
		t.Error("Generated IDs should be unique")
	}
}

func TestExtractSubject(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "subject present",
			body:     "Subject: Test Email\nBody content",
			expected: "Test Email",
		},
		{
			name:     "subject with colon case insensitive",
			body:     "subject: Lowercase subject\nBody",
			expected: "Lowercase subject",
		},
		{
			name:     "no subject",
			body:     "Body content\nMore body",
			expected: "(no subject)",
		},
		{
			name:     "subject in middle",
			body:     "Headers\nSubject: Middle\nBody",
			expected: "Middle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSubject(tt.body)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

package main

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	dbPath := ":memory:"
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}

func TestSaveAndGetEmail(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	email := Email{
		ID:      "test123",
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
		Date:    "Mon, 01 Jan 2026 00:00:00 UTC",
	}

	err = SaveEmail(db, email)
	if err != nil {
		t.Fatalf("SaveEmail failed: %v", err)
	}

	retrieved, err := GetEmailByID(db, "test123")
	if err != nil {
		t.Fatalf("GetEmailByID failed: %v", err)
	}

	if retrieved.ID != email.ID {
		t.Errorf("Expected ID %s, got %s", email.ID, retrieved.ID)
	}
	if retrieved.From != email.From {
		t.Errorf("Expected From %s, got %s", email.From, retrieved.From)
	}
	if retrieved.To != email.To {
		t.Errorf("Expected To %s, got %s", email.To, retrieved.To)
	}
	if retrieved.Subject != email.Subject {
		t.Errorf("Expected Subject %s, got %s", email.Subject, retrieved.Subject)
	}
	if retrieved.Body != email.Body {
		t.Errorf("Expected Body %s, got %s", email.Body, retrieved.Body)
	}
}

func TestGetAllEmails(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	emails := []Email{
		{
			ID:      "1",
			From:    "a@example.com",
			To:      "b@example.com",
			Subject: "A",
			Body:    "Body A",
			Date:    "Mon, 01 Jan 2026 00:00:00 UTC",
		},
		{
			ID:      "2",
			From:    "c@example.com",
			To:      "d@example.com",
			Subject: "B",
			Body:    "Body B",
			Date:    "Mon, 02 Jan 2026 00:00:00 UTC",
		},
	}

	for _, email := range emails {
		if err := SaveEmail(db, email); err != nil {
			t.Fatalf("SaveEmail failed: %v", err)
		}
	}

	retrieved, err := GetAllEmails(db)
	if err != nil {
		t.Fatalf("GetAllEmails failed: %v", err)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 emails, got %d", len(retrieved))
	}

	found1 := false
	found2 := false
	for _, email := range retrieved {
		if email.ID == "1" {
			found1 = true
		}
		if email.ID == "2" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Error("Not all expected emails were retrieved")
	}
}

func TestDeleteEmail(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	email := Email{
		ID:      "delete_me",
		From:    "from@example.com",
		To:      "to@example.com",
		Subject: "Delete Me",
		Body:    "Body",
		Date:    "Mon, 01 Jan 2026 00:00:00 UTC",
	}

	err = SaveEmail(db, email)
	if err != nil {
		t.Fatalf("SaveEmail failed: %v", err)
	}

	err = DeleteEmail(db, "delete_me")
	if err != nil {
		t.Fatalf("DeleteEmail failed: %v", err)
	}

	_, err = GetEmailByID(db, "delete_me")
	if err == nil {
		t.Error("Expected error when getting deleted email")
	}
}

func TestDeleteAllEmails(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	email := Email{
		ID:      "test",
		From:    "from@example.com",
		To:      "to@example.com",
		Subject: "Subject",
		Body:    "Body",
		Date:    "Mon, 01 Jan 2026 00:00:00 UTC",
	}

	err = SaveEmail(db, email)
	if err != nil {
		t.Fatalf("SaveEmail failed: %v", err)
	}

	err = DeleteAllEmails(db)
	if err != nil {
		t.Fatalf("DeleteAllEmails failed: %v", err)
	}

	emails, err := GetAllEmails(db)
	if err != nil {
		t.Fatalf("GetAllEmails failed: %v", err)
	}

	if len(emails) != 0 {
		t.Errorf("Expected 0 emails after deletion, got %d", len(emails))
	}
}

func TestCountEmails(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	count, err := CountEmails(db)
	if err != nil {
		t.Fatalf("CountEmails failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected initial count 0, got %d", count)
	}

	email := Email{
		ID:      "count_test",
		From:    "from@example.com",
		To:      "to@example.com",
		Subject: "Subject",
		Body:    "Body",
		Date:    "Mon, 01 Jan 2026 00:00:00 UTC",
	}

	err = SaveEmail(db, email)
	if err != nil {
		t.Fatalf("SaveEmail failed: %v", err)
	}

	count, err = CountEmails(db)
	if err != nil {
		t.Fatalf("CountEmails failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestDatabasePersistence(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "lazysmtp_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	db1, err := InitDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	email := Email{
		ID:      "persist",
		From:    "from@example.com",
		To:      "to@example.com",
		Subject: "Persistence Test",
		Body:    "Test Body",
		Date:    "Mon, 01 Jan 2026 00:00:00 UTC",
	}

	err = SaveEmail(db1, email)
	if err != nil {
		t.Fatalf("SaveEmail failed: %v", err)
	}
	db1.Close()

	db2, err := InitDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("InitDB failed on reopen: %v", err)
	}
	defer db2.Close()

	retrieved, err := GetEmailByID(db2, "persist")
	if err != nil {
		t.Fatalf("GetEmailByID failed after reopen: %v", err)
	}

	if retrieved.Subject != "Persistence Test" {
		t.Errorf("Persistence failed: expected %s, got %s", "Persistence Test", retrieved.Subject)
	}
}

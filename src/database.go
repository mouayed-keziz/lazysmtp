package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS emails (
		id TEXT PRIMARY KEY,
		from_address TEXT NOT NULL,
		to_address TEXT NOT NULL,
		subject TEXT,
		body TEXT,
		date TEXT NOT NULL,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SaveEmail(db *sql.DB, email Email) error {
	query := `
	INSERT INTO emails (id, from_address, to_address, subject, body, date)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, email.ID, email.From, email.To, email.Subject, email.Body, email.Date)
	return err
}

func GetAllEmails(db *sql.DB) ([]Email, error) {
	query := `
	SELECT id, from_address, to_address, subject, body, date
	FROM emails
	ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		err := rows.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Date)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}

func GetEmailByID(db *sql.DB, id string) (*Email, error) {
	query := `
	SELECT id, from_address, to_address, subject, body, date
	FROM emails
	WHERE id = ?
	`
	row := db.QueryRow(query, id)
	var email Email
	err := row.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Date)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

func DeleteEmail(db *sql.DB, id string) error {
	query := `DELETE FROM emails WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func DeleteAllEmails(db *sql.DB) error {
	query := `DELETE FROM emails`
	_, err := db.Exec(query)
	return err
}

func CountEmails(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM emails`
	row := db.QueryRow(query)
	var count int
	err := row.Scan(&count)
	return count, err
}

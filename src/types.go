package main

import (
	"database/sql"
)

type AppState struct {
	SelectedEmailIndex int
	Emails             []Email
	SMTP               *SMTPServer
	DB                 *sql.DB
	NewEmailChan       chan struct{}
	Mode               string // "text" or "html"
	ShowPopup          bool
	PopupScroll        int
}

type Email struct {
	ID      string
	From    string
	To      string
	Subject string
	Body    string
	Date    string
	Headers map[string]string
}

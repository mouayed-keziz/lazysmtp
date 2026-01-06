package main

import (
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
)

type Backend struct {
	db     *sql.DB
	notify chan struct{}
}

func NewBackend(db *sql.DB, notify chan struct{}) *Backend {
	return &Backend{
		db:     db,
		notify: notify,
	}
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{db: bkd.db, notify: bkd.notify}, nil
}

type Session struct {
	db     *sql.DB
	notify chan struct{}
	from   string
	to     string
	body   strings.Builder
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	_, err := io.Copy(&s.body, r)
	if err != nil {
		return err
	}

	id := generateID()
	email := Email{
		ID:      id,
		From:    s.from,
		To:      s.to,
		Subject: extractSubject(s.body.String()),
		Body:    s.body.String(),
		Date:    time.Now().Format(time.RFC1123),
	}

	err = SaveEmail(s.db, email)
	if s.notify != nil {
		select {
		case s.notify <- struct{}{}:
		default:
		}
	}

	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = ""
	s.body.Reset()
}

func (s *Session) Logout() error {
	return nil
}

type SMTPServer struct {
	server  *smtp.Server
	port    int
	db      *sql.DB
	notify  chan struct{}
	running bool
}

func NewSMTPServer(port int, db *sql.DB, notify chan struct{}) *SMTPServer {
	return &SMTPServer{
		port:   port,
		db:     db,
		notify: notify,
	}
}

func (s *SMTPServer) Start() error {
	if s.running {
		return nil
	}

	backend := NewBackend(s.db, s.notify)

	s.server = smtp.NewServer(backend)
	s.server.Addr = fmt.Sprintf(":%d", s.port)
	s.server.Domain = "localhost"
	s.server.ReadTimeout = 10 * time.Second
	s.server.WriteTimeout = 10 * time.Second
	s.server.MaxMessageBytes = 1024 * 1024
	s.server.MaxRecipients = 50
	s.server.AllowInsecureAuth = true

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != smtp.ErrServerClosed {
		}
	}()

	s.running = true
	return nil
}

func (s *SMTPServer) Stop() {
	if s.server != nil && s.running {
		s.server.Close()
		s.running = false
	}
}

func (s *SMTPServer) IsRunning() bool {
	return s.running
}

func (s *SMTPServer) Port() int {
	return s.port
}

func (s *SMTPServer) Toggle() error {
	if s.IsRunning() {
		s.Stop()
		return nil
	}
	return s.Start()
}

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func extractSubject(body string) string {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "subject:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "(no subject)"
}

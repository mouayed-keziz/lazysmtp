package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/awesome-gocui/gocui"
)

var (
	port   = flag.Int("port", 2525, "SMTP server port")
	dbPath = flag.String("db", "", "Path to SQLite database (default: XDG data directory)")
)

func main() {
	flag.Parse()

	dbPathToUse := *dbPath
	if dbPathToUse == "" {
		dbPathToUse = GetDefaultDBPath()
	}

	db, err := InitDB(dbPathToUse)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	newEmailChan := make(chan struct{}, 100)

	state := &AppState{
		SelectedEmailIndex: -1,
		Emails:             []Email{},
		SMTP:               NewSMTPServer(*port, db, newEmailChan),
		DB:                 db,
		NewEmailChan:       newEmailChan,
		Mode:               "text",
	}

	if err := state.SMTP.Start(); err != nil {
	}

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}
	defer g.Close()

	if err := SetKeybindings(g, state); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	if err := SetLayout(g, state); err != nil {
		log.Fatalf("Failed to set layout: %v", err)
	}

	go func() {
		for range newEmailChan {
			emails, _ := GetAllEmails(db)
			state.Emails = emails
			g.Update(func(_g *gocui.Gui) error {
				updateEmailList(_g, state)
				updateServerInfo(_g, state)
				return nil
			})
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		state.SMTP.Stop()
		g.Close()
		os.Exit(0)
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("GUI error: %v", err)
	}

	state.SMTP.Stop()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func updateView(g *gocui.Gui) error {
	return nil
}

func updateEmailList(g *gocui.Gui, state *AppState) error {
	v, err := g.View("emails")
	if err != nil {
		return err
	}
	v.Clear()

	emails, err := GetAllEmails(state.DB)
	if err != nil {
		return err
	}
	state.Emails = emails

	for i, email := range emails {
		maxLen := 30
		to := email.To
		subject := email.Subject

		if len(to) > maxLen {
			to = to[:maxLen-3] + "..."
		}
		if len(subject) > maxLen {
			subject = subject[:maxLen-3] + "..."
		}

		dateStr := formatHumanDate(email.Date)

		prefix := " "
		if i == state.SelectedEmailIndex {
			prefix = ">"
			fmt.Fprintf(v, "\x1b[0;34m%s %s | %s | %s\x1b[0m\n", prefix, to, subject, dateStr)
		} else {
			fmt.Fprintf(v, "%s %s | %s | %s\n", prefix, to, subject, dateStr)
		}
	}

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func printTable(v *gocui.View, headers []string, rows [][]string, colWidths []int) {
	color := "\x1b[0;36m"
	reset := "\x1b[0m"

	top := color + "┌" + reset
	middle := color + "├" + reset
	bottom := color + "└" + reset
	for i, w := range colWidths {
		top += color + strings.Repeat("─", w+2) + reset
		middle += color + strings.Repeat("─", w+2) + reset
		bottom += color + strings.Repeat("─", w+2) + reset
		if i < len(colWidths)-1 {
			top += color + "┬" + reset
			middle += color + "┼" + reset
			bottom += color + "┴" + reset
		}
	}
	top += color + "┐" + reset
	middle += color + "┤" + reset
	bottom += color + "┘" + reset

	fmt.Fprintln(v, top)

	headerRow := color + "│" + reset
	for i, header := range headers {
		padded := fmt.Sprintf(" %-*s ", colWidths[i], header)
		headerRow += color + padded + reset + color + "│" + reset
	}
	fmt.Fprintln(v, headerRow)

	fmt.Fprintln(v, middle)

	for _, row := range rows {
		rowStr := color + "│" + reset
		for i, cell := range row {
			padded := fmt.Sprintf(" %-*s ", colWidths[i], truncateString(cell, colWidths[i]))
			rowStr += color + padded + reset + color + "│" + reset
		}
		fmt.Fprintln(v, rowStr)
	}

	fmt.Fprintln(v, bottom)
}

func formatHumanDate(dateStr string) string {
	t, err := time.Parse(time.RFC1123, dateStr)
	if err != nil {
		return dateStr
	}
	now := time.Now()
	diff := now.Sub(t)

	if diff < 24*time.Hour {
		if diff < time.Hour {
			if diff < time.Minute {
				return "just now"
			}
			return fmt.Sprintf("%dm ago", int(diff.Minutes()))
		}
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}

	return t.Format("2 Jan")
}

func updateServerInfo(g *gocui.Gui, state *AppState) error {
	v, err := g.View("server")
	if err != nil {
		return err
	}
	v.Clear()

	status := "Stopped"
	statusColor := "\x1b[0;31m"
	if state.SMTP.IsRunning() {
		status = "Running"
		statusColor = "\x1b[0;32m"
	}

	modeColor := "\x1b[0;33m"
	if state.Mode == "html" {
		modeColor = "\x1b[0;35m"
	}

	fmt.Fprintf(v, "Status: %s%s\x1b[0m\n", statusColor, status)
	fmt.Fprintf(v, "\x1b[0;36mPort:\x1b[0m %d\n", state.SMTP.Port())
	fmt.Fprintf(v, "\x1b[0;36mEmails:\x1b[0m %d\n", len(state.Emails))
	fmt.Fprintf(v, "Mode: %s%s\x1b[0m\n", modeColor, state.Mode)
	fmt.Fprintf(v, "\n\x1b[0;33m[SPACE]\x1b[0m Toggle Server")
	fmt.Fprintf(v, "\n\x1b[0;33m[m]\x1b[0m Toggle Mode")

	return nil
}

func updateMainView(g *gocui.Gui, state *AppState) error {
	v, err := g.View("main")
	if err != nil {
		return err
	}
	v.Clear()

	if state.SelectedEmailIndex >= 0 && state.SelectedEmailIndex < len(state.Emails) {
		email := state.Emails[state.SelectedEmailIndex]
		fmt.Fprintf(v, "\x1b[1;36mEmail Details:\x1b[0m\n")

		emailRows := [][]string{
			{"ID", email.ID},
			{"From", email.From},
			{"To", email.To},
			{"Subject", email.Subject},
			{"Date", email.Date},
		}

		if contentType := email.Headers["Content-Type"]; contentType != "" {
			emailRows = append(emailRows, []string{"Content-Type", contentType})
		}

		printTable(v, []string{"Field", "Value"}, emailRows, []int{15, 38})

		var bodyContent string
		if state.Mode == "text" {
			bodyContent = htmlToText(email.Body)
		} else {
			bodyContent = email.Body
		}

		fmt.Fprintf(v, "\n\x1b[1;36mBody (%s mode):\x1b[0m\n%s\n", state.Mode, bodyContent)
	} else {
		fmt.Fprint(v, GetColoredASCIIArt())
		fmt.Fprintf(v, "\n\n\x1b[1;33mFeatures:\x1b[0m\n")
		fmt.Fprintf(v, "\x1b[0;36m•\x1b[0m Lightweight SMTP server for testing\n")
		fmt.Fprintf(v, "\x1b[0;36m•\x1b[0m Capture and view emails in real-time\n")
		fmt.Fprintf(v, "\x1b[0;36m•\x1b[0m HTML and text mode support\n")
		fmt.Fprintf(v, "\x1b[0;36m•\x1b[0m SQLite database for persistence\n")
		fmt.Fprintf(v, "\x1b[0;36m•\x1b[0m Keyboard-driven TUI interface\n")
		fmt.Fprintf(v, "\x1b[0;36m•\x1b[0m Perfect for development and testing\n")

		fmt.Fprintf(v, "\n\n\x1b[1;33mControls:\x1b[0m\n")
		printTable(v, []string{"Key", "Action"}, [][]string{
			{"j/k", "Navigate emails"},
			{"ESC", "Go back to home"},
			{"d", "Delete selected email"},
			{"SPACE", "Toggle server"},
			{"m", "Toggle text/html"},
			{"q / Ctrl+C", "Quit application"},
		}, []int{15, 26})
	}

	return nil
}

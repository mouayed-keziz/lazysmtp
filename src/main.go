package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awesome-gocui/gocui"
)

var (
	port     = flag.Int("port", 2525, "SMTP server port")
	dbPath   = flag.String("db", "", "Path to SQLite database (default: XDG data directory)")
	asciiArt = `
 ___      _______  _______  __   __  _______  __   __  _______  _______ 
|   |    |   _   ||       ||  | |  ||       ||  |_|  ||       ||       |
|   |    |  |_|  ||____   ||  |_|  ||  _____||       ||_     _||    _  |
|   |    |       | ____|  ||       || |_____ |       |  |   |  |   |_| |
|   |___ |       || ______||_     _||_____  ||       |  |   |  |    ___|
|       ||   _   || |_____   |   |   _____| || ||_|| |  |   |  |   |    
|_______||__| |__||_______|  |___|  |_______||_|   |_|  |___|  |___|    
                                        
   SMTP Testing Tool for Developers
`
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

	if err := SetKeybindings(g, state, asciiArt); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	if err := SetLayout(g, state, asciiArt); err != nil {
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

func updateMainView(g *gocui.Gui, state *AppState, art string) error {
	v, err := g.View("main")
	if err != nil {
		return err
	}
	v.Clear()

	if state.SelectedEmailIndex >= 0 && state.SelectedEmailIndex < len(state.Emails) {
		email := state.Emails[state.SelectedEmailIndex]
		fmt.Fprintf(v, "\x1b[1;36mID:\x1b[0m %s\n", email.ID)
		fmt.Fprintf(v, "\x1b[1;36mFrom:\x1b[0m %s\n", email.From)
		fmt.Fprintf(v, "\x1b[1;36mTo:\x1b[0m %s\n", email.To)
		fmt.Fprintf(v, "\x1b[1;36mSubject:\x1b[0m %s\n", email.Subject)
		fmt.Fprintf(v, "\x1b[1;36mDate:\x1b[0m %s\n", email.Date)

		// Show content type if available
		if contentType := email.Headers["Content-Type"]; contentType != "" {
			fmt.Fprintf(v, "\x1b[1;36mContent-Type:\x1b[0m %s\n", contentType)
		}

		var bodyContent string
		if state.Mode == "text" {
			bodyContent = htmlToText(email.Body)
		} else {
			bodyContent = email.Body
		}

		fmt.Fprintf(v, "\n\x1b[1;36mBody (%s mode):\x1b[0m\n%s\n", state.Mode, bodyContent)
	} else {
		fmt.Fprint(v, art)
		fmt.Fprintf(v, "\n\n\x1b[1;33mControls:\x1b[0m\n")
		fmt.Fprintf(v, "\x1b[0;36mj/k\x1b[0m - Navigate emails\n")
		fmt.Fprintf(v, "\x1b[0;36mESC\x1b[0m - Go back to home\n")
		fmt.Fprintf(v, "\x1b[0;36md\x1b[0m - Delete email\n")
		fmt.Fprintf(v, "\x1b[0;36mSPACE\x1b[0m - Toggle server\n")
		fmt.Fprintf(v, "\x1b[0;36mm\x1b[0m - Toggle text/html mode\n")
		fmt.Fprintf(v, "\x1b[0;36mq\x1b[0m / \x1b[0;36mCtrl+C\x1b[0m - Quit\n")
	}

	return nil
}

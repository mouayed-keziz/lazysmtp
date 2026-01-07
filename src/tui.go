package main

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func SetLayout(g *gocui.Gui, state *AppState) error {
	maxX, maxY := g.Size()

	leftPanelWidth := 50

	if v, err := g.SetView("server", 0, 0, leftPanelWidth, maxY/3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "lazySMTP Config"
		v.Wrap = false
		v.FrameColor = gocui.ColorCyan
		v.TitleColor = gocui.ColorCyan
	}

	if v, err := g.SetView("emails", 0, maxY/3, leftPanelWidth, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Emails"
		v.Wrap = false
		v.FrameColor = gocui.ColorCyan
		v.TitleColor = gocui.ColorCyan
	}

	if v, err := g.SetView("main", leftPanelWidth+1, 0, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "lazySMTP"
		v.Wrap = true
		v.FrameColor = gocui.ColorGreen
		v.TitleColor = gocui.ColorGreen
	}

	popupWidth := 60
	popupHeight := 20
	if v, err := g.SetView("popup", (maxX-popupWidth)/2, (maxY-popupHeight)/2, (maxX+popupWidth)/2, (maxY+popupHeight)/2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keybindings"
		v.Wrap = true
		v.FrameColor = gocui.ColorYellow
		v.TitleColor = gocui.ColorYellow
		v.Editable = false
		v.Highlight = false
		v.SelBgColor = gocui.ColorDefault
		v.SelFgColor = gocui.ColorDefault
	}

	if state.ShowPopup {
		if err := updatePopupView(g, state); err != nil {
			return err
		}
		_, err := g.SetCurrentView("popup")
		if err != nil {
			return err
		}
	} else {
		if popupView, err := g.View("popup"); err == nil && popupView != nil {
			if err := g.DeleteView("popup"); err != nil && err != gocui.ErrUnknownView {
				return err
			}
		}
		_, err := g.SetCurrentView("main")
		if err != nil {
			return err
		}
	}

	if err := updateServerInfo(g, state); err != nil {
		return err
	}
	if err := updateEmailList(g, state); err != nil {
		return err
	}
	if err := updateMainView(g, state); err != nil {
		return err
	}

	return nil
}

func SetKeybindings(g *gocui.Gui, state *AppState) error {
	if err := g.SetKeybinding("", 'q', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.ShowPopup {
			state.ShowPopup = false
			state.PopupScroll = 0
			if err := SetLayout(gui, state); err != nil {
				return err
			}
			return nil
		}
		return quit(gui, v)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.ShowPopup {
			state.ShowPopup = false
			state.PopupScroll = 0
			if err := SetLayout(gui, state); err != nil {
				return err
			}
			return nil
		}
		return quit(gui, v)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.ShowPopup {
			state.ShowPopup = false
			state.PopupScroll = 0
			if err := SetLayout(gui, state); err != nil {
				return err
			}
			return nil
		}
		state.SelectedEmailIndex = -1
		if err := updateEmailList(gui, state); err != nil {
			return err
		}
		if err := updateMainView(gui, state); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'j', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		emails, _ := GetAllEmails(state.DB)
		if len(emails) > 0 && state.SelectedEmailIndex < len(emails)-1 {
			state.SelectedEmailIndex++
			if err := updateEmailList(gui, state); err != nil {
				return err
			}
			if err := updateMainView(gui, state); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'k', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.SelectedEmailIndex > 0 {
			state.SelectedEmailIndex--
			if err := updateEmailList(gui, state); err != nil {
				return err
			}
			if err := updateMainView(gui, state); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'd', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.SelectedEmailIndex >= 0 && state.SelectedEmailIndex < len(state.Emails) {
			email := state.Emails[state.SelectedEmailIndex]
			if err := DeleteEmail(state.DB, email.ID); err != nil {
				return err
			}

			emails, _ := GetAllEmails(state.DB)
			state.Emails = emails
			if state.SelectedEmailIndex >= len(emails) {
				state.SelectedEmailIndex = len(emails) - 1
			}

			if err := updateEmailList(gui, state); err != nil {
				return err
			}
			if err := updateMainView(gui, state); err != nil {
				return err
			}
			if err := updateServerInfo(gui, state); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if err := state.SMTP.Toggle(); err != nil {
			return err
		}
		if err := updateServerInfo(gui, state); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'm', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.Mode == "text" {
			state.Mode = "html"
		} else {
			state.Mode = "text"
		}
		if err := updateServerInfo(gui, state); err != nil {
			return err
		}
		if err := updateMainView(gui, state); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'x', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.ShowPopup {
			state.ShowPopup = false
			state.PopupScroll = 0
			if err := SetLayout(gui, state); err != nil {
				return err
			}
			return nil
		}
		state.ShowPopup = true
		state.PopupScroll = 0
		if err := SetLayout(gui, state); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("popup", 'j', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		keybindings := []struct {
			key     string
			action  string
			context string
		}{
			{"x", "Open keybindings popup", "Always"},
			{"j/k", "Navigate emails down/up", "Email list"},
			{"j/k", "Scroll popup down/up", "Popup"},
			{"ESC", "Go back to home / Close popup", "When viewing email / Popup"},
			{"d", "Delete selected email", "Email list"},
			{"SPACE", "Toggle SMTP server on/off", "Always"},
			{"m", "Toggle text/html mode", "Always"},
			{"q", "Quit / Close popup", "Always / Popup"},
			{"Ctrl+C", "Quit / Close popup", "Always / Popup"},
		}

		popupView, err := gui.View("popup")
		if err != nil {
			return err
		}
		_, maxY := popupView.Size()
		visibleLines := maxY - 6
		maxScroll := len(keybindings) - visibleLines

		if maxScroll > 0 && state.PopupScroll < maxScroll {
			state.PopupScroll++
		}
		if err := updatePopupView(gui, state); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("popup", 'k', gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		if state.PopupScroll > 0 {
			state.PopupScroll--
		}
		if err := updatePopupView(gui, state); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func updatePopupView(g *gocui.Gui, state *AppState) error {
	v, err := g.View("popup")
	if err != nil {
		return err
	}
	v.Clear()

	keybindings := []struct {
		key    string
		action string
	}{
		{"x", "Open / Close keybindings popup"},
		{"j/k", "Navigate emails down/up"},
		{"j/k", "Scroll popup down/up"},
		{"ESC", "Go back to home / Close popup"},
		{"d", "Delete selected email"},
		{"SPACE", "Toggle SMTP server on/off"},
		{"m", "Toggle text/html mode"},
		{"q", "Quit application / Close popup"},
		{"Ctrl+C", "Quit application / Close popup"},
	}

	_, maxY := v.Size()
	visibleLines := maxY - 4

	for i := state.PopupScroll; i < len(keybindings) && i-state.PopupScroll < visibleLines; i++ {
		kb := keybindings[i]
		fmt.Fprintf(v, "\x1b[0;33m%12s\x1b[0m  %s\n", kb.key, kb.action)
	}

	if len(keybindings) > visibleLines {
		fmt.Fprintf(v, "\n\x1b[0;90mShowing %d-%d of %d keybindings\x1b[0m", state.PopupScroll+1, min(state.PopupScroll+visibleLines, len(keybindings)), len(keybindings))
	}

	return nil
}

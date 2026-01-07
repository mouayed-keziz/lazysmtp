package main

import (
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
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
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

	return nil
}

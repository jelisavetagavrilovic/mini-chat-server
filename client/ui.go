package main

import (
	"net"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var commands = []string{"/users", "/quit", "/help"}

func NewChatUI(app *tview.Application, conn net.Conn, activeUsers *[]string) (*tview.TextView, *tview.InputField) {
	// create the message view
	messageView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)
		
	messageView.SetBorder(true).SetTitle(fmt.Sprintf("Chat - %s", myName))
	messageView.SetChangedFunc(func() { 
		messageView.ScrollToEnd()
		app.Draw() 
	})

	// create the input field
	input := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorBlack). 
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldWidth(0)

	messageView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyTAB:
			app.SetFocus(input)
			return nil

		// case (event.Key() == tcell.KeyUp) && (event.Modifiers()&(tcell.ModMeta|tcell.ModCtrl)) != 0:
		// 	messageView.ScrollTo(0, 0)
		// 	return nil

		// case (event.Key() == tcell.KeyDown) && (event.Modifiers()&(tcell.ModMeta|tcell.ModCtrl)) != 0:
		// 	messageView.ScrollToEnd()
		// 	return nil

		case event.Key() == tcell.KeyUp:
			row, _ := messageView.GetScrollOffset()
			messageView.ScrollTo(row-1, 0)
			return nil

		case event.Key() == tcell.KeyDown:
			row, _ := messageView.GetScrollOffset()
			messageView.ScrollTo(row+1, 0)
			return nil
		}

		return event
	})

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(messageView)
			return nil
		}
		return event
	})

	// autocomplete functions
	input.SetAutocompleteFunc(func(currentText string) []string {
		currentText = strings.TrimSpace(currentText)
		if currentText == "" {
			return nil
		}

		suggestions := []string{}

		if strings.HasPrefix(currentText, "@") {
			for _, u := range *activeUsers {
				if strings.HasPrefix(u, strings.TrimPrefix(currentText, "@")) {
					suggestions = append(suggestions, "@" + u)
				}
			}
		}

		if strings.HasPrefix(currentText, "/") {
			for _, c := range commands {
				if strings.HasPrefix(c, currentText) {
					suggestions = append(suggestions, c)
				}
			}
		}

		return suggestions
	})

	// handle commands and messages
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		text := input.GetText()
		if text == "" {
			return
		}

		if text == "/quit" {
			SendMessage(conn, text)
			app.Stop()
			return
		}

		if text == "/users" {
			SendMessage(conn, text)
			AppendSystemMessage(messageView, "Active users: " + strings.Join(*activeUsers, ", "))
			input.SetText("")
			return
		}

		if text == "/help" {
			AppendSystemMessage(messageView, `Available commands:
			@user msg - send a private message msg to user
			/users    - show list of active users
			/quit     - leave the chat
			/help     - show this help message
			`)
			input.SetText("")
			return
		}

		SendMessage(conn, text)
		AppendMessage(messageView, text, true, strings.HasPrefix(text, "@"))
		input.SetText("")
	})

	return messageView, input
}

package main

import (
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var commands = []string{"/quit", "/users"}

func NewChatUI(app *tview.Application, conn net.Conn, activeUsers *[]string) (*tview.TextView, *tview.InputField) {
	// create the message view
	messageView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() { app.Draw() })
	messageView.SetBorder(true).SetTitle("Chat")

	// create the input field
	input := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)

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
			AppendSystemMessage(messageView, "Active users: " + strings.Join(*activeUsers, ", "))
			input.SetText("")
			return
		}

		SendMessage(conn, text)
		AppendMessage(messageView, text, true, strings.HasPrefix(text, "@"))
		input.SetText("")
	})

	return messageView, input
}

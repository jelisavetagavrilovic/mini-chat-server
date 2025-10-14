package main

import (
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Kreira TextView i InputField za chat
func NewChatUI(app *tview.Application, conn net.Conn) (*tview.TextView, *tview.InputField) {
	messageView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() { app.Draw() })
	messageView.SetBorder(true).SetTitle("Chat")

	input := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := input.GetText()
			if text != "" {
				SendMessage(conn, text)
				AppendMessage(messageView, text, true, false)
				input.SetText("")
			}
		}
	})

	return messageView, input
}

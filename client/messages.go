package main

import (
	"fmt"
	// "strings"
	"time"

	"github.com/rivo/tview"
)

// AppendSystemMessage adds a system message to the chat view
func AppendSystemMessage(view *tview.TextView, msg string) {
	timestamp := time.Now().Format("15:04")
	fmt.Fprintf(view, "[%s] %s[-]\n", timestamp, msg)
}

// AppendMessage adds a message to the chat view
func AppendMessage(view *tview.TextView, msg string, isMe bool, isPrivate bool) {
	timestamp := time.Now().Format("15:04")

	switch {
	case isMe:
		fmt.Fprintf(view, "[#61AFEF][%s] You: %s[-]\n", timestamp, msg)
	case isPrivate:
		fmt.Fprintf(view, "[#C678DD][%s] (Private) %s[-]\n", timestamp, msg[10:])
	default:
		fmt.Fprintf(view, "[green][%s] %s[-]\n", timestamp, msg)
	}
}
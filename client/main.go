package main

import (
	// "fmt"
	"log"
	"net"

	"github.com/rivo/tview"
)

func main() {
	// connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	// ensure connection closes when program exits
	defer conn.Close()

    name := askName(conn)
	setMyName(name)

	activeUsers := []string{} // local list of active users
	app := tview.NewApplication() // create a new TUI application

	// create message view and input field
	messageView, input := NewChatUI(app, conn, &activeUsers)
	// layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messageView, 0, 1, false).
		AddItem(input, 1, 0, true)

	// start goroutine to handle incoming messages
	go clientReader(conn, app, messageView, input, &activeUsers)

	// run the TUI application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal(err)
	}
}

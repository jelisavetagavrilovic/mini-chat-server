package main

import (
	// "fmt"
	"log"
	"net"

	"github.com/rivo/tview"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

    name := askName(conn)
	setMyName(name)

	app := tview.NewApplication()
	messageView, input := NewChatUI(app, conn)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messageView, 0, 1, false).
		AddItem(input, 1, 0, true)

	go clientReader(conn, messageView)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal(err)
	}
}

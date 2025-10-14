package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/rivo/tview"
)

var myName string

func setMyName(name string) {
	myName = name
}

func askName(conn net.Conn) string {
	fmt.Print("Enter your name: ")
	var name string
	fmt.Scanln(&name)
	conn.Write([]byte(name + "\n"))
	return name
}

func clientReader(conn net.Conn, view *tview.TextView) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(view, "[red]Disconnected from server.[-]\n")
			return
		}
		msg = strings.TrimSpace(msg)

		if strings.HasSuffix(msg, "has joined the chat") || strings.HasSuffix(msg, "has left the chat") || strings.HasSuffix(msg, "User not found") {
			AppendSystemMessage(view, msg)
			continue
		}

		sender := parseSender(msg)

		if sender == myName {
			continue
		}

		if strings.HasPrefix(msg, "[Private]") {
			AppendMessage(view, msg, false, true)
		} else {
			AppendMessage(view, msg, false, false)
		}
	}
}


func SendMessage(conn net.Conn, msg string) {
	conn.Write([]byte(msg + "\n"))
}

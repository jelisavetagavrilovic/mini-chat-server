package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/rivo/tview"
)

// client username
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

// clientReader reads messages from the server and updates the UI
func clientReader(conn net.Conn, view *tview.TextView, activeUsers *[]string) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(view, "[red]Disconnected from server.[-]\n")
			return
		}
		msg = strings.TrimSpace(msg)

		// existing users - list of users already logged in
        if strings.HasPrefix(msg, "Active users: ") {
            users := strings.TrimPrefix(msg, "Active users: ")
            *activeUsers = strings.Split(users, ", ")
            continue
        }

		// system messages - user joined
		if strings.HasSuffix(msg, "has joined the chat") {
			username := strings.TrimSuffix(msg, " has joined the chat")
			*activeUsers = appendUser(*activeUsers, username)
			AppendSystemMessage(view, msg)
			continue
		}

		// system messages - user left
		if strings.HasSuffix(msg, "has left the chat") {
			username := strings.TrimSuffix(msg, " has left the chat")
			*activeUsers = removeUser(*activeUsers, username)
			AppendSystemMessage(view, msg)
			continue
		}

		sender := parseSender(msg)

		if sender == myName {
			continue 
		}

		// display messages
		if strings.HasPrefix(msg, "[Private]") {
			AppendMessage(view, msg, false, true)
		} else {
			AppendMessage(view, msg, false, false)
		}
	}
}

// appendUser adds a user to the active list if not already present
func appendUser(list []string, user string) []string {
	for _, u := range list {
		if u == user {
			return list
		}
	}
	return append(list, user)
}

// removeUser removes a user from the active list
func removeUser(list []string, user string) []string {
	newList := []string{}
	for _, u := range list {
		if u != user {
			newList = append(newList, u)
		}
	}
	return newList
}

// SendMessage sends a message to the server
func SendMessage(conn net.Conn, msg string) {
	conn.Write([]byte(msg + "\n"))
}

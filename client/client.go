package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/rivo/tview"
)

var (
	myName string
	reader *bufio.Reader
)

func setMyName(name string) {
	myName = name
}

func askName(conn net.Conn) string {
	var name string
	reader = bufio.NewReader(conn)
	
	for {
		name = ""
		fmt.Print("Enter your username: ")
		fmt.Scanln(&name)
		conn.Write([]byte(name + "\n"))

		// read server response
        response, _ := reader.ReadString('\n')
        response = strings.TrimSpace(response)

        if response == "NAME_ACCEPTED" {
            return name
        } else if response == "NAME_TAKEN" {
			fmt.Println("Username is already taken, please choose another one.")
		} else if response == "NOT_NAME" {
			fmt.Println("Please choose a username.")
		} else {
			fmt.Println("Unexpected server response.")
		}
	}
}

// clientReader reads messages from the server and updates the UI
func clientReader(conn net.Conn, app *tview.Application, view *tview.TextView, input *tview.InputField, activeUsers *[]string) {
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            app.QueueUpdateDraw(func() {
                fmt.Fprintf(view, "[red]Disconnected from server.[-]\n")
                input.SetDisabled(true)
            })
            return
        }

        // handle case where multiple lines arrive in one read
        lines := strings.Split(msg, "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line == "" {
                continue
            }

            app.QueueUpdateDraw(func() {
                // existing users - list of users already logged in
                if strings.HasPrefix(line, "Active users: ") {
                    users := strings.TrimPrefix(line, "Active users: ")
                    *activeUsers = strings.Split(users, ", ")
                    return
                }

                // system messages - user joined
                if strings.HasSuffix(line, "has joined the chat") {
                    username := strings.TrimSuffix(line, " has joined the chat")
                    *activeUsers = appendUser(*activeUsers, username)
                    AppendSystemMessage(view, line)
                    return
                }

                // system messages - user left
                if strings.HasSuffix(line, "has left the chat") {
                    username := strings.TrimSuffix(line, " has left the chat")
                    *activeUsers = removeUser(*activeUsers, username)
                    AppendSystemMessage(view, line)
                    return
                }

				if strings.HasSuffix(line, "User not found") {
					AppendSystemMessage(view, line)
					return
				}

                sender := parseSender(line)
                if sender == myName {
                    return
                }

				// display messages
                if strings.HasPrefix(line, "[Private]") {
                    AppendMessage(view, line, false, true)
                } else {
                    AppendMessage(view, line, false, false)
                }
            })
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

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
                colonIdx := strings.Index(line, ":")
                if colonIdx == -1 {
                    AppendSystemMessage(view, "Unexpected server response.")
                    return
                }

                from := strings.TrimSpace(line[:colonIdx])
                content := strings.TrimSpace(line[colonIdx+1:])

                // system messages
                if from == "System" {
                    // existing users - list of users already logged in
                    if strings.HasPrefix(content, "Active users: ") {
                        users := strings.TrimPrefix(content, "Active users: ")
                        *activeUsers = strings.Split(users, ", ")
                        return
                    }

                    // user joined
                    if strings.HasSuffix(content, "has joined the chat") {
                        username := strings.TrimSuffix(content, " has joined the chat")
                        *activeUsers = appendUser(*activeUsers, username)
                        AppendSystemMessage(view, content)
                        return
                    }

                    // user left
                    if strings.HasSuffix(content, "has left the chat") {
                        username := strings.TrimSuffix(content, " has left the chat")
                        *activeUsers = removeUser(*activeUsers, username)
                        AppendSystemMessage(view, content)
                        return
                    }

                    // user not found
                    if strings.HasSuffix(content, "User not found") {
                        AppendSystemMessage(view, content)
                        return
                    }

                }

                if from == myName {
                    return
                }

				// display messages
                if strings.HasPrefix(content, "[Private]") {
                    AppendMessage(view, content, false, true)
                } else {
                    AppendMessage(view, content, false, false)
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

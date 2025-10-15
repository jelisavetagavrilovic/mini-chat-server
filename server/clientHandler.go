package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

type Client struct {
    Name string
    Conn net.Conn
}

func handleClient(conn net.Conn) {
    // close connection when function exits
    defer conn.Close()

    nameReader := bufio.NewReader(conn)
    name, _ := nameReader.ReadString('\n')
    name = strings.TrimSpace(name)


    clientsMux.Lock()

    // collect existing users to send them to the new user
    var existingUsers []string
    for uname := range clients {
        existingUsers = append(existingUsers, uname)
    }

    // add a new client to clients map
    clients[name] = Client{Name: name, Conn: conn}
    clientsMux.Unlock()

    if len(existingUsers) > 0 {
        // send the list to the new user
        conn.Write([]byte("Active users: " + strings.Join(existingUsers, ", ") + "\n"))
    }

    messages <- fmt.Sprintf("%s has joined the chat", name)

    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        msg = strings.TrimSpace(msg)

        // command - "quit"
        if strings.HasPrefix(msg, "/quit") {
            conn.Write([]byte("Goodbye!\n"))

            // update clients map
            clientsMux.Lock()
            delete(clients, name)
            clientsMux.Unlock()

            messages <- fmt.Sprintf("%s has left the chat", name)
            return // defer - close the connection
        }

        // private message if starts with "@"
        if strings.HasPrefix(msg, "@") {
            parts := strings.SplitN(msg, " ", 2)

            if len(parts) == 2 {
                targetName := strings.TrimPrefix(parts[0], "@")
                clientsMux.Lock()
                if target, ok := clients[targetName]; ok {
                	target.Conn.Write([]byte(fmt.Sprintf("[Private] %s: %s\n", name, parts[1])))
                } else {
                    conn.Write([]byte("User not found\n"))
                }
                clientsMux.Unlock()
                continue
            }
        }

        // broadcast
        messages <- fmt.Sprintf("%s: %s", name, msg)
    }

    // remove client from map if he disconnect (cmd + c)
    clientsMux.Lock()
    delete(clients, name)
    clientsMux.Unlock()
    
    messages <- fmt.Sprintf("%s has left the chat", name)
}
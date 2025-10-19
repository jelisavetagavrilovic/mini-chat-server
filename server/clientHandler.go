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
    var name string

    // username selection loop
    for {
        n, err := nameReader.ReadString('\n')
        if err != nil {
            return
        }
        name = strings.TrimSpace(n)

        if name == "" {
            conn.Write([]byte("NOT_NAME\n"))
            continue
        }

        clientsMux.Lock()
        _, exists := clients[name]
        clientsMux.Unlock()

        if exists {
            conn.Write([]byte("NAME_TAKEN\n"))
            continue
        }

        conn.Write([]byte("NAME_ACCEPTED\n"))
        break
    }

    clientsMux.Lock()

    // add a new client to clients map
    clients[name] = Client{Name: name, Conn: conn}

    // collect current users
    var existingUsers []string
    for uname := range clients {
        existingUsers = append(existingUsers, uname)
    }
    clientsMux.Unlock()

    SendMessage("System", "", fmt.Sprintf("%s has joined the chat", name))

    // send the list to the new user
    SendMessage("System", name, "Active users: " + strings.Join(existingUsers, ", "))
    

    // start reading user messages
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        msg = strings.TrimSpace(msg)

        // handle quit command 
        if strings.HasPrefix(msg, "/quit") {
            clientsMux.Lock()
            delete(clients, name)
            clientsMux.Unlock()
            SendMessage("System", "", fmt.Sprintf("%s has left the chat", name))
            return
        }

        // handle user list command 
        if strings.HasPrefix(msg, "/users") {
            var existingUsers []string
            clientsMux.Lock()
            for uname := range clients {
                existingUsers = append(existingUsers, uname)
            }
            clientsMux.Unlock()

            SendMessage("System", name, "Active users: " + strings.Join(existingUsers, ", "))
            continue
        }

        // handle private messages
        if strings.HasPrefix(msg, "@") {
            parts := strings.SplitN(msg, " ", 2)
            if len(parts) == 2 {
                targetName := strings.TrimPrefix(parts[0], "@")
                clientsMux.Lock()
                if target, ok := clients[targetName]; ok {
                    SendMessage(name, target.Name, fmt.Sprintf("[Private] %s: %s\n", name, parts[1]))
                } else {
                    SendMessage("System", name, "User not found")
                }
                clientsMux.Unlock()
                continue
            }
        }

        // broadcast
        SendMessage(name, "", fmt.Sprintf("%s: %s\n", name, msg))
    }

    // remove client from map if he disconnect (cmd + c)
    clientsMux.Lock()
    delete(clients, name)
    clientsMux.Unlock()
    SendMessage("System", "", fmt.Sprintf("%s has left the chat", name))
    return
}
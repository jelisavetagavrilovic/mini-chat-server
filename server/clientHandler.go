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
    defer conn.Close()
    nameReader := bufio.NewReader(conn)
    name, _ := nameReader.ReadString('\n')
    name = strings.TrimSpace(name)

    clientsMux.Lock()
    clients[name] = Client{Name: name, Conn: conn}
    clientsMux.Unlock()

    messages <- fmt.Sprintf("%s has joined the chat", name)

    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        msg = strings.TrimSpace(msg)

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

    clientsMux.Lock()
    delete(clients, name)
    clientsMux.Unlock()
    
    messages <- fmt.Sprintf("%s has left the chat", name)
}

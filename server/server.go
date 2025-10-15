package main

import "sync"

// global variables
var (
    clients    = make(map[string]Client) // map to store active clients
    clientsMux sync.Mutex                // mutex to protect clients map
    messages   = make(chan string)       // channel for messages to broadcast
)

// broadcaster listens on messages channel and sends message to all connected clients
func broadcaster() {
    for {
        msg := <-messages
        clientsMux.Lock()
        for _, client := range clients {
            client.Conn.Write([]byte(msg + "\n"))
        }
        clientsMux.Unlock()
    }
}

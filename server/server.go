package main

import (
    "sync"
    "fmt"
)

type Message struct {
    From    string
    To      string 
    Content string
}

// global variables
var (
    clients    = make(map[string]Client) // map to store active clients
    clientsMux sync.Mutex                // mutex to protect clients map
    messages   = make(chan Message)      // channel for messages to broadcast
)

// dispatcher listens on messages channel and sends message 
func dispatcher() {
    for {
        msg := <- messages
        clientsMux.Lock()
        if msg.To == "" {
            // broadcast 
            for _, c := range clients {
                c.Conn.Write([]byte(fmt.Sprintf("%s: %s\n", msg.From, msg.Content)))
            }
        } else {
            // private or system message
            if target, ok := clients[msg.To]; ok {
                target.Conn.Write([]byte(fmt.Sprintf("%s: %s\n", msg.From, msg.Content)))
            }
        }
        clientsMux.Unlock()
    }
}

// helper function for sending message in channel
func SendMessage(from, to, content string) {
    messages <- Message{
        From:    from,
        To:      to,
        Content: content,
    }
}

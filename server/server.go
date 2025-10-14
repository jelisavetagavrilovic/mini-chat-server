package main

import "sync"

var (
    clients    = make(map[string]Client)
    clientsMux sync.Mutex
    messages   = make(chan string)
)

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

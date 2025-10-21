package main

import (
    "fmt"
    "net"
)

func main() {
    // start listening on TCP port 8080
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    // ensure listener is closed when main exits
    defer listener.Close()

    fmt.Println("Server started on port 8080")

    // starts a goroutine that continuously listens on the 'messages' channel 
    // and forwards messages to the appropriate clients 
	go dispatcher()

    // accept incoming connections
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        // handle each client in a separate goroutine
        go handleClient(conn)
    }
}

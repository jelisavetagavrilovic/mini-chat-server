package main

import (
    "bufio"
	"fmt"
    "net"
    "os"
    "strings"
)

func readMessages(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("\nDisconnected from server.")
            os.Exit(0)
        }

        fmt.Print(msg) 
    }
}

func readInput(conn net.Conn) {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        if !scanner.Scan() {
            break
        }

        text := scanner.Text()
        text = strings.TrimSpace(text)
        if text == "" {
            continue
        }
		
        conn.Write([]byte(text + "\n"))
    }
}
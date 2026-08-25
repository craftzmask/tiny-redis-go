package main

import (
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
)

type RESPType struct {
}

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	// Handle multi-clients, so the code enters an infinite loop
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		// Instantly spawn handleClient in the background, so it won't block
		// the current thread and go back to l.Accept() immediately
		// This makes the server concurrent (multi-threaded)
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)
	for {
		_, err := conn.Read(buff)
		if err != nil {
			break
		}

		resp := "+PONG\r\n"
		str := strings.ToLower(string(buff))
		tokens := strings.Split(str, "\r\n")
		if slices.Contains(tokens, "echo") {
			content := tokens[len(tokens)-2]
			resp = fmt.Sprintf("$%d\r\n%s\r\n", len(content), content)
		}

		_, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Println("Error writing to connection: ", err)
			return
		}
	}
}

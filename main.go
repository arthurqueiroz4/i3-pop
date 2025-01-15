package main

import (
	"github.com/golang-collections/collections/stack"
	"go.i3wm.org/i3/v4"
	"log"
	"net"
	"strings"
	"sync"
)

// TODO: Go ahead with double stack
// TODO: Is it better Unix or TCP socket?
// TODO: Think other way for ignore events sent by myself

var x int

func main() {

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	msg := strings.TrimSpace(string(buf[:n]))
	log.Printf("Message received: %s", msg)

	if msg == "back" {
		mutex.Lock()
		defer mutex.Unlock()

		if s.Len() == 0 {
			log.Println("Stack is empty, no workspace to go back to")
			return
		}

		workspaceName := s.Pop().(string)
		log.Printf("Switching to workspace: %s", workspaceName)
		x++
		if _, err := i3.RunCommand("workspace " + workspaceName); err != nil {
			log.Printf("Failed to switch workspace: %v", err)
		}
	} else {
		log.Printf("Invalid command: %s", msg)
	}
}

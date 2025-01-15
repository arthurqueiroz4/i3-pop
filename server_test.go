package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
)

func TestServerTCP(t *testing.T) {
	msgToSend := "ping"

	serverReady := make(chan struct{})
	var wg sync.WaitGroup

	h := func(conn net.Conn) {
		defer conn.Close()
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			return
		}

		msg := strings.TrimSpace(string(buf[:n]))
		assert.NotEmpty(t, msg)
		assert.Equal(t, msgToSend, msg)

		_, err = conn.Write([]byte("pong"))
		assert.Nil(t, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		sr := NewServer("tcp", ":42333")
		defer sr.Stop()
		sr.Start(h, serverReady)
	}()

	<-serverReady

	conn, err := net.Dial("tcp", ":42333")
	assert.Nil(t, err)
	defer conn.Close()

	_, err = conn.Write([]byte(msgToSend))
	assert.Nil(t, err)

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	assert.Nil(t, err)
	response := strings.TrimSpace(string(buf[:n]))
	assert.Equal(t, "pong", response)
}

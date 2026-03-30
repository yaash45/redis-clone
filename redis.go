package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/yaash45/redis/internal/command"
	"github.com/yaash45/redis/internal/store"
)

var PORT = ":6379"

func main() {

	// Initialize new in-memory key-value store
	kvs := store.NewStore()

	// Set-up a tcp listener
	listener, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("Could not set up listener: %s", err.Error())
	}

	log.Printf("Listening on port %s...", PORT)
	defer listener.Close()

	// Process commands until server is terminated
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting con:", err)
		}

		_, err = conn.Write([]byte("Waiting to connect with server...\n\n"))

		if err != nil {
			log.Println("[warning] server write error")
		}

		handleConnection(conn, kvs)
	}
}

// Handles a tcp connection and facilitates client interaction with the key-value store
func handleConnection(conn net.Conn, kvs *store.KVStore) {

	log.Printf("Accepted new connection: %s\n", conn.RemoteAddr().String())

	_, err := conn.Write([]byte("Connected. Type commands or 'exit'.\n"))

	if err != nil {
		log.Println("server write error")
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		_, err := conn.Write([]byte("> "))

		if err != nil {
			log.Println("server write error")
			return
		}

		message, err := reader.ReadString('\n')

		if err != nil {
			log.Printf("request message read error: %s", err.Error())
			return
		}

		trimmedMessage := strings.TrimSpace(message)

		if strings.ToLower(trimmedMessage) == "exit" {
			break
		}

		cmd, err := command.Parse(message)

		if err != nil {
			log.Printf("bad input error: %s", err.Error())
		}

		result, err := kvs.ProcessCmd(&cmd)

		var response string

		if err != nil {
			response = fmt.Sprintf("%s\n\n", err.Error())
		} else {
			response = fmt.Sprintf("%s\n\n", result)
		}

		_, err = conn.Write([]byte(response))

		if err != nil {
			log.Println("server write error:", err)
		}
	}

	log.Printf("Ending connection with %s", conn.RemoteAddr().String())

}

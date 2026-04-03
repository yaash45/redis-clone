package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

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

	// Casting into a tcp listener to allow setting a connection accept deadline
	tcpListener := listener.(*net.TCPListener)

	log.Printf("Listening on port %s...", PORT)

	defer tcpListener.Close()

	// Keep track of the client connections and track write error counts
	var clientsAndErrorCounts map[net.Conn]int = make(map[net.Conn]int)

	// Process commands until server is terminated
	for {
		// Accept new connections, but set a deadline before moving on
		// to check for any other clients trying to connect
		tcpListener.SetDeadline(time.Now().Add(time.Millisecond * 10))
		conn, err := tcpListener.Accept()

		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				// Simply go ahead and server the clients that are
				// already connected to the server
			} else {
				// Not a timeout error, so something went wrong
				log.Fatalf("Error accepting con: %s", err)
			}
		}

		// If a timeout happens, the connection is nil, so check before
		// adding it to our list of tracked clients
		if conn != nil {
			val, ok := clientsAndErrorCounts[conn]

			if ok {
				// Client connection is already being tracked, nothing to do here
			} else {
				clientsAndErrorCounts[conn] = 0
				// Acknowledge the connection to the client
				_, err := conn.Write([]byte("Connected. Type commands or 'exit'.\n\n> "))

				if err != nil {
					// Something went wrong
					if val > 3 {
						log.Printf("max attepts to write to %s exhausted. closing connection", conn.RemoteAddr().String())
						conn.Close()
						delete(clientsAndErrorCounts, conn)
					} else {
						log.Println("server write error")
						clientsAndErrorCounts[conn] += 1
					}
					continue
				}
			}
		}

		var code int

		// Serve all clients one-by-one
		for c := range clientsAndErrorCounts {
			code = handleConnection(c, kvs)

			switch code {
			case 0:
				// Simple timeout, keep going
			case 1:
				// Exit was issued, close connection normally
				c.Write([]byte("Connection closed.\n"))
				c.Close()
				delete(clientsAndErrorCounts, c)
			default:
				// Some error occured
				c.Write([]byte("Server error.\n"))
				c.Close()
				delete(clientsAndErrorCounts, c)
			}
		}
	}
}

// Handles a tcp connection and facilitates client interaction with the key-value store
//
// It returns:
//   - 0 if there is nothing to read, and not close the connection
//   - 1 benign exit, so just close the connection
//   - -1 if there is an error
func handleConnection(conn net.Conn, kvs *store.KVStore) int {
	reader := bufio.NewReader(conn)

	// set a 500 millisecond read deadline
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

	message, err := reader.ReadString('\n')

	if err != nil {

		var netErr net.Error

		// Detect read timeouts or Ctrl+C keyboard interrupts and treat
		// them as benign conditions
		if err == io.EOF {
			// Keyboard interrupt (ctrl+c) was isssued, close connection with client
			return 1
		} else if errors.As(err, &netErr) && netErr.Timeout() {
			// Nothing to read, proceed normally and relinquish control of the thread
			return 0
		} else {
			// Something truly went wrong, return an error code
			log.Printf("[error] Reading client message failed: %s", err.Error())
			return -1
		}
	}

	log.Printf("Responding to: %s\n", conn.RemoteAddr().String())

	// Clean up the message and check if exit is requested
	trimmedMessage := strings.TrimSpace(message)

	if strings.ToLower(trimmedMessage) == "exit" {
		return 1
	}

	// Parse the message into a Command
	cmd, err := command.Parse(message)

	if err != nil {
		log.Printf("bad input error: %s", err.Error())
		return -1
	}

	// Process the parsed command
	result, err := kvs.ProcessCmd(&cmd)

	var response string

	// Prepare the response string and write the result/error to the client
	if err != nil {
		response = fmt.Sprintf("%s\n\n> ", err.Error())
	} else {
		response = fmt.Sprintf("%s\n\n> ", result)
	}

	_, err = conn.Write([]byte(response))

	if err != nil {
		log.Println("server write error:", err)
		return -1
	}

	return 0
}

package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/yaash45/redis/internal/client"
	"github.com/yaash45/redis/internal/command"
	"github.com/yaash45/redis/internal/status"
	"github.com/yaash45/redis/internal/store"
)

var PORT = ":6379"
var TCP_ACCEPTANCE_DURATION = time.Millisecond * 50

func main() {

	// Initialize new in-memory key-value store
	kvs := store.NewStore()

	// Initialize a registry to track all the clients connected to the server
	clientRegistry := client.NewRegistry()

	// Set-up a tcp listener
	listener, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("Could not set up listener: %s", err.Error())
	}

	// Casting into a tcp listener to allow setting a connection accept deadline
	tcpListener := listener.(*net.TCPListener)

	log.Printf("Listening on port %s...", PORT)

	defer tcpListener.Close()

	// Process commands until server is terminated
	for {
		// Accept new connections, but set a deadline before moving on
		// to check for any other clients trying to connect
		tcpListener.SetDeadline(time.Now().Add(TCP_ACCEPTANCE_DURATION))
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
			c := clientRegistry.Register(conn)

			// Acknowledge the connection to the client
			writeResult := c.Write([]byte("Connected. Type commands or 'exit'.\n\n> "))

			if writeResult.Status() == status.FatalErr {
				// Something went wrong, remove client
				clientRegistry.Remove(c)
			}
		}

		var clientsToRemove []*client.Client

		// Serve all clients one-by-one
		for _, c := range clientRegistry.Clients {
			code := handleConnection(c, kvs)

			switch code {
			case status.Success:
				// Succesful operation, keep going
				log.Printf("[success] remote: %s", c.RemoteAddr())

			case status.Timeout:
				// Timeout, keep serving other clients

			case status.Close:
				// Exit was issued, close connection normally
				log.Printf("[close] remote: %s", c.RemoteAddr())

				if res := c.Write([]byte("Connection closed.\n\n")); res.Status() != status.Success {
					log.Println("write failed during close: ", res.Error())
				}
				clientsToRemove = append(clientsToRemove, c)

			case status.BadRequestErr:
				// Non-fatal bad client request
				log.Printf("[bad request error] remote: %s", c.RemoteAddr())

			case status.ServerErr:
				// Some non-fatal server error occurred
				log.Printf("[server error] remote: %s", c.RemoteAddr())

			case status.FatalErr:
				// Some fatal error occurred, close connection
				log.Printf("[fatal error] remote: %s", c.RemoteAddr())

				if res := c.Write([]byte("Fatal error. Closing connection.\n\n")); res.Status() != status.Success {
					log.Println("write failed during close: ", res.Error())
				}
				clientsToRemove = append(clientsToRemove, c)
			}
		}

		// Remove all the clients that were marked for removal
		for _, c := range clientsToRemove {
			clientRegistry.Remove(c)
		}
	}
}

// Handles a tcp connection and facilitates client interaction with the key-value store
func handleConnection(c *client.Client, kvs *store.KVStore) status.StatusCode {

	res := c.Read('\n')

	switch res.Status() {
	// Keyboard interrupt (ctrl+c) was isssued, close connection with client
	case status.Close:
		return status.Close

	// Nothing to read, proceed normally and relinquish control of the thread
	case status.Timeout:
		return status.Timeout

	// Something truly went wrong, return an error code
	case status.ServerErr:
		return status.ServerErr
	}

	trimmedMessage := strings.TrimSpace(res.Message())

	if strings.ToLower(trimmedMessage) == "exit" {
		// Client explicitly sent the exit command
		return status.Close
	}

	// Parse the message into a Command
	cmd, err := command.Parse(trimmedMessage)

	if err != nil {
		if res := c.Write(fmt.Appendf(nil, "Bad input error: %s\n\n> ", err.Error())); res.Status() != status.Success {
			log.Println("Server write error: ", res.Error())
			return status.FatalErr
		}
		return status.BadRequestErr
	}

	// Process the parsed command
	kvResult, err := kvs.ProcessCmd(&cmd)

	// Prepare the response string and write the result/error to the client
	if err != nil {
		if res := c.Write(fmt.Appendf(nil, "%s\n\n> ", err.Error())); res.Status() != status.Success {
			log.Println("Server write error: ", res.Error())
			return status.FatalErr
		}

		return status.BadRequestErr
	}

	response := fmt.Sprintf("%s\n\n> ", kvResult)

	res = c.Write([]byte(response))

	if res.Status() == status.ServerErr {
		log.Println("Server write error:", res.Error())
		return status.FatalErr
	}

	return status.Success
}

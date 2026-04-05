package client

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"

	"github.com/yaash45/redis/internal/result"
	"github.com/yaash45/redis/internal/status"
)

var READ_TIMEOUT = time.Millisecond * 50

// Simple wrapper to abstract a Client connection
type Client struct {
	conn net.Conn
	rbuf *bufio.Reader
}

// Create a new Client instance
func NewClient(conn net.Conn) *Client {
	newBuffer := bufio.NewReader(conn)

	return &Client{
		conn: conn,
		rbuf: newBuffer,
	}
}

// Read a slice of bytes from the client before a read
// deadline specified by the `READ_TIMEOUT` global variable
func (client *Client) Read(delim byte) result.Result {

	client.conn.SetReadDeadline(time.Now().Add(READ_TIMEOUT))

	msg, err := client.rbuf.ReadBytes(delim)

	if err != nil {

		var netErr net.Error

		if err == io.EOF {
			return result.NewResult().WithStatus(status.Close)
		} else if errors.As(err, &netErr) && netErr.Timeout() {
			return result.NewResult().WithStatus(status.Timeout)
		} else {
			return result.NewResult().WithMessage(msg).WithStatus(status.ServerErr).WithErr(err)
		}
	}

	return result.NewResult().WithMessage(msg).WithStatus(status.Success)
}

// Write the given bytes to the Client
func (client *Client) Write(b []byte) result.Result {

	_, err := client.conn.Write(b)

	if err != nil {
		return result.NewResult().WithStatus(status.ServerErr).WithErr(err)
	}

	return result.NewResult().WithStatus(status.Success)
}

// Get the remote address of the Client
func (client *Client) RemoteAddr() string {
	return client.conn.RemoteAddr().String()
}

// Close the connection with the Client
func (client *Client) Close() {
	client.conn.Close()
}

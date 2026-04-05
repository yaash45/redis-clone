package client

import "net"

// Tracks all the clients registered with the system
type ClientRegistry struct {
	Clients map[string]*Client
}

// Get a new registry instance
func NewRegistry() *ClientRegistry {
	return &ClientRegistry{
		Clients: make(map[string]*Client),
	}
}

// Register a new client, or return existing client from the
// same remote address
func (registry *ClientRegistry) Register(conn net.Conn) *Client {

	c, ok := registry.Clients[conn.RemoteAddr().String()]

	if !ok {
		c = NewClient(conn)
		registry.Clients[c.RemoteAddr()] = c
	}

	return c
}

// Remove client from registry
func (registry *ClientRegistry) Remove(c *Client) {
	c.Close()
	delete(registry.Clients, c.RemoteAddr())
}

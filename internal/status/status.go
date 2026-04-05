package status

// A simple code indicating the status of an operation
type StatusCode int

// The possible set of codes and what they mean
const (
	// Successful request
	Success StatusCode = iota

	// Client hit a timeout deadline
	Timeout

	// Client closed the connection
	Close

	// Client error occurred
	BadRequestErr

	// Some error occurred
	ServerErr

	// Some fatal error occurred
	FatalErr
)

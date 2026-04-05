package result

import "github.com/yaash45/redis/internal/status"

// Simple wrapper to pass around values, status, and errors
type Result struct {
	message []byte
	status  status.StatusCode
	err     error
}

// Build a base Result wrapper object
func NewResult() Result {
	return Result{}
}

// Build new Result with a given message
func (res Result) WithMessage(message []byte) Result {
	return Result{
		message: message,
		status:  res.status,
		err:     res.err,
	}
}

// Build new Result with a given status
func (res Result) WithStatus(s status.StatusCode) Result {
	return Result{
		message: res.message,
		status:  s,
		err:     res.err,
	}
}

// Build new Result with a given error
func (res Result) WithErr(err error) Result {
	return Result{
		message: res.message,
		status:  res.status,
		err:     err,
	}
}

// Get the message of the Result
func (res Result) Message() string {
	return string(res.message)
}

// Get the status of the Result
func (res Result) Status() status.StatusCode {
	return res.status
}

// Get the error string of the Result
func (res Result) Error() string {
	if res.err != nil {
		return res.err.Error()
	}
	return ""
}

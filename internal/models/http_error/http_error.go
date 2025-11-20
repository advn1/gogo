package httperror

import "fmt"

// Custom Http Error
type HttpError struct {
	Message string
	Code    int
}

// Formatting Http Error
func (e *HttpError) Error() string {
	return fmt.Sprintf("%s Code: %d", e.Message, e.Code)
}
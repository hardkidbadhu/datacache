package Errs

import "fmt"

type AppErr struct {
	Message string `json:"message"`
	Err     string `json:"error,omitempty"`
}

func (ae *AppErr) Error() string {
	return fmt.Sprintf("error: %s - (%s)", ae.Message, ae.Err)
}

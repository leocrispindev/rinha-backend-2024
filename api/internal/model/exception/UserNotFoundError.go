package exception

import "fmt"

type UserNotFound struct {
	Message string
}

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("Transaction error: ", e.Message)
}

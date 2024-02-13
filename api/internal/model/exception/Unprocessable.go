package exception

import "fmt"

type UnprocessableEntity struct {
	Status  int
	Message string
}

func (e *UnprocessableEntity) Error() string {
	return fmt.Sprintf("Transaction error: ", e.Message)
}

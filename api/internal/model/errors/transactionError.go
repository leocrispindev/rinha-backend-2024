package errors

import "fmt"

type TransactionError struct {
	Message string
}

func (e *TransactionError) Error() string {
	return fmt.Sprintf("Transaction error: ", e.Message)
}

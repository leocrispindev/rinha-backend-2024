package model

import (
	"errors"
	"time"
)

type Transaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em,omitempty"`
}

func (t *Transaction) Validate() error {
	if t.Type != "c" && t.Type != "d" {
		return errors.New("invalid transaction [type]")
	}

	if t.Description == "" || len(t.Description) > 10 {
		return errors.New("invalid transaction [description]")
	}

	return nil
}

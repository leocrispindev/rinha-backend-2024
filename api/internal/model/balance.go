package model

import "time"

type Balance struct {
	Total int       `json:"total"`
	Limit int       `json:"limite"`
	Date  time.Time `json:"data_extrato, omitempty"`
}

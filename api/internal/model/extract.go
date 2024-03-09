package model

import "time"

type Extract struct {
	Saldo        ExtractBalance `json:"saldo"`
	Transactions []Transaction  `json:"ultimas_transacoes"`
}

type ExtractBalance struct {
	Total        int       `json:"total"`
	BalanceLimit int       `json:"limite"`
	Date         time.Time `json:"data_extrato, omitempty"`
}

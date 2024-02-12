package model

type Extract struct {
	Saldo        Balance       `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

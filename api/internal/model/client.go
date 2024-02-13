package model

type Client struct {
	Id           int `json:"-"`
	Balance      int `json:"limite"`
	BalanceLimit int `json:"saldo"`
}

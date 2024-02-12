package model

type Client struct {
	Id           int `json:"id,omitempty"`
	Balance      int `json:"balance"`
	BalanceLimit int `json:"balanceLimit"`
}

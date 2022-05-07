package model

type Wallet struct {
	ID      int `json:"-"`
	UserID  int `json:"user_id"`
	Balance int `json:"balance"`
}

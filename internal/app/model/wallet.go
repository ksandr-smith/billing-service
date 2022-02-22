package model

type Wallet struct {
	ID      int `json:"id,omitempty"`
	UserID  int `json:"user_id"`
	Balance int `json:"balance"`
}

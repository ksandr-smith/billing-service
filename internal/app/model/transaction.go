package model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	UUID      uuid.UUID `json:"uuid"`
	WalletID  int       `json:"wallet_id"`
	Amount    int       `json:"amount"`
	Desc      string    `json:"transaction_type"`
	Timestamp time.Time `json:"created"`
}

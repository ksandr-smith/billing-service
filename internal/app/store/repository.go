package store

import (
	"billing-service/internal/app/model"
	"net/url"
)

type WalletRepository interface {
	Create(*model.Wallet) error
	FindByUserId(int) (*model.Wallet, error)
	GetOrCreate(*model.Wallet) (bool, error)
	UpdateBalance(*model.Wallet, int) error
}

type TransactionRepository interface {
	Create(*model.Transaction) error
	FindByWalletId(int, url.Values) ([]model.Transaction, error)
}

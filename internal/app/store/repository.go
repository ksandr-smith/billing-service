package store

import (
	"billing-service/internal/app/model"
	"context"
	"net/url"
)

type WalletRepository interface {
	Create(ctx context.Context, w *model.Wallet) error
	FindByUserId(ctx context.Context, id int) (*model.Wallet, error)
	GetOrCreate(ctx context.Context, w *model.Wallet) (bool, error)
	UpdateBalance(ctx context.Context, w *model.Wallet, amount int) error
}

type TransactionRepository interface {
	Create(ctx context.Context, t *model.Transaction) error
	FindByWalletId(ctx context.Context, id int, values url.Values) ([]model.Transaction, error)
}

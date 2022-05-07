package sqlstore

import (
	"billing-service/internal/app/model"
	"context"
	"errors"
)

type WalletRepository struct {
	store *Store
}

func (r *WalletRepository) Create(ctx context.Context, w *model.Wallet) error {
	return r.store.db.QueryRow(ctx, "INSERT INTO wallet (user_id) VALUES ($1) RETURNING id, balance", w.UserID).Scan(&w.ID, &w.Balance)
}

func (r *WalletRepository) FindByUserId(ctx context.Context, userId int) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID: userId,
	}

	if _, err := r.GetOrCreate(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *WalletRepository) GetOrCreate(ctx context.Context, w *model.Wallet) (bool, error) {
	if err := r.store.db.QueryRow(ctx, "SELECT id, balance FROM wallet WHERE user_id = $1", w.UserID).Scan(&w.ID, &w.Balance); err != nil {
		err2 := r.Create(ctx, w)
		if err2 != nil {
			return false, err2
		}
		return true, nil
	}
	return false, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, w *model.Wallet, amount int) error {
	if amount < 0 && w.Balance < -amount {
		return errors.New("not enough funds")
	}
	w.Balance += amount
	if err := r.store.db.QueryRow(ctx, "UPDATE wallet SET balance = $1 WHERE user_id = $2 RETURNING balance", w.Balance, w.UserID).Scan(&w.Balance); err != nil {
		return nil
	}
	return nil
}

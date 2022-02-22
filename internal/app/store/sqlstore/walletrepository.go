package sqlstore

import (
	"billing-service/internal/app/model"
	"errors"
)

type WalletRepository struct {
	store *Store
}

func (r *WalletRepository) Create(w *model.Wallet) error {
	return r.store.db.QueryRow("INSERT INTO wallet (user_id) VALUES ($1) RETURNING id, balance", w.UserID).Scan(&w.ID, &w.Balance)
}

func (r *WalletRepository) FindByUserId(userId int) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID: userId,
	}

	if _, err := r.GetOrCreate(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *WalletRepository) GetOrCreate(w *model.Wallet) (bool, error) {
	if err := r.store.db.QueryRow("SELECT id, balance FROM wallet WHERE user_id = $1", w.UserID).Scan(&w.ID, &w.Balance); err != nil {
		err2 := r.Create(w)
		if err2 != nil {
			return false, err2
		}
		return true, nil
	}
	return false, nil
}

func (r *WalletRepository) UpdateBalance(w *model.Wallet, amount int) error {
	if amount < 0 && w.Balance < -amount {
		return errors.New("not enough funds")
	}
	w.Balance += amount
	if err := r.store.db.QueryRow("UPDATE wallet SET balance = $1 WHERE user_id = $2", w.Balance, w.UserID).Err(); err != nil {
		return err
	}
	return nil
}

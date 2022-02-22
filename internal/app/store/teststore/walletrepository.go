package teststore

import (
	"billing-service/internal/app/model"
	"errors"
)

type WalletRepository struct {
	store   *Store
	wallets map[int]*model.Wallet
}

func (r *WalletRepository) Create(w *model.Wallet) error {

	r.wallets[w.UserID] = w
	w.ID = len(r.wallets)

	return nil
}

func (r *WalletRepository) FindByUserId(userId int) (*model.Wallet, error) {
	w, ok := r.wallets[userId]
	if !ok {
		return nil, errors.New("not found")
	}
	return w, nil
}

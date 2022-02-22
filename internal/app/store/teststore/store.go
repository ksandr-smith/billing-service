package teststore

import (
	"billing-service/internal/app/model"
	"billing-service/internal/app/store"
)

type Store struct {
	walletRepository *WalletRepository
}

func New() *Store {
	return nil
}

func (s *Store) Wallet() store.WalletRepository {
	if s.walletRepository != nil {
		return s.walletRepository
	}

	s.walletRepository = &WalletRepository{
		store:   s,
		wallets: make(map[int]*model.Wallet),
	}

	return s.walletRepository
}

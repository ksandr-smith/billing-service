package sqlstore

import (
	"billing-service/internal/app/store"
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db                    *sql.DB
	walletRepository      *WalletRepository
	transactionRepository *TransactionRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Wallet() store.WalletRepository {
	if s.walletRepository != nil {
		return s.walletRepository
	}

	s.walletRepository = &WalletRepository{
		store: s,
	}

	return s.walletRepository
}

func (s *Store) Transaction() store.TransactionRepository {
	if s.transactionRepository != nil {
		return s.transactionRepository
	}

	s.transactionRepository = &TransactionRepository{
		store: s,
	}

	return s.transactionRepository
}

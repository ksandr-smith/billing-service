package sqlstore

import (
	"billing-service/internal/app/model"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type TransactionRepository struct {
	store *Store
}

func (r *TransactionRepository) Create(ctx context.Context, t *model.Transaction) error {
	return r.store.db.QueryRow(ctx,
		"INSERT INTO transactions (wallet_id, transaction_type, amount) VALUES ($1, $2, $3) RETURNING uuid, created",
		t.WalletID,
		t.Desc,
		t.Amount).Scan(&t.UUID, &t.Timestamp)
}

func (r *TransactionRepository) FindByWalletId(ctx context.Context, walletId int, query url.Values) ([]model.Transaction, error) {
	var transactions []model.Transaction
	sql := fmt.Sprintf("SELECT * FROM transactions WHERE wallet_id = %d", walletId)

	if sort := query.Get("sort"); sort != "" {
		direction := "ASC"
		sorts := strings.Split(sort, ",")
		next := false
		for _, sort = range sorts {
			if strings.Contains(sort, "-") {
				sort = strings.Replace(sort, "-", "", 1)
				direction = "DESC"
			}
			if !next {
				sql = fmt.Sprintf("%s ORDER BY %s %s", sql, sort, direction)
				next = true
			} else {
				sql = fmt.Sprintf("%s, %s %s", sql, sort, direction)
			}
		}
	}

	var currentPage int
	var perPage int

	if page := query.Get("page"); page == "" {
		currentPage = 1
	} else {
		currentPage, _ = strconv.Atoi(page)
	}

	if offset := query.Get("offset"); offset != "" {
		perPage, _ = strconv.Atoi(offset)
	}

	if perPage != 0 {
		sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perPage, (currentPage-1)*perPage)
	}

	rows, err := r.store.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.UUID, &t.WalletID, &t.Amount, &t.Desc, &t.Timestamp); err != nil {
			return transactions, err
		}
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return transactions, err
	}
	return transactions, nil
}

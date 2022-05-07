package apiserver

import (
	"billing-service/internal/app/store/sqlstore"
	postgresql "billing-service/pkg/client"
	"context"
	"net/http"
)

func Start(config *Config) error {
	postgres, err := postgresql.NewClient(context.TODO(), 3, *config.Store)
	if err != nil {
		return err
	}

	store := sqlstore.New(postgres)

	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}

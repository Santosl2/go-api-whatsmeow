package services

import (
	"context"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

var sqlStoreInstance *sqlstore.Container

func SQLStore() *sqlstore.Container {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if sqlStoreInstance == nil {
		container, err := sqlstore.New(ctx, "sqlite3", "file:whatsmeow.db?_foreign_keys=on", nil)

		if err != nil {
			panic(err)
		}

		sqlStoreInstance = container
	}

	return sqlStoreInstance
}

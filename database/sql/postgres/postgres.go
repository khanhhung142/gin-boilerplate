package postgres

import (
	"context"
	"database/sql"
	"gin-boilerplate/config"
	"log"

	_ "github.com/lib/pq"
)

type postgresClient struct {
	Client *sql.DB
}

var client postgresClient

func InitClient(ctx context.Context, config config.Config) {
	db, err := sql.Open("postgres", config.Database.ConnectString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	client = postgresClient{
		Client: db,
	}
}

func PostgresClient() postgresClient {
	return client
}

package postgres

import (
	"context"
	"database/sql"
	"gin-boilerplate/config"
	"gin-boilerplate/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap/zapcore"
)

var client *sql.DB

func InitClient(ctx context.Context, config *config.Config) {
	db, err := sql.Open("postgres", config.Database.ConnectString)
	if err != nil {
		logger.Fatal(ctx, "Failed to open a DB connection: ", zapcore.Field{Key: "error", Type: zapcore.StringType, String: err.Error()})
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(ctx, "Failed to ping DB: ", zapcore.Field{Key: "error", Type: zapcore.StringType, String: err.Error()})
	}

	client = db
}

func PostgresClient() *sql.DB {
	return client
}

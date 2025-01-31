package db

import (
	"context"
	"habbit-tracker/config"
	"habbit-tracker/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, config *config.Config) *gorm.DB {
	connection, err := gorm.Open(postgres.Open(config.Database.ConnectString), &gorm.Config{})
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", zap.Error(err))
	}

	return connection
}

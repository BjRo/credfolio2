// Package database provides PostgreSQL database connection management using Bun ORM.
package database

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"backend/internal/config"
)

// Connect establishes a connection to the PostgreSQL database using Bun ORM.
// It returns a configured *bun.DB instance ready for queries.
func Connect(ctx context.Context, cfg config.DatabaseConfig) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.URL())))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

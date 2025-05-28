package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/og11423074s/asyncapi/config"
	"time"
)

func NewPostgres(conf *config.Config) (*sql.DB, error) {
	dsn := conf.DatabaseUrl()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//if err != nil {
	//	return nil, fmt.Errorf("error creating database driver: %w", err)
	//}
	//
	//m, err := migrate.NewWithDatabaseInstance(
	//	fmt.Sprintf("file:///%s/migrations", conf.ProjectRoot),
	//	"postgres", driver)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating migration instance: %w", err)
	//}
	//
	//if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//	return nil, fmt.Errorf("error running migrations: %w", err)
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	return db, nil
}

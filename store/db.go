package store

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/jparkkennaby/docker-kafka-go/store/gen"
)

//go:embed migrations
var fs embed.FS

type Config struct {
	DatabaseURL string
}

type DB struct {
	conn    *pgxpool.Pool
	queries *gen.Queries
}

func NewDB(ctx context.Context, config Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %v", err)
	}
	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   Logger{},
		LogLevel: tracelog.LogLevelDebug,
	}

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("pool new: %v", err)
	}

	err = ensureSchema(conn)
	if err != nil {
		return nil, fmt.Errorf("ensure schema: %v", err)
	}

	return &DB{
		conn:    conn,
		queries: gen.New(conn),
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

func (db *DB) AddEventRecord(ctx context.Context, request gen.AddEventRecordParams) error {
	return db.queries.AddEventRecord(ctx, request)
}

func (db *DB) GetEventRecord(ctx context.Context, eventID int64) (int64, error) {
	return db.queries.GetEventRecord(ctx, eventID)
}

func ensureSchema(conn *pgxpool.Pool) error {
	sourceInstance, err := httpfs.New(http.FS(fs), "migrations")
	if err != nil {
		return fmt.Errorf("invalid soure instance: %v", err)
	}

	db := stdlib.OpenDBFromPool(conn)
	defer db.Close()

	targetInstance, err := cockroachdb.WithInstance(db, new(cockroachdb.Config))
	if err != nil {
		return fmt.Errorf("invalid target instance: %v", err)
	}
	defer targetInstance.Close()

	m, err := migrate.NewWithInstance("httpfs", sourceInstance, "cockroachdb", targetInstance)
	if err != nil {
		return fmt.Errorf("migrate new: %v", err)
	}
	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return sourceInstance.Close()
}

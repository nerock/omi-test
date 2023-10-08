package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nerock/omi-test/logger/config"
	"github.com/nerock/omi-test/logger/internal"
	"github.com/nerock/omi-test/logger/internal/listener"
	"github.com/nerock/omi-test/logger/internal/storage"
	"github.com/nerock/omi-test/logger/pkg"
	"go.uber.org/zap"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	cfg, err := config.FromFile("config.json")
	if err != nil {
		stdlog.Fatal("load config:", err)
	}

	log, err := pkg.NewZapLogger(cfg.Log)
	if err != nil {
		stdlog.Fatal("could not load zap logger:", err)
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
			cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB))
	if err != nil {
		log.Fatal("could not connect to postgres", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("could not close db", zap.Error(err))
		}
	}()
	if err := db.Ping(); err != nil {
		log.Fatal("could not ping postgres", zap.Error(err))
	}
	if err := runMigrations(db); err != nil {
		log.Fatal("could not run db migrations", zap.Error(err))
	}

	nc, err := nats.Connect(cfg.Nats.Uri)
	if err != nil {
		log.Fatal("could not connect to nats", zap.Error(err))
	}
	defer nc.Close()

	st := storage.NewAuditPostgresStorage(db)
	svc := internal.NewAuditService(st)
	listener.Subscribe(ctx, nc, cfg.Nats.Topic, svc, log)

	fmt.Printf("Listening to events on NATS topic: %s\n", cfg.Nats.Topic)
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"github.com/nerock/omi-test/account/internal/interceptor"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nerock/omi-test/account/config"
	"github.com/nerock/omi-test/account/internal"
	"github.com/nerock/omi-test/account/internal/api"
	"github.com/nerock/omi-test/account/internal/storage"
	"github.com/nerock/omi-test/account/pkg"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
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

	ec, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		log.Fatal("could not create proto encoded connection to nats", zap.Error(err))
	}
	defer ec.Close()

	st := storage.NewAccountPostgresStorage(db)
	svc := internal.NewAccountService(st)
	api := api.NewAccountGrpcApi(interceptor.NewAccountNatsInterceptor(svc, ec, cfg.Nats.Topic, log))

	srv := pkg.NewGrpcServer(cfg.Server.GrpcPort, log, api)

	gw, err := pkg.NewGrpcGW(ctx, cfg.Server, api)
	if err != nil {
		log.Fatal("could not setup grpc gateway", zap.Error(err))
	}

	go func() {
		if err := srv.Serve(); err != nil {
			log.Fatal("could not serve grpc server", zap.Error(err))
		}
	}()

	go func() {
		if err := gw.Serve(); err != nil {
			log.Fatal("could not serve grpc gw", zap.Error(err))
		}
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	closeChan := make(chan struct{})
	go func() {
		if err := gw.Shutdown(ctx); err != nil {
			log.Error("could not gracefully shutdown grpc gateway server", zap.Error(err))
		}

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("could not gracefully close grpc client", zap.Error(err))
		}

		closeChan <- struct{}{}
	}()

	<-closeChan
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

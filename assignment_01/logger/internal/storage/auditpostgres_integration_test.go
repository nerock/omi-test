//go:build integration
// +build integration

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nerock/omi-test/logger/config"
	"github.com/nerock/omi-test/logger/internal"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestAuditLogStorage(t *testing.T) {
	Convey("SaveLog", t, func() {
		ctx := context.Background()
		Convey("with a valid config", func() {
			cfg, err := config.FromFile("../../config_test.json")
			So(err, ShouldBeNil)

			Convey("and a valid connection to db", func() {
				db, err := sql.Open("postgres",
					fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
						cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB))
				So(err, ShouldBeNil)
				defer func() {
					So(db.Close(), ShouldBeNil)
				}()
				So(db.Ping(), ShouldBeNil)

				Convey("run the migrations", func() {
					m, err := setupMigrations(db)
					So(err, ShouldBeNil)
					So(m.Up(), ShouldBeNil)
					defer func() {
						So(m.Down(), ShouldBeNil)
					}()

					Convey("start storage", func() {
						st := NewAuditPostgresStorage(db)

						Convey("we can save new logs in the db", func() {
							newLog := internal.AuditLog{
								EventTypeID: 10,
								EventType:   "ACTION_HAPPENED",
								Timestamp:   time.Now(),
								UserIP:      "192.168.1.1",
								Data:        []byte(`{"changes": "theChange"}`),
							}

							So(st.SaveLog(ctx, newLog), ShouldBeNil)

							Convey("and the log was created", func() {
								var count int
								err := db.QueryRowContext(ctx, "SELECT count(*) FROM audit_logs").Scan(&count)
								So(err, ShouldBeNil)
								So(count, ShouldEqual, 1)
							})
						})
					})
				})
			})
		})
	})
}

func setupMigrations(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil

}

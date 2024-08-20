package database

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"time"
)

//go:embed schema.sql
var Schema string

const interval = 1 * time.Hour

// ConfigureSqlite sets some sensible defaults for sqlite
func ConfigureSqlite(db *sql.DB) error {
	pragmas := []string{
		"busy_timeout = 5000",
		"journal_mode = WAL",
		"synchronous = NORMAL",
		"cache_size = 1000000000", // 1GB
		"foreign_keys = true",
		"temp_store = memory",
		"mmap_size = 3000000000",
	}

	for _, pragma := range pragmas {
		_, err := db.Exec("PRAGMA " + pragma)
		if err != nil {
			return err
		}
	}
	return nil
}

// ScheduleCleanup deletes checks older than deleteOlderThan
// every interval. The format of deleteOlderThan is a string
// that can be parsed by sqlite, for example "-1 days".
func ScheduleCleanup(ctx context.Context, db *sql.DB, deleteOlderThan string) {
	queries := New(db)
	for {
		slog.Info("cleaning up old checks")

		tx, err := db.Begin()
		if err != nil {
			slog.Error("unable to start transaction", "error", err)
			continue
		}

		queriesTx := queries.WithTx(tx)

		if err := queriesTx.ChecksCleanup(ctx, deleteOlderThan); err != nil {
			slog.Error("unable to cleanup checks", "error", err)
		}

		changes, err := queriesTx.Changes(ctx)
		if err != nil {
			slog.Error("unable to get changes", "error", err)
		}

		if err := tx.Commit(); err != nil {
			slog.Error("unable to commit transaction", "error", err)
		}

		slog.Info("cleaned up old checks", "changes", changes)

		time.Sleep(interval)
	}
}

package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"minwa/internal/checker"
	"minwa/internal/database"
	"minwa/internal/notify"
	"minwa/internal/web"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var port = flag.String("port", "8080", "port")
var dbName = flag.String("dbname", "db.sqlite", "database name")

func main() {
	flag.Parse()

	ctx := context.Background()

	db, err := sql.Open("sqlite3", *dbName)
	if err != nil {
		slog.Error("unable to open sqlite db", "error", err)
		os.Exit(-1)
	}

	if err := database.ConfigureSqlite(db); err != nil {
		slog.Error("unable to open sqlite db", "error", err)
		os.Exit(-1)
	}

	if _, err := db.ExecContext(ctx, database.Schema); err != nil {
		slog.Error("unable to exec schema", "error", err)
		os.Exit(-1)
	}

	go checker.ScheduleCheck(
		ctx,
		db,
		notify.Config{
			From:  os.Getenv("MAIL_FROM"),
			To:    os.Getenv("MAIL_TO"),
			Token: os.Getenv("POSTMARK_TOKEN"),
		},
		1*time.Minute,
	)

	go database.ScheduleCleanup(ctx, db, "-1 days")

	hs := web.NewHttpServer(db)

	slog.Info("starting server", "url", "http://localhost:"+*port)
	http.ListenAndServe(":"+*port, hs.Server)
}

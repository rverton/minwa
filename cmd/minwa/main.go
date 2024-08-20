package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"minwa/internal/checker"
	"minwa/internal/database"
	"minwa/internal/web"
	"net/http"
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
	}

	if _, err := db.ExecContext(ctx, database.Schema); err != nil {
		slog.Error("unable to exec schema", "error", err)
	}

	go checker.ScheduleCheck(ctx, db, 30*time.Second)

	hs := web.NewHttpServer(db)

	slog.Info("starting server", "url", "http://localhost:"+*port)
	http.ListenAndServe(":"+*port, hs.Server)
}

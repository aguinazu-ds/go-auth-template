package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var Bun *bun.DB

func CreateDatabase(
	dbname string,
	dbuser string,
	dbpassword string,
	dbhost string,
) (*sql.DB, error) {
	hostArr := strings.Split(dbhost, ":")
	host := hostArr[0]
	port := "5432"
	if len(hostArr) > 1 {
		port = hostArr[1]
	}
	uri := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbuser,
		dbpassword,
		dbname,
		host,
		port,
	)
	slog.Info("Connecting to database " + dbname)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	slog.Info("Pinging database " + dbname)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	slog.Info("Connected to database " + dbname)
	return db, nil
}

func Init() error {
	var (
		host   = os.Getenv("DB_HOST")
		user   = os.Getenv("DB_USER")
		pass   = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
	)
	db, err := CreateDatabase(dbname, user, pass, host)
	if err != nil {
		log.Fatal(err)
	}
	Bun = bun.NewDB(db, pgdialect.New())
	if os.Getenv("APP_DEBUG") == "true" {
		slog.Info("Enabling bun debug")
		Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return nil
}

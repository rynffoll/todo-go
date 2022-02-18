package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	migrate "github.com/rubenv/sql-migrate"
	"net/http"
	"todo-go/todos"

	"github.com/jmoiron/sqlx"
)

// TODO: kafka
// TODO: https://cobra.dev/ --- cmd package
// TODO: OpenTelemetry
func main() {
	setupLogger()

	// TODO extract settings to external config (env?)
	db, err := setupDB()
	if err != nil {
		log.Fatal().Err(err).Stack().Msg("DB setup failed")
		return
	}
	defer db.Close()

	migrations(db.DB)

	s := &todos.TodoService{DB: db}

	r := todos.Route(s)
	r.Mount("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal().Err(err).Stack().Msg("Server failed")
	}
}

func setupLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setupDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Stack().Msg("DB open failed")
		return db, err
	}
	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Stack().Msg("DB ping failed")
		return db, err
	}

	db.SetMaxOpenConns(5)

	return db, err
}

func migrations(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{Dir: "migrations"}
	if _, err := migrate.Exec(db, "postgres", migrations, migrate.Up); err != nil {
		log.Fatal().Err(err).Stack().Msg("Migrations failed")
		return
	}
	log.Info().Msg("Migrations completed")
}

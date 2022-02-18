package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	migrate "github.com/rubenv/sql-migrate"
	"net/http"
	"time"
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

	r := chi.NewRouter()

	r.Use(hlog.NewHandler(log.Logger))
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	r.Use(hlog.RemoteAddrHandler("ip"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

	r.Mount("/", todos.Route(s))
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

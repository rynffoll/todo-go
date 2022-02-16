package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"todo-go/todos"

	"github.com/jmoiron/sqlx"
)

func main() {
	// TODO extract settings to external config (env?)
	db, err := sqlx.Open("postgres", "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.WithError(err).Error("Error during creation DB")
		return
	}
	db.SetMaxOpenConns(5)

	migrations(db.DB)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	s := &todos.TodoService{DB: db}
	r.Mount("/todos", todos.TodosRoute(s))

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.WithError(err).Error("Error during starting server")
	}
}

func migrations(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.WithError(err).Error("Migrations complete!")
		return
	}
	log.WithField("n", n).Info("Migrations complete!")
}

package todos

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	ERROR_CODE = -1
)

type Service interface {
	Add(ctx context.Context, todo Todo) (Todo, error)
	GetAll(ctx context.Context) ([]Todo, error)
	Get(ctx context.Context, id int) (Todo, error)
	Update(ctx context.Context, t Todo) (int, error)
	Remove(ctx context.Context, id int) (int, error)
}

type TodoService struct {
	DB *sql.DB
}

func (r *TodoService) Add(ctx context.Context, todo Todo) (Todo, error) {
	stmt, err := r.DB.PrepareContext(
		ctx,
		`INSERT INTO todos (title, done, date)
                 VALUES ($1, $2, $3)
                 RETURNING id`,
	)
	if err != nil {
		log.WithError(err).Warn("Error during creating prepared stmt")
		return todo, err
	}
	defer stmt.Close()
	var id int
	if err := stmt.QueryRowContext(ctx, todo.Title, todo.Done, todo.Date).Scan(&id); err != nil {
		log.WithError(err).Warn("Error during reading row")
		return todo, err
	}
	todo.ID = id
	return todo, err
}

func (r *TodoService) GetAll(ctx context.Context) ([]Todo, error) {
	todos := make([]Todo, 0)
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, title, done, date
                   FROM todos`,
	)
	if err != nil {
		log.WithError(err).Warn("Error during crearing query")
		return todos, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id    int
			title string
			done  bool
			date  time.Time
		)
		if err := rows.Scan(&id, &title, &done, &date); err != nil {
			log.WithError(err).Warn("Error during reading row")
			return todos, err
		}
		todos = append(todos, Todo{ID: id, Title: title, Done: done, Date: date})
	}
	return todos, err
}

func (r *TodoService) Get(ctx context.Context, id int) (Todo, error) {
	var todo Todo
	stmt, err := r.DB.PrepareContext(
		ctx,
		`SELECT title, done, date
                   FROM todos
                  WHERE id = $1`,
	)
	if err != nil {
		log.WithError(err).Warn("Error during creating stmt")
		return todo, err
	}
	var (
		title string
		done  bool
		date  time.Time
	)
	if err := stmt.QueryRowContext(ctx, id).Scan(&title, &done, &date); err != nil {
		log.WithError(err).Warn("Error during reading row")
		return todo, err
	}
	todo = Todo{ID: id, Title: title, Done: done, Date: date}
	return todo, err
}

func (r *TodoService) Update(ctx context.Context, todo Todo) (int, error) {
	stmt, err := r.DB.PrepareContext(
		ctx,
		`UPDATE todos
                    SET title = $1, 
                        done  = $2, 
                        date  = $3
                  WHERE id    = $4`,
	)
	if err != nil {
		log.WithError(err).Warn("Error during creating prepared stmt")
		return ERROR_CODE, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, todo.Title, todo.Done, todo.Date, todo.ID)
	if err != nil {
		log.WithError(err).Warn("Error during execution stmt")
		return ERROR_CODE, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.WithError(err).Warn("Error during reading result")
		return ERROR_CODE, err
	}
	return int(count), err
}

func (r *TodoService) Remove(ctx context.Context, id int) (int, error) {
	stmt, err := r.DB.PrepareContext(
		ctx,
		`DELETE FROM todos
                  WHERE id = $1`,
	)
	if err != nil {
		log.WithError(err).Warn("Error during creation prepared stmt")
		return ERROR_CODE, err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.WithError(err).Warn("Error during execution stmt")
		return ERROR_CODE, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.WithError(err).Warn("Error during reading result")
		return ERROR_CODE, err
	}
	return int(count), nil
}

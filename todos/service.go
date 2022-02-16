package todos

import (
	"context"

	"github.com/jmoiron/sqlx"
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
	DB *sqlx.DB
}

func (r *TodoService) Add(ctx context.Context, todo Todo) (Todo, error) {
	stmt, err := r.DB.PrepareNamedContext(
		ctx,
		`INSERT INTO todos (title, done, date)
                 VALUES (:title, :done, :date)
                 RETURNING id`,
	)
	defer stmt.Close()

	if err != nil {
		log.WithError(err).Warn("Error during creating prepared stmt")
		return todo, err
	}

	if err = stmt.Get(&todo.ID, todo); err != nil {
		log.WithError(err).Warn("Error during reading row")
		return todo, err
	}

	return todo, err
}

func (r *TodoService) GetAll(ctx context.Context) ([]Todo, error) {
	todos := []Todo{}

	err := r.DB.SelectContext(
		ctx,
		&todos,
		`SELECT id, title, done, date
                   FROM todos`,
	)
	if err != nil {
		log.WithError(err).Warn("Error during exec query")
		return todos, err
	}

	return todos, err
}

func (r *TodoService) Get(ctx context.Context, id int) (Todo, error) {
	todo := Todo{}

	err := r.DB.GetContext(
		ctx,
		&todo,
		`SELECT title, done, date
                   FROM todos
                  WHERE id = $1`,
		id,
	)
	if err != nil {
		log.WithError(err).Warn("Error during exec stmt")
		return todo, err
	}

	return todo, err
}

func (r *TodoService) Update(ctx context.Context, todo Todo) (int, error) {
	res, err := r.DB.NamedExecContext(
		ctx,
		`UPDATE todos
                    SET title = :title, 
                        done  = :done, 
                        date  = :date
                  WHERE id    = :id`,
		todo,
	)
	if err != nil {
		log.WithError(err).Warn("Error during update")
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
	res, err := r.DB.ExecContext(
		ctx,
		`DELETE FROM todos
                  WHERE id = $1`,
		id,
	)

	if err != nil {
		log.WithError(err).Warn("Error during creation prepared stmt")
		return ERROR_CODE, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.WithError(err).Warn("Error during reading result")
		return ERROR_CODE, err
	}

	return int(count), nil
}

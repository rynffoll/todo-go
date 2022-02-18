package todos

import (
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//go:embed sql/add.sql
var addQuery string

//go:embed sql/getAll.sql
var getAllQuery string

//go:embed sql/get.sql
var getQuery string

//go:embed sql/update.sql
var updateQuery string

//go:embed sql/remove.sql
var removeQuery string

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
	stmt, err := r.DB.PrepareNamedContext(ctx, addQuery)
	defer stmt.Close()
	if err != nil {
		return todo, errors.WithStack(err)
	}
	if err = stmt.Get(&todo.ID, todo); err != nil {
		return todo, errors.WithStack(err)
	}
	return todo, errors.WithStack(err)
}

func (r *TodoService) GetAll(ctx context.Context) ([]Todo, error) {
	todos := []Todo{}
	err := r.DB.SelectContext(ctx, &todos, getAllQuery)
	if err != nil {
		return todos, errors.WithStack(err)
	}
	return todos, errors.WithStack(err)
}

func (r *TodoService) Get(ctx context.Context, id int) (Todo, error) {
	todo := Todo{}
	err := r.DB.GetContext(ctx, &todo, getQuery, id)
	if err != nil {
		return todo, errors.WithStack(err)
	}
	return todo, errors.WithStack(err)
}

func (r *TodoService) Update(ctx context.Context, todo Todo) (int, error) {
	res, err := r.DB.NamedExecContext(ctx, updateQuery, todo)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return int(count), errors.WithStack(err)
}

func (r *TodoService) Remove(ctx context.Context, id int) (int, error) {
	res, err := r.DB.ExecContext(ctx, removeQuery, id)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return int(count), errors.WithStack(err)
}

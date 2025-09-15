package repository

import (
	"context"

	"to-do-list/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository interface {
	List(ctx context.Context) ([]model.Todo, error)
	Create(ctx context.Context, title string) (model.Todo, error)
	Update(ctx context.Context, id int, title *string, done *bool) (model.Todo, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type PostgresTodoRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTodoRepository(pool *pgxpool.Pool) *PostgresTodoRepository {
	return &PostgresTodoRepository{pool: pool}
}

func (r *PostgresTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
	rows, err := r.pool.Query(ctx, "SELECT id,title,done,created_at FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *PostgresTodoRepository) Create(ctx context.Context, title string) (model.Todo, error) {
	var t model.Todo
	err := r.pool.QueryRow(ctx, "INSERT INTO todos(title) VALUES($1) RETURNING id,title,done,created_at", title).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	return t, err
}

func (r *PostgresTodoRepository) Update(ctx context.Context, id int, title *string, done *bool) (model.Todo, error) {
	var t model.Todo
	err := r.pool.QueryRow(ctx, "UPDATE todos SET title=COALESCE($1,title), done=COALESCE($2,done) WHERE id=$3 RETURNING id,title,done,created_at", title, done, id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return t, pgx.ErrNoRows
	}
	return t, err
}

func (r *PostgresTodoRepository) Delete(ctx context.Context, id int) (bool, error) {
	cmd, err := r.pool.Exec(ctx, "DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}

package service

import (
	"context"
	"errors"

	"to-do-list/internal/model"
	"to-do-list/internal/repository"
)

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) List(ctx context.Context) ([]model.Todo, error) {
	return s.repo.List(ctx)
}

func (s *TodoService) Create(ctx context.Context, title string) (model.Todo, error) {
	if title == "" {
		return model.Todo{}, errors.New("title is required")
	}
	return s.repo.Create(ctx, title)
}

func (s *TodoService) Update(ctx context.Context, id int, title *string, done *bool) (model.Todo, error) {
	return s.repo.Update(ctx, id, title, done)
}

func (s *TodoService) Delete(ctx context.Context, id int) (bool, error) {
	return s.repo.Delete(ctx, id)
}

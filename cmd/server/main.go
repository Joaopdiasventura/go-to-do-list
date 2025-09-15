package main

import (
	"context"
	"log"

	"to-do-list/internal/api"
	"to-do-list/internal/config"
	"to-do-list/internal/db"
	"to-do-list/internal/repository"
	"to-do-list/internal/service"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	repo := repository.NewPostgresTodoRepository(pool)
	svc := service.NewTodoService(repo)
	router := api.NewRouter()
	ctrl := api.NewTodoController(svc)
	ctrl.Register(router)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

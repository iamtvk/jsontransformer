package main

import (
	"database/sql"
	"log"

	"github.com/iamtvk/jsontransformer/internal/config"
	"github.com/iamtvk/jsontransformer/internal/repository/postgres"
	"github.com/iamtvk/jsontransformer/internal/service"
	"github.com/iamtvk/jsontransformer/internal/transport/http"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("postgre", cfg.DbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()
	cache := service.NewCacheLayer()
	repo := postgres.NewPostgreSQLRepository(db)
	transformerService := service.NewTransformerService(repo, cfg, cache)

	httpHandler := http.NewHandler(transformerService)

	httpHandler.Transform()
}

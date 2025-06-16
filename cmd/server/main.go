package main

import (
	"database/sql"
	"github.com/iamtvk/jsontransformer/internal/config"
	"github.com/iamtvk/jsontransformer/internal/repository/postgres"
	"github.com/iamtvk/jsontransformer/internal/service"
	httptransport "github.com/iamtvk/jsontransformer/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()
	cache := service.NewCacheLayer()
	repo := postgres.NewPostgreSQLRepository(db)
	transformerService := service.NewTransformerService(repo, cfg, cache)
	httptransport.NewHandler(transformerService)
	http.Handle()

}

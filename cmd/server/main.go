package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/iamtvk/jsontransformer/api/proto/transformerPb"
	"github.com/iamtvk/jsontransformer/internal/config"
	"github.com/iamtvk/jsontransformer/internal/repository/postgres"
	"github.com/iamtvk/jsontransformer/internal/service"
	grpcTransport "github.com/iamtvk/jsontransformer/internal/transport/grpc"
	httpTransport "github.com/iamtvk/jsontransformer/internal/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()

	// test db conn
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db:", err)
	}
	log.Print("db conn successfull")
	cache := service.NewCacheLayer()
	repo := postgres.NewPostgreSQLRepository(db)
	transformerService := service.NewTransformerService(repo, cfg, cache)
	httpHandler := httpTransport.NewHandler(transformerService)
	grpcServer := grpcTransport.NewServer(transformerService)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		startHttpServer(cfg.HTTPPort, httpHandler)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRpcServer(cfg.GRPCPort, grpcServer)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("*** Application server started running, press CTRL C to shutdown ***")
	<-sigChan

	log.Printf("Shuting Down")
	wg.Wait()
	log.Printf("Shutdown complete")
}

func startHttpServer(port string, handler http.Handler) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("http server starting on port:%s ", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server failed to start: %v", err)
	}

}
func startGRpcServer(port string, handler *grpcTransport.Server) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("error starting grpc server on port: %s, error: %v", port, err)
	}
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(4 * 1024 * 1024), // 4MB
		grpc.MaxSendMsgSize(4 * 1024 * 1024),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterTransformerServiceServer(server, handler)
	reflection.Register(server)
	log.Printf("gRPC server starting on port: %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("grpc server failed to start, error: %v", err)
	}
}

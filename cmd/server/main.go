package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
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
	_ "github.com/lib/pq"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		startHttpServer(ctx, cfg.HTTPPort, httpHandler)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRpcServer(ctx, cfg.GRPCPort, grpcServer)
	}()

	// sigChan := make(chan os.Signal, 1)
	log.Printf("*** Application server started running, press CTRL C to shutdown ***")
	<-ctx.Done()

	log.Printf("Shuting Down")
	wg.Wait()
	log.Printf("Shutdown complete")
}

func startHttpServer(ctx context.Context, port string, handler http.Handler) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		<-ctx.Done()
		log.Printf("Shutting down http server on port %s ...", port)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}
	}()
	log.Printf("http server starting on port:%s ", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server failed to start: %v", err)
	}
}

func startGRpcServer(ctx context.Context, port string, handler *grpcTransport.Server) {
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

	go func() {
		<-ctx.Done()
		log.Printf("Shuttingdown gRPC server on port %s ...", port)
		done := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
			log.Println("gRPC server shutdown gracefully")
		case <-time.After(5 * time.Second):
			log.Println("gRPC server shutdown shutdown timeout, forcing stop")
			server.Stop()
		}
	}()
	log.Printf("gRPC server starting on port: %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("grpc server failed to start, error: %v", err)
	}
}

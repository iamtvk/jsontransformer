package grpctest

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"transformPb"
)

func main() {
	// Create a gRPC connection using the new API
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := transformerPb.NewTransformerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Example usage
	resp, err := client.Transform(ctx, &transformerPb.TransformRequest{
		ScriptIdentifier: "my-script",
		Data:             []byte(`{"key":"value"}`),
	})
	if err != nil {
		log.Fatalf("Transform error: %v", err)
	}

	log.Printf("Transform result: %s", string(resp.GetResult()))
}

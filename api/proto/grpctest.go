package main

import (
	"context"
	"log"
	"time"

	"github.com/iamtvk/jsontransformer/api/proto/transformerPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a gRPC connection using the new API
	conn, err := grpc.NewClient("localhost:9090",
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
		ScriptIdentifier: "address-script-001",
		Data: []byte(`
  "FirstName": "Fred",
  "Surname": "Smith",
  "Age": 28,
  "Address": {
    "Street": "Hursley Park",
    "City": "Winchester",
    "Postcode": "SO21 2JN"
  },
  "Phone": [
    {
      "type": "home",
      "number": "0203 544 1234"
    },
    {
      "type": "office",
      "number": "01962 001234"
    },
    {
      "type": "office",
      "number": "01962 001235"
    },
    {
      "type": "mobile",
      "number": "077 7700 1234"
    }
  ],
  "Email": [
    {
      "type": "office",
      "address": [
        "fred.smith@my-work.com",
        "fsmith@my-work.com"
      ]
    },
    {
      "type": "home",
      "address": [
        "freddy@my-social.com",
        "frederic.smith@very-serious.com"
      ]
    }
  ],
  "Other": {
    "Over 18 ?": true,
    "Misc": null,
    "Alternative.Address": {
      "Street": "Brick Lane",
      "City": "London",
      "Postcode": "E1 6RF"
    }
  }
}`),
	})
	if err != nil {
		log.Fatalf("Transform error: %v", err)
	}

	log.Printf("Transform result: %s", string(resp.GetResult()))
}

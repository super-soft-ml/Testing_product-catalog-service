package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc"

	"product-catalog-service/internal/services"
	productv1 "product-catalog-service/proto/product/v1"
)

func main() {
	project := os.Getenv("SPANNER_PROJECT")
	if project == "" {
		project = "test-project"
	}
	instance := os.Getenv("SPANNER_INSTANCE")
	if instance == "" {
		instance = "test-instance"
	}
	database := os.Getenv("SPANNER_DATABASE")
	if database == "" {
		database = "product-catalog"
	}
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instance, database)

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		log.Fatalf("spanner.NewClient: %v", err)
	}
	defer client.Close()

	opts := services.NewOptions(client)
	handler := opts.ProductHandler()

	grpcServer := grpc.NewServer()
	productv1.RegisterProductServiceServer(grpcServer, handler)

	addr := ":50051"
	if a := os.Getenv("GRPC_ADDR"); a != "" {
		addr = a
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("gRPC server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

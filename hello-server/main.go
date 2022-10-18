package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/rafarlopes/grpc-with-go/hello"
	"google.golang.org/grpc"
)

type server struct {
	hello.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Printf("Received: %s", req.GetName())

	return &hello.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, &server{})

	log.Printf("\nserver listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

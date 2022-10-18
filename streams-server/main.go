package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rafarlopes/grpc-with-go/streams"
	"google.golang.org/grpc"
)

type server struct {
	streams.UnimplementedGreeterServer
}

func (s *server) SayHello(stream streams.Greeter_SayHelloServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Received: %s", req.GetName())

		if err := stream.Send(&streams.HelloResponse{
			Message: fmt.Sprintf("Hello %s", req.GetName()),
		}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	streams.RegisterGreeterServer(s, &server{})

	log.Printf("\nserver listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

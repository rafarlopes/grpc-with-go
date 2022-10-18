package main

import (
	"context"
	"log"
	"time"

	"github.com/rafarlopes/grpc-with-go/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := hello.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &hello.HelloRequest{Name: "test"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Greeting: %s", resp.GetMessage())
}

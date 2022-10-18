package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/rafarlopes/grpc-with-go/streams"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := streams.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := c.SayHello(ctx)
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan interface{})

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(wait)
				return
			}
			if err != nil {
				log.Fatalf("\nmessage failed %v", err)
			}

			log.Printf("\ngot message %s", resp.GetMessage())
		}
	}()

	messages := []*streams.HelloRequest{
		{Name: "aaa"},
		{Name: "bbb"},
		{Name: "ccc"},
		{Name: "ddd"},
	}

	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			log.Fatalf("\nfailed to send message %v", err)
		}
	}

	stream.CloseSend()

	<-wait
}

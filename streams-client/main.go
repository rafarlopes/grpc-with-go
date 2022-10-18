package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/rafarlopes/grpc-with-go/streams"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mu := &sync.Mutex{}

	messages := []*streams.HelloRequest{
		{Name: "aaa"},
		{Name: "bbb"},
		{Name: "ccc"},
		{Name: "ddd"},
	}

	backoff.Retry(func() error {
		log.Println("starting...")
		conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		defer conn.Close()

		c := streams.NewGreeterClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		stream, err := c.SayHello(ctx)
		if err != nil {
			return err
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
					log.Printf("\nmessage failed %v", err)
					return
				}

				log.Printf("\ngot message %s", resp.GetMessage())
				mu.Lock()
				messages = messages[1:]
				mu.Unlock()
			}
		}()

		for _, msg := range messages {
			if err := stream.Send(msg); err != nil {
				return fmt.Errorf("\nfailed to send message %v", err)
			}

			time.Sleep(10 * time.Second)
		}

		stream.CloseSend()

		<-wait

		return nil

	}, backoff.NewConstantBackOff(5*time.Second))
}

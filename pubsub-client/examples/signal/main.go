package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/getjump/pubsub/pubsub-client/pkg"
	"github.com/getjump/pubsub/pubsub-client/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:9000", "the address to connect to")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %w", err)
	}

	defer conn.Close()

	c := pkg.NewPubSubClient(conn, "testClient")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	go func() {
		mc, err := c.Subscribe(ctx, "test")
		if err != nil {
			log.Fatalf("error on subscribe: %w", err)
		}

		for m := range mc {
			log.Println("Message", m.Data)
		}
	}()

	for {
		err = c.PublishAsync(ctx, "test", &proto.Signal{}, true)

		if err != nil {
			log.Fatalf("error on publish: %w", err)
		}

		<-time.After(time.Millisecond * 50)
	}
}

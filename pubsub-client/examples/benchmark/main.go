package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/getjump/pubsub/pubsub-client/pkg"
	"github.com/getjump/pubsub/pubsub-client/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr        = flag.String("addr", "localhost:9000", "the address to connect to")
	concurrency = flag.Int("concurrency", 1, "concurrency")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %w", err)
	}

	defer conn.Close()

	c := pkg.NewPubSubClient(conn, "testClient")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	timeStart := time.Now()

	readerCounter := uint64(0)
	writerCounter := uint64(0)

	count := 1000

	wgr := sync.WaitGroup{}
	wgw := sync.WaitGroup{}

	wgr.Add(*concurrency)
	wgw.Add(*concurrency)

	for i := 0; i < *concurrency; i++ {
		go reader(ctx, c, count, &readerCounter, &wgr)
		go writer(ctx, c, count, &writerCounter, &wgw)
	}

	wgr.Wait()
	fmt.Printf("RPS reader %f\n", float64(atomic.LoadUint64(&readerCounter))/time.Since(timeStart).Seconds())

	wgw.Wait()
	fmt.Printf("RPS writer %f\n", float64(atomic.LoadUint64(&writerCounter))/time.Since(timeStart).Seconds())
	fmt.Println(time.Since(timeStart).Seconds())
}

func reader(ctx context.Context, client *pkg.PubSubClient, count int, counter *uint64, wg *sync.WaitGroup) {
	defer wg.Done()

	mc, err := client.Subscribe(ctx, "bench")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < *concurrency*count; i++ {
		select {
		case <-mc:
			// case <-time.After(time.Millisecond * 50):
			// 	fmt.Println("timed out")
			// 	break
		}

		atomic.AddUint64(counter, 1)
		// fmt.Println(atomic.LoadUint64(counter))
	}
}

func writer(ctx context.Context, client *pkg.PubSubClient, count int, counter *uint64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count; i++ {
		err := client.PublishSync(ctx, "bench", &proto.Signal{}, false)
		if err != nil {
			fmt.Println(err)
			return
		}
		atomic.AddUint64(counter, 1)
		fmt.Println(i)
	}
}

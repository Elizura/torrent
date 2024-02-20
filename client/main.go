package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	pb "simplebittorrent/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDownloadFileClient(conn)
	var wg sync.WaitGroup
	wg.Add(1)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		defer wg.Done()
		if _, err := c.Download(ctx, &pb.DownloadRequest{}); err != nil {
			log.Fatalf("error downloading: %v", err)
		}
	}()

	wg.Wait()
}

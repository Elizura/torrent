package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	// "sync"

	pb "simplebittorrent/rpc"

	svr "simplebittorrent/server"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	pb.UnimplementedDownloadFileServer
}

func (s *server) Download(ctx context.Context, in *pb.DownloadRequest) (*pb.DownloadReply, error) {
	handleConnection()
	return &pb.DownloadReply{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDownloadFileServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func handleConnection() {
	done := make(chan struct{})

	go func() {
		defer close(done)
		svr.Seeder()
	}()
	svr.StartDownload("torrent-files/debian-11.6.0-amd64-netinst.iso.torrent")
	<-done
}

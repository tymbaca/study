package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/tymbaca/study/grpc-go/hsh"
	pb "github.com/tymbaca/study/grpc-go/message"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "the server port")
)

type Server struct {
	// pb.HasherServer
	pb.UnimplementedHasherServer
}

func (s *Server) Hash(ctx context.Context, in *pb.Input) (*pb.Output, error) {
	sum, _ := hsh.Hash(in.Data)
	log.Printf("new hash: '%X'", sum)
	// log.Printf("new hash (as string): '%v', '%s'", string(sum), string(sum))
	res := &pb.Output{Data: sum}
	return res, nil
}

func main() {
	flag.Parse()

	s := grpc.NewServer()
	pb.RegisterHasherServer(s, &Server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	log.Printf("listening port: %d", *port)
	log.Fatalf("failed to serve: %v", s.Serve(lis))
}

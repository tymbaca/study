package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"net"

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
	h := sha256.New()

	// Write to hash. Never returns an error
	_, err := h.Write(in.Data)
	if err != nil {
		return nil, fmt.Errorf("probably earth is broken: %w", err)
	}

	sum := h.Sum(nil)
	log.Printf("new hash: '%x'", sum)
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

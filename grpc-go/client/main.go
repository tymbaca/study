package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/tymbaca/study/grpc-go/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:8080", "address of server to connect")
	name = flag.String("name", "client", "client name")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewHasherClient(conn)
	in := &pb.Input{Data: []byte("Tigran")}
	res, err := c.Hash(context.Background(), in)
	if err != nil {
		panic(err)
	}

	log.Printf("Original message: '%s', Hash: '%s'", in.Data, res.Data)
}

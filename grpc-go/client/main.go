package main

import (
	"context"
	"flag"
	"log"
	"time"

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
	in := &pb.Input{Data: []byte(`The reason gRPC -- well, really protobufs -- doesn’t scale well in your example is that every entry of your repeated field results in protobuf needing to decode a separate field, and there is overhead related to that. You can see more details about the encoding of repeated fields in the docs here. You're using proto3, so at least you don't need to specify the [packed=true] option, although that helps somewhat if you're on proto2.

The reason switching to a string or bytes field speeds it up so much is that there is only a constant decoding cost for this field which doesn't scale with the amount of data that's encoded in the field (not sure about JS though, which might need to create a copy of the data, but clearly that is still much faster than actually parsing the data). Just make sure your protocol defines what format / endianness the data in the field is :-)

Answering your question at a higher level, sending multiple megabytes in a single API call is usually not an amazing idea anyway -- it ties up a thread on both the server and client for a long time which forces you to use multithreading or async code to get reasonable performance. (Admittedly might be less of an issue since you are used to writing async stuff on Node, but there's still only so many CPUs to burn on the server.)

Depending on what you're actually trying to do, a common pattern can be to write the data to a file in a shared storage system (S3, etc.) and pass the filename to the other service, which can then download it when it's actually needed.`)}

	// Start timer
	start := time.Now()
	// Calls
	for i := 0; i < 100; i++ {
		res, err := c.Hash(context.Background(), in)
		if err != nil {
			panic(err)
		}
		log.Printf("Original message: '%s', Hash: '%x'", in.Data, res.Data)
	}

	end := time.Since(start)
	log.Println(end)
}

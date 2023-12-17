package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileSystem struct{}

func (f *FileSystem) ReceiveFile(conn net.Conn) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	for {
		time.Sleep(1 * time.Second)

		var size int64
		binary.Read(conn, binary.LittleEndian, &size)

		n, err := io.CopyN(buf, conn, size)
		log.Printf("received %d bytes", n)
		if err != io.EOF {
			// it's ok to get EOF
			log.Println("Ending file receive")
			return buf, nil
		} else if err != nil {
			return nil, fmt.Errorf(
				"error while receiving file (chink with %d bytes): %s",
				n,
				err,
			)
		}
	}
}

func main() {
	fs := new(FileSystem)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	conn, err := lis.Accept()
	if err != nil {
		panic(err)
	}

	buf, err := fs.ReceiveFile(conn)
	if err != nil {
		panic(err)
	}
	_ = buf
	// log.Println(buf.Bytes())
}

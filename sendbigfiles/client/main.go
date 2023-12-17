package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func sendFile(conn net.Conn, size int64) error {
	data := make([]byte, size)
	// data := bytes.NewBuffer()
	_, err := rand.Reader.Read(data)
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, size)
	w, err := io.CopyN(conn, bytes.NewReader(data), size)
	if err != nil {
		return err
	}

	log.Printf("written %d bytes", w)
	return nil
}

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	log.Fatal(sendFile(conn, 2*1024*1024*1024))
}

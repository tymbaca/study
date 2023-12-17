package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	sock, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	conn, err := sock.Accept()
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	addr := conn.RemoteAddr()
	log.Printf("Conn with addr: %s", addr.String())

	conn.Write([]byte("FUCK YOU"))
	buf := make([]byte, 1024)
	for {
		r, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		log.Printf("Read %v bytes: %s", r, string(buf))
		if strings.Contains(string(buf), "FUCK YOU") {
			time.Sleep(1 * time.Second)
			sendBuf := []byte("NO FUCK YOU CLIENT")
			w, err := conn.Write(sendBuf)
			if err != nil {
				panic(err)
			}

			log.Printf("Send %v bytes: %s", w, string(sendBuf))
		}
	}
}

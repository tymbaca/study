package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	addr := conn.RemoteAddr()
	log.Printf("Connected with: %s", addr.String())

	buf := make([]byte, 1024)
	for {
		r, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		log.Printf("Read %v bytes: %s", r, string(buf))
		if strings.Contains(string(buf), "FUCK YOU") {
			time.Sleep(1 * time.Second)
			sendBuf := []byte("NO FUCK YOU SERVER")
			w, err := conn.Write(sendBuf)
			if err != nil {
				panic(err)
			}

			log.Printf("Send %v bytes: %s", w, string(sendBuf))
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tymbaca/study/grpc-go/hsh"
)

type Input struct {
	Data string `json:"data"`
}

type Output struct {
	Data string `json:"data"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		in := &Input{}
		dec.Decode(in)

		sum, _ := hsh.Hash([]byte(in.Data))
		log.Printf("new hash: '%X'", sum)

		out := &Output{Data: fmt.Sprintf("%x", sum)}
		data, err := json.Marshal(out)
		if err != nil {
			panic(err)
		}

		w.Write(data)
	})
	http.ListenAndServe(":8080", nil)
}

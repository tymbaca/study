package main

import (
	"log"
	"net/http"
)

func ClickHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// counter := 0

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/click", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte(fmt.Sprint(counter)))
		// counter++
		log.Println(r.FormValue("item"))
	})
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/ pinged")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ping-pong"))
		log.Println("/ pinged")
	})
	adr := "localhost:2001"
	log.Println("running at " + adr)

	err := http.ListenAndServe(adr, nil)

	if err != nil {
		log.Fatal("oh no, something went wrong" + err.Error())
	}
}

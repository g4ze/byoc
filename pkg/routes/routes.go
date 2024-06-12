package routes

import (
	"log"
	"net/http"

	"github.com/g4ze/byoc/pkg/handlers"
)

func Server() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/ pinged")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ping-pong"))
		log.Println("/ pinged")
	})

	http.HandleFunc("/make-cluster", handlers.Make_cluster)
	http.HandleFunc("/delete-cluster", handlers.Delete_cluster)
	http.HandleFunc("/deploy-container", handlers.Deploy_container)
	adr := "localhost:2001"
	log.Println("running at " + adr)
	err := http.ListenAndServe(adr, nil)

	if err != nil {
		log.Fatal("oh no, something went wrong" + err.Error())
	}
}

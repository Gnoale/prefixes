package main

import (
	"log"
	"net/http"
	"words/api"
	"words/repository"
)

func main() {

	mux := http.NewServeMux()

	h := api.NewHandler(repository.NewMemRepository())

	mux.HandleFunc("POST /words/{word}", h.InsertWord)
	mux.HandleFunc("GET /words/{prefix}", h.FindPrefix)
	log.Fatal(http.ListenAndServe(":8080", mux))

}

package main

import (
	"log"
	"net/http"
	"os"
	"words/api"
	"words/repository"
)

func main() {

	mux := http.NewServeMux()

	storeEnvSelector := os.Getenv("STORE_KIND")
	store, err := repository.StoreFactory(storeEnvSelector)
	if err != nil {
		panic(err)
	}
	h := api.NewHandler(store)

	mux.HandleFunc("POST /words/{word}", h.InsertWord)
	mux.HandleFunc("GET /words/{prefix}", h.FindPrefix)
	mux.HandleFunc("GET /words/", h.List)
	log.Fatal(http.ListenAndServe(":8080", mux))

}

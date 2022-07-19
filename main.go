package main

import (
	"createanswers/answers"
	"createanswers/cachedb"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	if passed := cachedb.InitRedis(); !passed {
		log.Fatal("failed to intialise cache db")
	}
	address := ":8080"
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/answer", answers.CreateAnswer).Methods("POST")
	router.HandleFunc("/answers/history/{key}", answers.GetHistory).Methods("GET")
	router.HandleFunc("/answers/{key}", answers.GetAnswerDetails).Methods("GET")
	router.HandleFunc("/answers", answers.UpdateAnswer).Methods("PATCH")
	router.HandleFunc("/answers/{key}", answers.DeleteAnswer).Methods("DELETE")
	log.Printf("main: starting the server on %v", address)
	log.Fatal(http.ListenAndServe(address, router))
}

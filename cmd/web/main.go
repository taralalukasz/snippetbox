package main

import (
	"log"
	"net/http"
)

func main() {
	//create mux server and register function to corresponding endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//this is how you start a web server
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000",mux)
	log.Fatal(err)
}
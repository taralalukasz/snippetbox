package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path != "/") {
		//this is an equivalent of 
		// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// w.Header().Set("X-Content-Type-Options", "nosniff")
		// w.WriteHeader(code)
		http.NotFound(w, r)  
		return
	}
	w.Write([]byte("Hello from snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a snippet"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create new snippet"))
}

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
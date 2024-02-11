package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Displaying snippet with id %v", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "POST") {
		//this is adding headers
		w.Header().Set("Allow", "POST")
		//you can call this method only once, it ultimately sets the header and response
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Method Not Allowed"))

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
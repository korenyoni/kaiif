package main

import (
	"net/http"
	"log"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	content, err := Asset(r.URL.Path[1:])
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	w.Write(content)
}

func serveHTTP(c chan error) {
	http.HandleFunc("/", handler)
	c <- http.ListenAndServe("0.0.0.0:8080", nil)
}

func main() {
	c := make(chan error)
	go serveHTTP(c)
	err := <-c
	if err != nil {
		log.Fatal(err)
	}
}

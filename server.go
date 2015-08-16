package main

import (
	"fmt"
	"time"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5000 * time.Millisecond)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

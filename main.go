package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/integralist/go-requester/requester"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No config file path provided")
		os.Exit(1)
	}

	configPath := os.Args[1] // zero index is the binary command itself

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requester.Process(w, r, configPath)
	})

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	// Handler
	fs := http.FileServer(http.Dir("."))

	mux.Handle("/", fs)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Println("Address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Errorf("Error fo ListenAndServer", err)
	}

}

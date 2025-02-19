package main

import (
	"net/http"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("contentType", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK, ready!"))

}

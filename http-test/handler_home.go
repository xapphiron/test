package main

import (
	"fmt"
	"net/http"
)

func handlerHome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Server is running.")
}
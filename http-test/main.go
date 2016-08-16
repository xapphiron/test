package main

import (
	"fmt"
	"net/http"
)

func writeResponseMessage(w http.ResponseWriter, status int, message string) {
	w.Header()["Content-Type"] = []string { "text/html; charset=UTF-8" }
	w.WriteHeader(status)
	fmt.Fprintf(w, "<div style=\"font:normal 16px Arial;\">" + message + "</div>")
}

func main() {
	fmt.Println("Starting server.")
    http.HandleFunc("/", handlerHome)
    http.HandleFunc("/im", imInterfaceHandler);
    http.ListenAndServe(":8080", nil)
}

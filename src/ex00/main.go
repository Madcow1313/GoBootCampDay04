package main

import (
	"fmt"
	"net/http"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r, r.Body)
	fmt.Fprintf(w, "Hi")
}

func main() {
	http.HandleFunc("/", simpleHandler)
	http.ListenAndServe(":80", nil)
}

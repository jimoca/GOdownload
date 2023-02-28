package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	dir := "./files"

	http.Handle("/", http.FileServer(http.Dir(dir)))

	log.Printf("Listening on port: %s...\n", port)
	log.Printf("Directory: %s\n", dir)

	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

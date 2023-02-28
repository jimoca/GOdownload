package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	var path string
	flag.StringVar(&path, "path", "", "Path")
	flag.Parse()
	fmt.Printf("path: %s\n", path)

	port := "8080"
	dir := path + "files"
	http.Handle("/", http.FileServer(http.Dir(dir)))

	log.Printf("Listening on port: %s...\n", port)
	log.Printf("Directory: %s\n", dir)

	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

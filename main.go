package main

import (
	"fmt"
	"github.com/georgerogers42/chatty-go/lib/chatty"
	"log"
	"net/http"
	"os"
)

func main() {
	http.Handle("/", chatty.App)
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	port := os.ExpandEnv(":$PORT")
	fmt.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

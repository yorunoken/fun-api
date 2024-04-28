package main

import (
	"fmt"
	"fun-api/handlers"
	"fun-api/utils"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/user", handlers.User)
	http.HandleFunc("/mood", handlers.Mood)
	http.HandleFunc("/hello", handlers.Hello)
	http.HandleFunc("/card", handlers.Card)
	http.HandleFunc("/media/", utils.MediaRedirector)

	fmt.Println("Listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(utils.GetPort("3000"), nil))
}

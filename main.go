package main

import (
	"fmt"
	"fun-api/api"
	"fun-api/handlers"
	"fun-api/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/juju/ratelimit"
)

var limiter = ratelimit.NewBucketWithRate(100, 100)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	go refreshToken()

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/mood", handlers.Mood)
	http.HandleFunc("/hello", handlers.Hello)
	http.HandleFunc("/card", handlers.Card)

	http.HandleFunc("/api/user", api.User)
	http.HandleFunc("/api/token", func(w http.ResponseWriter, r *http.Request) {
		if limiter.TakeAvailable(1) == 0 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		api.Token(w, r)
	})

	http.HandleFunc("/media/", utils.MediaRedirector)

	fmt.Println("Listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(utils.GetPort("3000"), nil))
}

func refreshToken() {
	secret := os.Getenv("secret")
	baseUrl := os.Getenv("base_url")

	utils.SetToken(secret, baseUrl)

	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		utils.SetToken(secret, baseUrl)
	}
}

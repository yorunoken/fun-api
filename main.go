package main

import (
	"fmt"
	"fun-api/api"
	"fun-api/api/beatmap"
	"fun-api/api/user"
	"fun-api/handlers"
	"fun-api/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file (most likely because it's in prod)")
	}

	go refreshToken()

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/mood", handlers.Mood)
	http.HandleFunc("/hello", handlers.Hello)
	http.HandleFunc("/card", handlers.Card)

	http.HandleFunc("/api/averagecolor", api.AverageColor)
	http.HandleFunc("/api/user/skills", user.Skills)
	http.HandleFunc("/api/user/details", user.Details)
	http.HandleFunc("/api/user/tops", user.Tops)
	http.HandleFunc("/api/user/avatar", user.Avatar)

	http.HandleFunc("/api/beatmap/download", beatmap.Download)

	http.HandleFunc("/api/graph", api.Graph)
	http.HandleFunc("/api/token", api.Token)

	http.HandleFunc("/media/", utils.MediaRedirector)

	fmt.Println("Listening on http://localhost:3001")
	log.Fatal(http.ListenAndServe(utils.GetPort("3001"), nil))
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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fun-api/utils"
	"net/http"
	"os"
)

func Token(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if secret != os.Getenv("secret") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

	payloadData := map[string]interface{}{
		"grant_type":    "client_credentials",
		"client_id":     os.Getenv("client_id"),
		"client_secret": os.Getenv("client_secret"),
		"scope":         "public",
		"code":          "code",
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

	data, err := utils.Post("https://osu.ppy.sh/oauth/token", bytes.NewReader(payloadBytes))

	if err != nil {
		fmt.Println("Error fetching OAuth token:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to fetch OAuth token"}`))
		return
	}

	var userResponse map[string]string

	if err := json.Unmarshal(data, &userResponse); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

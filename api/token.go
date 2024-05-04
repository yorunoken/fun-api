package api

import (
	"bytes"
	"encoding/json"
	"fun-api/utils"
	"net/http"
	"os"
)

func Token(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if secret != os.Getenv("secret") {
		utils.WriteError(w, "Internal Server Error")
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
		utils.WriteError(w, "Error while converting JSON")
		return
	}

	data, err := utils.Post("https://osu.ppy.sh/oauth/token", bytes.NewReader(payloadBytes))

	if err != nil {
		utils.WriteError(w, "Failed to fetch OAuth token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

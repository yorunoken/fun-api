package user

import (
	"fmt"
	"fun-api/utils"
	"net/http"
	"os"
	"strings"
)

func Details(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	mode := strings.ToLower(r.URL.Query().Get("mode"))

	if mode == "" {
		mode = "osu"
	}

	if username == "" {
		utils.WriteError(w, "Missing username parameter")
		return
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "Accept: application/json",
		"Authorization": "Bearer " + os.Getenv("access_token"),
	}

	data, err := utils.Get(fmt.Sprintf("https://osu.ppy.sh/api/v2/users/%s/%s", username, mode), headers)

	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Failed to fetch user data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

package user

import (
	"fmt"
	"fun-api/utils"
	"net/http"
	"os"
)

func Details(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	mode := r.URL.Query().Get("mode")

	if username == "" || mode == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing username or mode"}`))
		return
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "Accept: application/json",
		"Authorization": "Bearer " + os.Getenv("access_token"),
	}

	data, err := utils.Get(fmt.Sprintf("https://osu.ppy.sh/api/v2/users/%s/%s", username, mode), headers)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to fetch user data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

package user

import (
	"fmt"
	"fun-api/utils"
	"net/http"
	"os"
	"strings"
)

func Tops(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	mode := r.URL.Query().Get("mode")
	scoreType := strings.ToLower(r.URL.Query().Get("type"))

	if scoreType != "best" && scoreType != "firsts" && scoreType != "recent" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "type parameter must be one of 'best', 'firsts', or 'recent'."}`))
		return
	}

	if userId == "" || mode == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing id or mode parameters"}`))
		return
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "Accept: application/json",
		"Authorization": "Bearer " + os.Getenv("access_token"),
	}

	data, err := utils.Get(fmt.Sprintf("https://osu.ppy.sh/api/v2/users/%s/scores/%s?mode=%s&limit=100", userId, scoreType, mode), headers)

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

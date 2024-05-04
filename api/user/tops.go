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
	mode := strings.ToLower(r.URL.Query().Get("mode"))
	scoreType := strings.ToLower(r.URL.Query().Get("type"))

	if mode == "" {
		mode = "osu"
	}

	if scoreType != "best" && scoreType != "firsts" && scoreType != "recent" {
		utils.WriteError(w, "Type parameter must be one of 'best', 'firsts', or 'recent'.")
		return
	}

	if userId == "" {
		utils.WriteError(w, "Missing id parameter")
		return
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "Accept: application/json",
		"Authorization": "Bearer " + os.Getenv("access_token"),
	}

	data, err := utils.Get(fmt.Sprintf("https://osu.ppy.sh/api/v2/users/%s/scores/%s?mode=%s&limit=100", userId, scoreType, mode), headers)

	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Failed to fetch user data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

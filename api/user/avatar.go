package user

import (
	"fmt"
	"fun-api/utils"
	"net/http"
)

func Avatar(w http.ResponseWriter, r *http.Request) {
	avatarUrl := r.URL.Query().Get("url")

	if avatarUrl == "" {
		utils.WriteError(w, "`url` parameter must be defined.")
		return
	}

	headers := map[string]string{
		"Content-Type": "image/png",
	}

	dataBytes, err := utils.Get(avatarUrl, headers)

	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Failed to fetch user avatar: %s", err))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(dataBytes)
}

package handlers

import (
	"fun-api/utils"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, utils.ServeHtml(r, "index"))
}

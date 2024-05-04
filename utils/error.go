package utils

import (
	"fmt"
	"net/http"
)

func WriteError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Println(err)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}

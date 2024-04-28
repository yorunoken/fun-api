package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	resp := HelloObj{
		Message: "Hello!",
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	}
}

type HelloObj struct {
	Message string `json:"message"`
}

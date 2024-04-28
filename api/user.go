package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func User(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	mode := r.URL.Query().Get("mode")

	token := os.Getenv("token")

	if username == "" || mode == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing username or mode"}`))
		return
	}

	data, err := fetchUserData(username, mode, token)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to fetch user data"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func fetchUserData(username string, mode string, token string) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("https://osu.ppy.sh/api/get_user?u=%s&m=%s&k=%s", username, mode, token)

	req, err := http.NewRequest("GET", reqUrl, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userData []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	if len(userData) > 0 {
		return userData[0], nil
	}

	return nil, fmt.Errorf("no user data found")
}

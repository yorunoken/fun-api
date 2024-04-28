package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func SetToken(secret string, baseUrl string) {
	rq, err := Get(fmt.Sprintf("%s/api/token?secret=%s", baseUrl, secret))

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var tokenResponse TokenResponse

	err = json.Unmarshal(rq, &tokenResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Refreshed token.")
	os.Setenv("access_token", tokenResponse.AccessToken)
}

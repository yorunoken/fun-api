package user

import (
	"encoding/json"
	"fmt"
	"fun-api/utils"
	"math"
	"net/http"
	"os"
	"strings"
)

func Skills(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	mode := strings.ToLower(r.URL.Query().Get("mode"))
	baseUrl := os.Getenv("base_url")

	if mode == "" {
		mode = "osu"
	}

	bytes, err := utils.Get(fmt.Sprintf("%s/api/user/tops?id=%s&type=best", baseUrl, userId))

	if err != nil {
		utils.WriteError(w, fmt.Sprintf("There was an error while making the request to api/user/tops: %s", err))
		return
	}

	// Define an array of interface
	var tops []map[string]interface{}

	if err := json.Unmarshal(bytes, &tops); err != nil {
		utils.WriteError(w, fmt.Sprintf("There was an error while decoding JSON: %s", err))
		return
	}

	calcStandardSkills(tops)
}

func calcValue(val float64) float64 {
	factor := math.Pow(8.0/(val/72.0+8.0), 10)

	return -101.0*factor + 101.0
}

func calcStandardSkills(scores []map[string]interface{}) (float64, float64, float64) {
	// acc := 0.0
	// aim := 0.0
	// speed := 0.0
	// weightSum := 0.0

	const (
		accNerf   = 1.3
		aimNerf   = 2.6
		speedNerf = 2.4
	)

	for i, score := range scores {
		fmt.Println(score["max_combo"])
		fmt.Println(i)

		// state := OsuScoreState{
		// 	MaxCombo: score.MaxCombo,
		// 	N300:     score.Statistics.Great,
		// 	N100:     score.Statistics.Ok,
		// 	N50:      score.Statistics.Meh,
		// 	Misses:   score.Statistics.Miss,
		// }
		// fmt.Println(state)
	}

	return 0, 0, 0
}

func calcTaikoSkills() {

}

func calcFruitsSkils() {

}

func calcManiaSkills() {

}

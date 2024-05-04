package user

/*
#cgo LDFLAGS: -L../../lib -lrosu_pp_ffi
#include "../../lib/rosu_pp_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fun-api/utils"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"
)

func Skills(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	mode := strings.ToLower(r.URL.Query().Get("mode"))
	baseUrl := os.Getenv("base_url")

	if mode == "" {
		mode = "osu"
	}

	bytes, err := utils.Get(fmt.Sprintf("%s/api/user/tops?id=%s&type=best&mode=%s", baseUrl, userId, mode))
	if err != nil {
		utils.WriteError(w, fmt.Sprintf("There was an error while making the request to api/user/tops: %s", err))
		return
	}

	// Define an array of interface
	var tops []Score
	if err := json.Unmarshal(bytes, &tops); err != nil {
		utils.WriteError(w, fmt.Sprintf("There was an error while decoding JSON: %s", err))
		return
	}

	dbPath := "/root/HanamiBot/src/data.db"
	if os.Getenv("DEV") == "1" {
		dbPath = "./test.db"
	}

	var db *sql.DB

	for {
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			utils.WriteError(w, fmt.Sprintf("There was an error while opening the database: %s\nThis is most likely caused because the database is locked.\nRetrying until it succeeds.", err))
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	if err := db.Ping(); err != nil {
		utils.WriteError(w, fmt.Sprintf("Failed to connect to the database: %s", err))
		return
	}

	performances := []C.calculateresult{}

	for _, top := range tops {
		beatmapId := top.Beatmap.ID

		var data string
		ok := utils.EntryExists(db, "maps", fmt.Sprint(beatmapId))
		if ok {
			var id string
			utils.GetEntry(db, "maps", fmt.Sprint(beatmapId)).Scan(&id, &data)
		} else {
			fmt.Printf("beatmap_id %d does not exist\n", beatmapId)
			bytes, err := utils.Get(fmt.Sprintf("%s/api/beatmap/download?id=%d", baseUrl, beatmapId))
			if err != nil {
				utils.WriteError(w, fmt.Sprintf("There was an error while downloading beatmap number %d: %s", beatmapId, err))
				return
			}

			data = string(bytes)
			_, err = utils.AddEntry(db, "maps", fmt.Sprint(beatmapId), []utils.DatabaseData{{Key: "data", Value: data}})
			if err != nil {
				utils.WriteError(w, fmt.Sprintf("There was an error while inserting beatmap %d into database: %s", beatmapId, err))
				return
			}
		}

		statistics := top.Statistics
		params := CalculatorParams{
			mapData: string(data),
			scoreParams: ScoreParams{
				mode:    uint(top.ModeInt),
				mods:    uint(utils.GetModsEnum(top.Mods)),
				acc:     top.Accuracy,
				n300:    uint(statistics.Count300),
				n100:    uint(statistics.Count100),
				n50:     uint(statistics.Count50),
				nKatu:   nilOrValueUint(statistics.CountKatu),
				nGeki:   nilOrValueUint(statistics.CountGeki),
				nMisses: uint(statistics.CountMiss),
				combo:   uint(top.MaxCombo),
			},
		}

		performance := Calculate(params)

		performances = append(performances, performance)
	}

	db.Close()

	if mode == "osu" {
		acc, aim, speed := calcStandardSkills(performances)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"mode": "osu", "acc": "%f", "aim": "%f", "speed": "%f"}`, acc, aim, speed)))
	}

	if mode == "taiko" {
		acc, strain := calcTaikoSkills(performances)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"mode": "osu", "acc": "%f", "strain": "%f"}`, acc, strain)))
	}

	if mode == "fruits" {
		acc, movement := calcFruitsSkils(tops, performances)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"mode": "osu", "acc": "%f", "movement": "%f"}`, acc, movement)))
	}

	if mode == "mania" {
		acc, strain := calcManiaSkills(tops, performances)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"mode": "osu", "acc": "%f", "strain": "%f"}`, acc, strain)))
	}
}

func nilOrValueUint(value *int) uint {
	if value != nil {
		return uint(*value)
	}
	return 0
}

func calcStandardSkills(attrs []C.calculateresult) (float64, float64, float64) {
	var acc, aim, speed, weightSum float64

	const (
		AccNerf   = 1.3
		AimNerf   = 2.6
		SpeedNerf = 2.4
	)

	for i, attr := range attrs {
		accValue := attr.ppAcc.t / AccNerf
		aimValue := attr.ppAim.t / AimNerf
		speedValue := attr.ppSpeed.t / SpeedNerf
		weight := math.Pow(0.95, float64(i))

		acc += float64(accValue) * weight
		aim += float64(aimValue) * weight
		speed += float64(speedValue) * weight
		weightSum += weight
	}

	return acc / weightSum, aim / weightSum, speed / weightSum
}

func calcTaikoSkills(attrs []C.calculateresult) (float64, float64) {
	var acc, strain, weightSum float64

	const (
		AccNerf        = 1.15
		DifficultyNerf = 2.8
	)

	for i, attr := range attrs {
		accValue := attr.ppAcc.t / AccNerf
		difficultyValue := attr.ppDifficulty.t / DifficultyNerf
		weight := math.Pow(0.95, float64(i))

		acc += float64(accValue) * weight
		strain += float64(difficultyValue)
		weightSum += weight
	}

	return acc / weightSum, strain / weightSum
}

func calcFruitsSkils(scores []Score, attrs []C.calculateresult) (float64, float64) {
	var acc, movement, weightSum float64

	const (
		AccBuff      = 2.0
		MovementNerf = 4.7
	)

	for i, attr := range attrs {
		score := scores[i]

		od := attr.od
		acc_ := score.Accuracy
		nObjects := attr.nObjects.t

		accExp := math.Pow((math.Pow(acc_/46.5, 6) / 55.0), 1.5)
		accAdj := 1.0 / (5.0 * math.Log1p(accExp) * 0.1)

		accValue := math.Pow(float64(attr.stars), accExp-accAdj) * math.Pow(float64(od/7.0), 0.25) * math.Pow(float64(nObjects/2000.0), 0.15) * AccBuff
		movementValue := attr.pp / MovementNerf
		weight := math.Pow(0.95, float64(i))

		acc += accValue * weight
		movement += float64(movementValue) * weight
		weightSum += weight
	}

	return acc / weightSum, movement / weightSum
}

func calcManiaSkills(scores []Score, attrs []C.calculateresult) (float64, float64) {
	var acc, strain, weightSum float64

	const (
		AccBuff        = 2.1
		DifficultyNerf = 0.6
	)

	for i, attr := range attrs {
		score := scores[i]

		od := attr.od
		nObjects := attr.nObjects.t
		acc_ := math.Pow(math.Pow(score.Accuracy/60.0, 4.5), 1.5)

		accValue := math.Pow(float64(attr.stars), acc_) * math.Pow(float64(od/7.0), 0.25) * math.Pow(float64(nObjects/2000.0), 0.15) * AccBuff
		difficultyValue := attr.ppDifficulty.t / DifficultyNerf
		weight := math.Pow(0.95, float64(i))

		acc += accValue * weight
		strain += float64(difficultyValue) * weight
		weightSum += weight
	}

	return acc / weightSum, strain / weightSum
}

func Calculate(rosu CalculatorParams) C.calculateresult {
	cMapData := C.CString(rosu.mapData)
	defer C.free(unsafe.Pointer(cMapData))

	var calculator *C.calculator
	C.calculator_from_data(&calculator, cMapData)
	defer C.calculator_destroy(&calculator)

	var scoreParams *C.scoreparams
	C.score_params_new(&scoreParams)
	C.score_params_mode(scoreParams, C.mode(rosu.scoreParams.mode))
	if rosu.scoreParams.mods > 0 {
		C.score_params_mods(scoreParams, C.uint(rosu.scoreParams.mods))
	}
	if rosu.scoreParams.acc > 0 {
		C.score_params_acc(scoreParams, C.double(rosu.scoreParams.acc))
	}
	if rosu.scoreParams.n300 > 0 {
		C.score_params_n300(scoreParams, C.uint(rosu.scoreParams.n300))
	}
	if rosu.scoreParams.n100 > 0 {
		C.score_params_n100(scoreParams, C.uint(rosu.scoreParams.n100))
	}
	if rosu.scoreParams.n50 > 0 {
		C.score_params_n50(scoreParams, C.uint(rosu.scoreParams.n50))
	}
	if rosu.scoreParams.combo > 0 {
		C.score_params_combo(scoreParams, C.uint(rosu.scoreParams.combo))
	}
	if rosu.scoreParams.nMisses > 0 {
		C.score_params_n_misses(scoreParams, C.uint(rosu.scoreParams.nMisses))
	}
	if rosu.scoreParams.nKatu > 0 {
		C.score_params_n_katu(scoreParams, C.uint(rosu.scoreParams.nKatu))
	}
	defer C.score_params_destroy(&scoreParams)

	calculationResult := C.calculator_calculate(calculator, scoreParams)
	return calculationResult
}

type ScoreParams struct {
	mode          uint
	mods          uint
	acc           float64
	n300          uint
	n100          uint
	n50           uint
	nMisses       uint
	nKatu         uint
	nGeki         uint
	combo         uint
	passedObjects uint
	clockRate     float64
}

type CalculatorParams struct {
	mapData     string
	scoreParams ScoreParams
}

type Score struct {
	Beatmap    BeatmapCompact    `json:"beatmap"`
	Beatmapset BeatmapsetCompact `json:"beatmapset"`
	User       UserCompact       `json:"user"`
	ID         int               `json:"id"`
	BestID     int               `json:"best_id"`
	UserID     int               `json:"user_id"`
	Accuracy   float64           `json:"accuracy"`
	Mods       []string          `json:"mods"`
	Score      int               `json:"score"`
	MaxCombo   int               `json:"max_combo"`
	Perfect    bool              `json:"perfect"`
	Statistics ScoreStatistics   `json:"statistics"`
	Passed     bool              `json:"passed"`
	PP         float64           `json:"pp"`
	Rank       string            `json:"rank"`
	CreatedAt  string            `json:"created_at"`
	Mode       string            `json:"mode"`
	ModeInt    int               `json:"mode_int"`
	Replay     bool              `json:"replay"`
	Weight     Weight            `json:"weight"`
}

type ScoreStatistics struct {
	Count50   int  `json:"count_50"`
	Count100  int  `json:"count_100"`
	Count300  int  `json:"count_300"`
	CountGeki *int `json:"count_geki,omitempty"`
	CountKatu *int `json:"count_katu,omitempty"`
	CountMiss int  `json:"count_miss"`
}

type Weight struct {
	Percentage float64 `json:"percentage"`
	Pp         float64 `json:"pp"`
}

type BeatmapCompact struct {
	BeatmapsetID     int     `json:"beatmapset_id"`
	DifficultyRating float64 `json:"difficulty_rating"`
	ID               int     `json:"id"`
	Mode             string  `json:"mode"`
	Status           string  `json:"status"`
	TotalLength      int     `json:"total_length"`
	UserID           int     `json:"user_id"`
	Version          string  `json:"version"`
}

type BeatmapsetCompact struct {
	Artist         string `json:"artist"`
	ArtistUnicode  string `json:"artist_unicode"`
	Covers         Covers `json:"covers"`
	Creator        string `json:"creator"`
	FavouriteCount int    `json:"favourite_count"`
	ID             int    `json:"id"`
	NSFW           bool   `json:"nsfw"`
	PlayCount      int    `json:"play_count"`
	PreviewURL     string `json:"preview_url"`
	Source         string `json:"source"`
	Status         string `json:"status"`
	Title          string `json:"title"`
	TitleUnicode   string `json:"title_unicode"`
	UserID         int    `json:"user_id"`
	Video          bool   `json:"video"`
	Checksum       string `json:"checksum,omitempty"`
}

type Covers struct {
	Cover       string `json:"cover"`
	Cover2x     string `json:"cover@2x"`
	Card        string `json:"card"`
	Card2x      string `json:"card@2x"`
	List        string `json:"list"`
	List2x      string `json:"list@2x"`
	SlimCover   string `json:"slimcover"`
	SlimCover2x string `json:"slimcover@2x"`
}

type UserCompact struct {
	AvatarURL     string `json:"avatar_url"`
	CountryCode   string `json:"country_code"`
	DefaultGroup  string `json:"default_group"`
	ID            int    `json:"id"`
	IsActive      bool   `json:"is_active"`
	IsBot         bool   `json:"is_bot"`
	IsDeleted     bool   `json:"is_deleted"`
	IsOnline      bool   `json:"is_online"`
	IsSupporter   bool   `json:"is_supporter"`
	LastVisit     string `json:"last_visit"`
	PMFriendsOnly bool   `json:"pm_friends_only"`
	ProfileColour string `json:"profile_colour,omitempty"`
	Username      string `json:"username"`
}

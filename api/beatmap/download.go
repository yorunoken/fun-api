package beatmap

import (
	"fmt"
	"fun-api/utils"
	"net/http"
)

func Download(w http.ResponseWriter, r *http.Request) {
	mapId := r.URL.Query().Get("id")

	bytes, err := utils.Get(fmt.Sprintf("https://osu.ppy.sh/osu/%s", mapId))
	if err != nil {
		utils.WriteError(w, fmt.Sprintf("There was an error while downloading the beatmap: %s", err))
	}

	w.Write(bytes)
}

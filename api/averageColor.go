package api

import (
	"encoding/json"
	"fmt"
	"fun-api/utils"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func AverageColor(w http.ResponseWriter, r *http.Request) {
	imageUrl := r.URL.Query().Get("image")
	if imageUrl == "" {
		utils.WriteError(w, "`image` parameter was not specified.")
		return
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Error getting image: %s", err))
		return
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Error decoding image: %s", err))
		return
	}

	rgb := averageColorCalculator(img)

	respJson := responseType{Rgb: rgb}

	jsonBytes, err := json.Marshal(respJson)
	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Error on JSON: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

type responseType struct {
	Rgb [3]uint8 `json:"rgb"`
}

func averageColorCalculator(img image.Image) [3]uint8 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var totalR, totalG, totalB, totalPixels int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			currPixel := img.At(x, y)
			r, g, b, _ := currPixel.RGBA()

			// Shift rgb bits by 8
			totalR += int(r >> 8)
			totalG += int(g >> 8)
			totalB += int(b >> 8)
			totalPixels++
		}
	}

	avgR := uint8(totalR / totalPixels)
	avgG := uint8(totalG / totalPixels)
	avgB := uint8(totalB / totalPixels)

	return [3]uint8{avgR, avgG, avgB}
}

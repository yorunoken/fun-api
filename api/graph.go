package api

import (
	"bytes"
	"fmt"
	"fun-api/utils"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func Graph(w http.ResponseWriter, r *http.Request) {
	dataPointsStr := r.URL.Query().Get("points")
	isUpsideDown := r.URL.Query().Get("upside") == "true"

	if dataPointsStr == "" {
		utils.WriteError(w, "Missing points parameter")
		return
	}

	dataPoints, err := parseDataPoints(dataPointsStr)
	if err != nil {
		utils.WriteError(w, fmt.Sprintf("Invalid dataPoints parameter: %s", err))
		return
	}

	pts := generateData(dataPoints)

	graph := plot.New()

	graph.BackgroundColor = color.RGBA{42, 34, 38, 255}

	// Make them insivible
	graph.X.Tick.Color = color.RGBA{0, 0, 0, 0}
	graph.Y.Tick.Color = color.RGBA{0, 0, 0, 0}
	graph.X.Color = color.RGBA{0, 0, 0, 0}
	graph.Y.Color = color.RGBA{0, 0, 0, 0}
	graph.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{})
	graph.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{})

	margin := 10
	xMin := 0 - margin
	xMax := len(dataPoints) - 1 + margin
	yMin := utils.MinValue(dataPoints) - margin
	yMax := utils.MaxValue(dataPoints) + margin

	graph.X.Min = float64(xMin)
	graph.X.Max = float64(xMax)
	graph.Y.Min = float64(yMin)
	graph.Y.Max = float64(yMax)

	line, points, err := plotter.NewLinePoints(pts)
	if err != nil {
		log.Fatal(err)
	}

	line.Color = color.RGBA{255, 204, 34, 255}
	line.Width = 2
	points.Color = color.RGBA{0, 0, 0, 0}

	graph.Add(line, points)

	var buffer bytes.Buffer
	board, err := graph.WriterTo(vg.Length(8*vg.Inch), vg.Length(2*vg.Inch), "png")
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := board.WriteTo(&buffer); err != nil {
		fmt.Println(err)
		return
	}

	if isUpsideDown {
		img, err := png.Decode(&buffer)

		if err != nil {
			fmt.Println(err)
			return
		}

		flippedImage := flipImage(img)

		var flippedBuf bytes.Buffer
		if err := png.Encode(&flippedBuf, flippedImage); err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "image/png")
		if _, err := w.Write(flippedBuf.Bytes()); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if _, err := w.Write(buffer.Bytes()); err != nil {
		fmt.Println(err)
		return
	}
}

func generateData(n []int) plotter.XYs {

	pts := make(plotter.XYs, len(n))
	for i := range pts {
		x := float64(i)
		y := float64(n[i])
		pts[i].X = x
		pts[i].Y = y
	}
	return pts
}

func parseDataPoints(dataPointsStr string) ([]int, error) {
	strSlice := strings.Split(dataPointsStr, ",")
	dataPoints := make([]int, len(strSlice))
	for i, str := range strSlice {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		dataPoints[i] = val
	}
	return dataPoints, nil
}

func flipImage(img image.Image) image.Image {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			pixel := img.At(x, bounds.Dy()-y-1)
			flipped.Set(x, y, pixel)
		}
	}

	return flipped
}

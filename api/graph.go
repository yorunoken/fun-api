package api

import (
	"bytes"
	"fmt"
	"image/color"
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

	if dataPointsStr == "" {
		http.Error(w, "Missing points parameter", http.StatusBadRequest)
		return
	}

	dataPoints, err := parseDataPoints(dataPointsStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid dataPoints parameter: %v", err), http.StatusBadRequest)
		return
	}

	// reverse dataPoints array
	for i, j := 0, len(dataPoints)-1; i < j; i, j = i+1, j-1 {
		dataPoints[i], dataPoints[j] = dataPoints[j], dataPoints[i]
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

	line, points, err := plotter.NewLinePoints(pts)
	if err != nil {
		log.Fatal(err)
	}

	line.Color = color.RGBA{255, 204, 34, 255}
	line.Width = 2
	points.Color = color.RGBA{0, 0, 0, 0}

	graph.Add(line, points)

	var buf bytes.Buffer

	board, err := graph.WriterTo(vg.Length(7*vg.Inch), vg.Length(2*vg.Inch), "png")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := board.WriteTo(&buf); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Fatal(err)
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

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func max(a int, b int) int {

	if a > b {
		return a
	} else {
		return b
	}

}

func min(a int, b int) int {

	if a < b {
		return a
	} else {
		return b
	}

}

func getIdentityFilter() ([][]float64, float64, float64) {

	return [][]float64{
			{0, 0, 0},
			{0, 1, 0},
			{0, 0, 0},
		},
		0.0, //bias
		1.0 //factor
}

func getEdgeFilter1() ([][]float64, float64, float64) {

	return [][]float64{
			{1, 0, -1},
			{0, 0, 0},
			{-1, 0, 1},
		},
		0.0, //bias
		1.0 //factor
}

func getEdgeFilter2() ([][]float64, float64, float64) {

	return [][]float64{
			{0, 1, 0},
			{1, -4, 1},
			{0, 1, 0},
		},
		0.0, //bias
		1.0 //factor
}

func getEdgeFilter3() ([][]float64, float64, float64) {

	return [][]float64{
			{-1, -1, -1},
			{-1, 8, -1},
			{-1, -1, -1},
		},
		0.0, //bias
		1.0 //factor
}

func getEdgeFilter4() ([][]float64, float64, float64) {

	return [][]float64{
			{-1, 0, 0, 0, 0},
			{0, -2, 0, 0, 0},
			{0, 0, 6, 0, 0},
			{0, 0, 0, -2, 0},
			{0, 0, 0, 0, -1},
		},
		0.0, //bias
		1.0 //factor
}

func getEmboss() ([][]float64, float64, float64) {

	return [][]float64{
			{-1, -1, 0},
			{-1, 0, 1},
			{0, 1, 1},
		},
		50.0, //bias
		1.0 //factor
}

func getEmboss2() ([][]float64, float64, float64) {

	return [][]float64{
			{-1, -1, -1, -1, 0},
			{-1, -1, -1, 0, 1},
			{-1, -1, 0, 1, 1},
			{-1, 0, 1, 1, 1},
			{0, 1, 1, 1, 1},
		},
		0.0, //bias
		1.0 //factor
}

func getBlur2() ([][]float64, float64, float64) {

	return [][]float64{
			{0, 0, 1, 0, 0},
			{0, 1, 1, 1, 0},
			{1, 1, 1, 1, 1},
			{0, 1, 1, 1, 0},
			{0, 0, 1, 0, 0},
		},
		0.0, //bias
		1.0 / 16 //factor
}

func getBlur() ([][]float64, float64, float64) {

	return [][]float64{
			{0.0625, 0.125, 0.0625},
			{0.125, 0.25, 0.125},
			{0.0625, 0.125, 0.0625},
		},
		8.0, //bias
		1 //factor
}

func getExcessiveEdge() ([][]float64, float64, float64) {

	return [][]float64{
			{1, 1, 1},
			{1, -7, 1},
			{1, 1, 1},
		},
		1.0, //bias
		1.0 //factor
}

func getEdge() ([][]float64, float64, float64) {

	return [][]float64{
			{-1, -1, -1},
			{0, 0, 0},
			{1, 1, 1},
		},
		100.0, //bias
		1.0 //factor
}

func getSharpen() ([][]float64, float64, float64) {

	return [][]float64{
			{-1, -1, -1, -1, -1},
			{-1, 2, 2, 2, -1},
			{-1, 2, 8, 2, -1},
			{-1, 2, 2, 2, -1},
			{-1, -1, -1, -1, -1},
		},
		1.0, //bias
		1.0 / 4 //factor
}

type filter func() ([][]float64, float64, float64)

func main() {

	m, w, h := getImageArray("lena.png")

	f := []filter{
		getIdentityFilter,
		getEdgeFilter1,
		getEdgeFilter2,
		getEdgeFilter3,
		getEdgeFilter4,
		getEmboss,
		getEmboss2,
		getBlur,
		getBlur2,
		getEdge,
		getExcessiveEdge,
		getSharpen,
	}

	for idx, method := range f {

		filter, bias, factor := method()

		fmt.Println(filter, bias, factor)

		stride := len(filter)

		result := applyFilter(w, h, m, filter, stride, factor, bias)

		outfile := fmt.Sprintf("output_%d.png", idx)
		img, _ := os.Create(outfile)
		defer img.Close()
		png.Encode(img, result)

	}

}

func getImageArray(imagePath string) (image.Image, int, int) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	reader, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	return m, w, h
}

func applyFilter(w int, h int, m image.Image, filter [][]float64, stride int, factor float64, bias float64) *image.RGBA {

	result := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {

			var red float64 = 0.0
			var green float64 = 0.0
			var blue float64 = 0.0

			for filterY := 0; filterY < stride; filterY++ {
				for filterX := 0; filterX < stride; filterX++ {
					imageX := (x - stride/2 + filterX + w) % w
					imageY := (y - stride/2 + filterY + h) % h
					r, g, b, _ := m.At(imageX, imageY).RGBA()

					red += (float64(r) / 257) * filter[filterY][filterX]
					green += (float64(g) / 257) * filter[filterY][filterX]
					blue += (float64(b) / 257) * filter[filterY][filterX]
				}
			}

			_r := min(max(int(factor*red+bias), 0), 255)
			_g := min(max(int(factor*green+bias), 0), 255)
			_b := min(max(int(factor*blue+bias), 0), 255)

			c := color.RGBA{
				uint8(_r),
				uint8(_g),
				uint8(_b),
				255,
			}
			result.Set(x, y, c)

		}

	}
	return result
}

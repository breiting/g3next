package generation

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"

	"github.com/breiting/g3next/noise"
	"github.com/g3n/engine/math32"
)

type NoiseMap struct {
	data   [][]float32
	Width  int
	Height int
}

func NewNoiseMap(seed int64, width, height int, ofs, scale float64, octaves int, persistance, lacunarity float64) NoiseMap {

	min := math.MaxFloat64
	max := -math.MaxFloat64

	// seed := int64(1587046793530293277)
	// offset := [2]float64{1, 1}
	offset := [2]float64{ofs, 0}

	// NewPerlin creates new Perlin noise generator
	// In what follows "alpha" is the weight when the sum is formed.
	// Typically it is 2, As this approaches 1 the function is noisier.
	// "beta" is the harmonic scaling/spacing, typically 2, n is the
	// number of iterations and seed is the math.rand seed value to use
	perlinEngine := noise.NewPerlin(2, 2, 3, seed)

	noiseMap := NoiseMap{
		Width:  width,
		Height: height,
	}

	noiseMap.data = make([][]float32, height)
	for i := 0; i < height; i++ {
		noiseMap.data[i] = make([]float32, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			amplitude := float64(1)
			frequency := float64(1)
			noiseHeight := float64(0)

			for i := 0; i < octaves; i++ {
				perlinValue := perlinEngine.Noise2D(
					float64(x)/scale*frequency+offset[0],
					float64(y)/scale*frequency+offset[1])
				noiseHeight += perlinValue * amplitude

				amplitude *= persistance
				frequency += lacunarity
			}
			if noiseHeight < min {
				min = noiseHeight
			}
			if noiseHeight > max {
				max = noiseHeight
			}

			noiseMap.data[x][y] = float32(noiseHeight)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			noiseMap.data[x][y] = normalize(float32(min), float32(max), noiseMap.data[x][y])
		}
	}

	return noiseMap
}

func (n *NoiseMap) Get(x, y int) float32 {
	// TODO check boundaries
	return n.data[x][y]
}

func (n *NoiseMap) GetColor(x, y int) color.RGBA {
	return getColor(n.data[x][y])
}

func (n *NoiseMap) WriteImage(w io.Writer) error {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{n.Width, n.Height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for y := 0; y < n.Height; y++ {
		for x := 0; x < n.Width; x++ {
			img.Set(x, y, getColor(n.data[x][y]))
		}
	}
	png.Encode(w, img)
	return nil
}

func interpolateColor(v float32, r1, g1, b1, r2, g2, b2 float32) color.RGBA {

	sample := math32.Color{
		R: r1 / 255.0,
		G: g1 / 255.0,
		B: b1 / 255.0,
	}
	sample.Lerp(&math32.Color{
		R: r2 / 255.0,
		G: g2 / 255.0,
		B: r2 / 255.0,
	}, v)

	return color.RGBA{
		R: uint8(sample.R * 255.0),
		G: uint8(sample.G * 255.0),
		B: uint8(sample.B * 255.0),
		A: 0xff,
	}
}

func getColor(v float32) color.RGBA {

	if v < 0.4 {
		return interpolateColor(v, 14, 0, 100, 0, 51, 100)
	}
	if v < 0.45 {
		// sand
		return interpolateColor(v, 100, 80, 0, 77, 100, 0)
	}
	if v < 0.6 {
		return interpolateColor(v, 0, 48, 23, 0, 95, 44)
	}
	if v < 0.9 {
		return interpolateColor(v, 50, 26, 20, 26, 21, 20)
	}
	// snow
	return interpolateColor(v, 72, 72, 72, 255, 255, 255)
}

// normalize maps the value from range [min..max] to [0..1]
func normalize(min, max, val float32) float32 {
	return (val - min) / (max - min)
}

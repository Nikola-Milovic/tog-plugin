package engine

import (
	"github.com/Nikola-Milovic/tog-plugin/math"
	rand2 "math/rand"
)

type CurveFunction func(dist, maxDist int, value float32) float32

type Imap struct {
	Height   int
	Width    int
	TileSize int
	Grid     [][]float32
}

func NewImap(width, height, tileSize int) *Imap {
	imap := Imap{Height: height, Width: width, TileSize: tileSize}
	grid := make([][]float32, width)
	imap.Grid = grid

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			imap.SetCell(x, y, 0)
		}
	}

	return &imap
}

func (imap *Imap) Normalize() {
	highest := float32(-10000)
	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			if imap.Grid[x][y] > highest {
				highest = imap.Grid[x][y]
			}
		}
	}
	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			if highest != 0 {
				imap.Grid[x][y] /= highest
			}
		}
	}
}

func (imap *Imap) NormalizeAndInvert() {
	highest := float32(-10000)
	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			if imap.Grid[x][y] > highest {
				highest = imap.Grid[x][y]
			}
		}
	}
	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			if highest != 0 {
				imap.Grid[x][y] /= highest
			}
			imap.Grid[x][y] = 1 - imap.Grid[x][y]
		}
	}
}

func (imap *Imap) SetCell(x, y int, newVal float32) {
	if x >= 0 && x < imap.Width && y >= 0 && y < imap.Height {
		if imap.Grid[x] == nil {
			imap.Grid[x] = make([]float32, imap.Height)
		}
		imap.Grid[x][y] = newVal
	}
}
func (imap *Imap) GetCell(x, y int) float32 {
	if x >= 0 && x < imap.Width && y >= 0 && y < imap.Height {
		return imap.Grid[x][y]
	} else {
		return 0.0
	}
}

func (imap *Imap) AddValue(x, y int, value float32) {
	if x >= 0 && x < imap.Width && y >= 0 && y < imap.Height {
		imap.Grid[x][y] += value
	}
}

func (imap *Imap) PropagateInfluence(centerX, centerY, radius int, influenceCalc CurveFunction, magnitude float32) {
	startX := centerX - (radius / 2)
	startY := centerY - (radius / 2)
	endX := centerX + (radius / 2)
	endY := centerY + (radius / 2)

	minX := math.Min(0, startX)
	maxX := math.Max(imap.Width, endX)
	minY := math.Min(0, startY)
	maxY := math.Max(imap.Height, endY)

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			distance := math.GetDistanceBetweenPoints(math.P(x, y), math.P(centerX, centerY))
			imap.Grid[x][y] = influenceCalc(distance, radius/2, 1.0) * magnitude
		}
	}
}

func (imap *Imap) GetHighestCell() (x, y int, value float32) {
	highValue := float32(0.0)
	highCellX := 0
	highCellY := 0

	offsetX := 0
	offsetY := 0

	offsetX += rand2.Intn(imap.Width-0) + 0
	offsetY += rand2.Intn(imap.Height-0) + 0

	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			actualX := (x + offsetX) % imap.Width
			actualY := (y + offsetY) % imap.Height

			value := imap.GetCell(actualX, actualY)
			if value > highValue {
				highCellX = actualX
				highCellY = actualY
				highValue = value
			}

		}
	}

	return highCellX, highCellY, highValue
}

func (imap *Imap) GetLowestValue() (x, y int, value float32) {
	//rand := rand2.Rand{}
	//rand.Seed(time.Now().UnixNano())

	lowValue := float32(10000.0)
	lowCellX := 0
	lowCellY := 0

	offsetX := 0
	offsetY := 0

	offsetX += rand2.Intn(imap.Width-0) + 0
	offsetY += rand2.Intn(imap.Width-0) + 0

	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			actualX := (x + offsetX) % imap.Width
			actualY := (x + offsetY) % imap.Height

			value := imap.GetCell(actualX, actualY)
			if value < lowValue {
				lowCellX = actualX
				lowCellY = actualY
				lowValue = value
			}

		}
	}

	return lowCellX, lowCellY, lowValue
}

func (imap *Imap) Clear() {
	for x := 0; x < imap.Width; x++ {
		for y := 0; y < imap.Height; y++ {
			imap.Grid[x][y] = 0.0
		}
	}
}

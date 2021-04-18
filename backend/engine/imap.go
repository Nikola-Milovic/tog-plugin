package engine

type CurveFunction func(dist, maxDist int, value float32) float32

type Imap struct {
	Height   int
	Width    int
	TileSize int
	Grid     [][]float32
}

func NewImap(width, height, tileSize int) Imap {
	imap := Imap{Height: height, Width: width, TileSize: tileSize}
	grid := make([][]float32, width)
	imap.Grid = grid

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			imap.SetCell(x, y, 0)
		}
	}

	return imap
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
			imap.Grid[x][y] /= highest
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
	for x := 0; x < imap.Width/imap.TileSize; x++ {
		for y := 0; y < imap.Height; y++ {
			imap.Grid[x][y] /= highest
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

	minX := Min(0, startX)
	maxX := Max(imap.Width, endX)
	minY := Min(0, startY)
	maxY := Max(imap.Height, endY)

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			distance := GetDistanceIncludingDiagonal(x, y, centerX, centerY)
			imap.Grid[x][y] = influenceCalc(distance, radius/2, 1.0) * magnitude
		}
	}
}

package engine

func AddIntoBiggerMap(src, dest *Imap, centerX, centerY int, magnitude float32) {
	startX := centerX - (src.Width / 2)
	startY := centerY - (src.Height / 2)

	for y := 0; y < src.Height; y++ {
		for x := 0; x < src.Width; x++ {
			targetX := x + startX
			targetY := y + startY

			if targetX >= 0 && targetY >= 0 && targetX < dest.Width &&
				targetY < dest.Height {
				dest.Grid[targetX][targetY] += src.GetCell(x, y) * magnitude
			}
		}
	}
}

func AddIntoSmallerMap(src, dest *Imap, centerX, centerY int, magnitude float32) {
	tarMapWidth := dest.Width
	tarMapHeight := dest.Height

	for y := 0; y < tarMapWidth; y++ {
		for x := 0; x < tarMapHeight; x++ {
			targetX := centerX - tarMapWidth/2 + x
			targetY := centerY - tarMapHeight/2 + y

			dest.AddValue(x, y, src.GetCell(targetX, targetY)*magnitude)
		}
	}
}

func TranslateCoordsOutsideofMapIntoMap(centerX, centerY, tarX, tarY, size int) (newX, newY int) {
	x := 0
	y := 0

	offset := size / 2

	if tarX > centerX+offset {
		x = size
	} else if tarX < centerX-offset {
		x = 0
	} else {
		x = size/2
	}

	if tarY > centerY+offset {
		y = size
	} else if tarY < centerY-offset {
		y = 0
	} else {
		y = size/2
	}

	return x, y
}


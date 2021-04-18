package engine


func AddIntoBiggerMap(src, dest Imap, centerX, centerY int, magnitude float32) Imap {
	startX := centerX - ( src.Width / 2)
	startY := centerY  - ( src.Height / 2)

	for y := 0; y < src.Height; y ++ {
		for x := 0; x < src.Width; x++ {
			targetX := x + startX
			targetY := y + startY

			if targetX >=0 && targetY >= 0 && targetX < dest.Width &&
				targetY < dest.Height {
				dest.Grid[targetX][targetY] += src.GetCell(x,y) * magnitude
			}
		}
	}
	return dest
}

func AddIntoSmallerMap(src, dest Imap, centerX, centerY int, magnitude float32) Imap {
	tarMapWidth := dest.Width
	tarMapHeight := dest.Height

	//startX := centerX  - ( tarMapWidth >> 1)
	//startY := centerY - ( tarMapHeight >> 1)

	//minX := Max(0,startX)
	//maxX := Min(tarMapWidth, src.Width - startX)
	//minY := Max(0, startY)
	//maxY := Min(tarMapHeight, src.Height - startY)


	for y := 0; y < tarMapWidth; y ++ {
		for x := 0; x < tarMapHeight; x++ {
			targetX := centerX -tarMapWidth/2 + x
			targetY := centerY -tarMapHeight/2 + y

			dest.AddValue(x, y, src.GetCell(targetX, targetY) * magnitude)
		}
	}
	return dest
}
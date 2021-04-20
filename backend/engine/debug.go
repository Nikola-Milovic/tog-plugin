package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintImapToFile(imap *Imap, title string, append bool) {
	f, err := os.OpenFile("temp.txt", os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	if !append {
		f, err = os.Create("../tests/temp.txt")
	}

	defer f.Close()

	var sb strings.Builder
	heading := fmt.Sprintf("------------ %s ----------- \n width %d, height %d \n", title, imap.Width, imap.Height)
	sb.WriteString(heading)
	for y := 0; y < imap.Height; y++ {
		for x := 0; x < imap.Width; x++ {
			s := fmt.Sprintf("%.2f ", imap.Grid[x][y])
			sb.WriteString(s)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	w := bufio.NewWriter(f)
	writtenBytes, err := w.WriteString(sb.String())
	check(err)
	fmt.Printf("wrote %d bytes\n", writtenBytes)

	w.Flush()
}


func PrintImapToFileWithStep(imap *Imap, title string, step int) {
	f, err := os.OpenFile("../tests/temp.txt", os.O_APPEND|os.O_WRONLY, 0644)
	check(err)

	defer f.Close()

	var sb strings.Builder
	heading := fmt.Sprintf("------------ %s ----------- \n width %d, height %d \n", title, imap.Width, imap.Height)
	sb.WriteString(heading)
	for y := 0; y < imap.Height; y += step {
		for x := 0; x < imap.Width; x += step {
			if imap.Grid[x][y] > 0.0 {
				s := "* "
				sb.WriteString(s)
			} else {
				s := ". "      //fmt.Sprintf("%.2f ", imap.Grid[x][y])
				sb.WriteString(s)
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	w := bufio.NewWriter(f)
	writtenBytes, err := w.WriteString(sb.String())
	check(err)
	fmt.Printf("wrote %d bytes\n", writtenBytes)

	w.Flush()
}



func check(e error) {
	if e != nil {
		panic(e)
	}
}
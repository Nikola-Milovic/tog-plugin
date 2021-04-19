package engine

import "fmt"

type ImapTemplate struct {
	Radius int
	Type   string
	Imap   *Imap
}

func (imap *ImapTemplate) String() string {
	return fmt.Sprintf("Radius : %d, type : %s, Imap size : %dx%d", imap.Radius, imap.Type, imap.Imap.Width, imap.Imap.Height)
}


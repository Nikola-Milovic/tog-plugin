package main

import "log"

func check(err error) {
	if err == nil {
		return
	}
	log.Fatalln(err)
}
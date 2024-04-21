package main

import (
	"log"
)

func main() {
	lp := ":7777"

	db, err := NewStorage()
	if err != nil {
		log.Fatalln(err)
	}
	s, err := NewAPIServer(lp, db)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
	//select {}
}

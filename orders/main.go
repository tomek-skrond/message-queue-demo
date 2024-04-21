package main

import (
	"log"
)

func main() {

	listenPort := ":8000"
	db, err := NewStorage()
	if err != nil {
		log.Fatalln(err)
	}
	s, err := NewAPIServer(listenPort, db)
	if err != nil {
		log.Fatalln()
	}

	go s.Start()

	select {}
}

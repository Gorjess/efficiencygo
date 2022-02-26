package main

import (
	"github.com/Gorjess/kitten/client/klient"
	"log"
)

func main() {
	err := klient.New().Run()
	if err != nil {
		log.Println(err)
	}
}

package main

import (
	"github.com/Gorjess/kitten/echo-server/srv"
	"log"
)

func main() {
	tes, err := srv.New()
	if err != nil {
		log.Fatalln("new typical echo echo-server failed:",
			err.Error())
		return
	}
	// run echo-server
	tes.Run()
}

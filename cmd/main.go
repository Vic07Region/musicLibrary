package main

import (
	"github.com/Vic07Region/musicLibrary/internal/pkg/app"
	"log"
)

func main() {
	//init app
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	//run app
	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}

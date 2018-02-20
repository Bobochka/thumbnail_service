package main

import (
	"log"

	"net/http"

	"github.com/Bobochka/thumbnail_service/lib/service"
)

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(cfg)
	app := &App{service: svc}

	http.HandleFunc("/thumbnail", app.thumbnail)

	log.Fatal(http.ListenAndServe(bindPort(), nil))
}

package main

import (
	"net/http"
	"orc-server/orc"
)

type Config struct {
	Tickcount int `json:"tickcount"`
}

var GlobalConfig Config

func main() {
	server := orc.GameServer{
		Port: 1004,
	}
	go server.Run()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8282", nil)
}

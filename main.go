package main

import (
	"fmt"
	"net/http"
	"orc-server/orc"
)

func main() {
	fmt.Println("config info")
	fmt.Printf("%+v\n", orc.GlobalConfig)

	orc.StartGlobal()

	server := orc.GameServer{
		Port: 8000,
	}
	go server.Run()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8282", nil)
}

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
		Port: orc.GlobalConfig.ServerPort,
	}
	go server.Run()
	fmt.Println("listen on port : ", orc.GlobalConfig.ServerPort)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8282", nil)
}

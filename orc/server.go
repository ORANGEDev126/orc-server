package orc

import (
	"fmt"
	"net"
	"strconv"
)

type GameServer struct {
	Port int
}

func (server *GameServer) Run() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		sock, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		session := NewSession(sock)
		session.playGround = &globalPlayGround
		RegisterGlobal(session)
	}
}

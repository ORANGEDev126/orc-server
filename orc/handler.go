package orc

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

func HandleMoveJogReq(session *Session, buf []byte) {
	msg := &MoveJogReqMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Println("unmarshal error on move jog req %v", err)
		return
	}

	playGround := session.playGround
	if playGround == nil {
		fmt.Println("play ground is nil on move jog req")
		return
	}

	moveJog := MoveJogChan{
		id:  session.id,
		dir: msg.Dir,
	}

	playGround.moveJog <- moveJog
}

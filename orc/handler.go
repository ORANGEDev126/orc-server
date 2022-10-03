package orc

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

func HandleMoveJogReq(session *Session, buf []byte) {
	msg := &MoveJogReqMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Println("unmarshal error on move jog req", err)
		return
	}

	playGround := session.playGround
	if playGround == nil {
		fmt.Println("play ground is nil on move jog req")
		return
	}

	moveJog := MoveJogChan{
		id:  session.GetId(),
		dir: msg.Dir,
	}

	playGround.moveJog <- moveJog
}

func HandleShootProjectileReq(session *Session, buf []byte) {
	msg := &ShootProjectileReqMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Println("unmarshal error on shoot projectile req", err)
		return
	}

	playGround := session.playGround
	if playGround == nil {
		fmt.Println("play ground is nil on shoot projectile req")
		return
	}

	reqChannel := ShootProjectileChan{
		playerId: session.GetId(),
		angle:    int(msg.GetAngle()),
	}

	playGround.shoot <- reqChannel
}

func HandleAttackReq(session *Session, buf []byte) {
	msg := &AttackReqMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Println("unmarshal error on attack projectile req", err)
		return
	}

	playGround := session.playGround
	if playGround == nil {
		fmt.Println("play ground is nil on attack")
		return
	}

	playGround.attackChan <- AttackChan{
		playerId: session.GetId(),
	}
}

func HandleDefenceReq(session *Session, buf []byte) {
	msg := &DefenceReqMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Println("unmarshal error on defence req", err)
		return
	}

	playGround := session.playGround
	if playGround == nil {
		fmt.Println("play ground is nill on defence req")
		return
	}

	playGround.defenceChan <- DefenceChan{
		playerId: session.GetId(),
	}
}

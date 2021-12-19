package orc

import (
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	session *Session
	currDir Direction
	x       float32
	y       float32
	speed   float32
	accel   float32
}

func (player Player) ToPlayerMessage() *PlayerMessage {
	return &PlayerMessage{
		Id: int64(player.GetId()),
		X:  player.x,
		Y:  player.y,
	}
}

func NewPlayer(s *Session) *Player {
	return &Player{
		session: s,
		currDir: Direction_NONE_DIR,
	}
}

func (player *Player) GetId() uint64 {
	return player.session.id
}

func (player *Player) GetLocation() (float32, float32) {
	return player.x, player.y
}

func (player *Player) SendMessage(id Protocol, msg proto.Message) {
	packetLen := HEADER_LENGTH + proto.Size(msg)
	buf := make([]byte, packetLen)

	binary.LittleEndian.PutUint32(buf, uint32(packetLen))
	binary.LittleEndian.PutUint32(buf, uint32(id))
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("marshal message error")
		return
	}

	copy(buf[HEADER_LENGTH:], msgBytes)
	player.session.Send(buf)
}

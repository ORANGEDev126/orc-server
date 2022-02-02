package orc

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	session  *Session
	currDir  Direction
	jogDir   Direction
	circle   Circle
	speed    float64
	maxSpeed float64
	accel    float64
}

func (player Player) ToPlayerMessage() *PlayerMessage {
	return &PlayerMessage{
		Id: int64(player.GetId()),
		X:  player.circle.point.x,
		Y:  player.circle.point.y,
	}
}

func NewPlayer(s *Session) *Player {
	return &Player{
		session: s,
		currDir: Direction_NONE_DIR,
		jogDir:  Direction_NONE_DIR,
		circle:  Circle{Point{0, 0}, GlobalConfig.PlayerRadius},
	}
}

func (player *Player) GetId() uint64 {
	return player.session.id
}

func (player *Player) SendMessage(id Protocol, msg proto.Message) {
	packetLen := HEADER_LENGTH + proto.Size(msg)
	buf := make([]byte, packetLen)

	binary.LittleEndian.PutUint32(buf, uint32(packetLen))
	binary.LittleEndian.PutUint32(buf[LENGTH_FIELD_LENGTH:], uint32(id))
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("marshal message error")
		return
	}

	copy(buf[HEADER_LENGTH:], msgBytes)
	player.session.Send(buf)
}

func (player *Player) ProjectileAttacked(projectileAngle int) {

}

func (player *Player) NextSpeed() float64 {
	accel := player.accel
	if player.jogDir == Direction_NONE_DIR {
		accel = -accel
	}

	v := player.speed + accel
	return Clamp(0, player.maxSpeed, v)
}

func (player *Player) NextDirection() Direction {
	if player.jogDir == Direction_NONE_DIR {
		return player.currDir
	}

	if player.currDir == Direction_NONE_DIR {
		return player.jogDir
	}

	diff := GetDiffDirection(player.currDir, player.jogDir)
	if diff == 0 {
		return player.currDir
	} else if diff == 4 {
		if rand.Int()%2 == 0 {
			return DirectionToClockwise(player.currDir)
		} else {
			return DirectionToAntiClockwise(player.currDir)
		}
	} else if diff < 4 {
		return DirectionToClockwise(player.currDir)
	} else {
		return DirectionToAntiClockwise(player.currDir)
	}
}

func (player *Player) NextPoint(nextSpeed float64, nextDirection Direction) Point {
	diagonalVal := nextSpeed * math.Sqrt2 / float64(2)

	if nextDirection == Direction_NORTH {
		return Point{player.circle.point.x + nextSpeed,
			player.circle.point.y}
	} else if nextDirection == Direction_NORTH_EAST {
		return Point{player.circle.point.x + diagonalVal,
			player.circle.point.y + diagonalVal}
	} else if nextDirection == Direction_EAST {
		return Point{player.circle.point.x + nextSpeed,
			player.circle.point.y}
	} else if nextDirection == Direction_EAST_SOUTH {
		return Point{player.circle.point.x + diagonalVal,
			player.circle.point.y - diagonalVal}
	} else if nextDirection == Direction_SOUTH {
		return Point{player.circle.point.x,
			player.circle.point.y - nextSpeed}
	} else if nextDirection == Direction_SOUTH_WEST {
		return Point{player.circle.point.x - diagonalVal,
			player.circle.point.y - diagonalVal}
	} else if nextDirection == Direction_WEST {
		return Point{player.circle.point.x - nextSpeed,
			player.circle.point.y}
	} else if nextDirection == Direction_WEST_NORTH {
		return Point{player.circle.point.x - diagonalVal,
			player.circle.point.y + diagonalVal}
	}

	return player.circle.point
}

func (player *Player) Move(nextSpeed float64,
	nextDirection Direction,
	point Point) {
	player.speed = nextSpeed
	player.currDir = nextDirection
	player.circle.point = point
}

func (player *Player) ToMoveObjectMessage() *MoveObjectNotiMessage_Object {
	return &MoveObjectNotiMessage_Object{
		Id:  int64(player.GetId()),
		X:   player.circle.point.x,
		Y:   player.circle.point.y,
		Dir: player.currDir,
	}
}

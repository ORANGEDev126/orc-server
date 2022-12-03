package orc

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	session        *Session
	currDir        Direction
	jogDir         Direction
	circle         Circle
	speed          float64
	maxSpeed       float64
	accel          float64
	attackRange    int
	attackDistance float64
	defenceTime    time.Time
	attackedTime   time.Time
	hp             int
}

func (player Player) ToPlayerMessage() *PlayerMessage {
	return &PlayerMessage{
		Id:       int64(player.GetId()),
		X:        player.circle.point.x,
		Y:        player.circle.point.y,
		RemainHp: int32(player.GetHP()),
		MaxHp:    int32(GlobalConfig.PlayerMaxHP),
	}
}

func NewPlayer(s *Session) *Player {
	return &Player{
		session:        s,
		currDir:        Direction_NONE_DIR,
		jogDir:         Direction_NONE_DIR,
		circle:         Circle{Point{0, 0}, GlobalConfig.PlayerRadius},
		accel:          GlobalConfig.Accel,
		maxSpeed:       GlobalConfig.MaxSpeed,
		attackRange:    GlobalConfig.PlayerAttackRange,
		attackDistance: GlobalConfig.PlayerAttackDistance,
		hp:             300,
	}
}

func (player *Player) GetId() uint64 {
	return player.session.id
}

func (player *Player) IsDefencing() bool {
	return int(time.Now().Sub(player.defenceTime).Milliseconds()) < GlobalConfig.DefenceDuration
}

func (player *Player) IsStiffen() bool {
	return int(time.Now().Sub(player.attackedTime).Milliseconds()) < GlobalConfig.DefenceDuration
}

func (player *Player) SendMessage(id Notification, msg proto.Message) {
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

func (player *Player) CurrSpeed() float64 {
	return player.speed
}

func (player *Player) CurrPoint() Point {
	return player.circle.point
}

func (player *Player) UpdateSpeed(deltaTime int64) float64 {
	accel := player.accel
	if player.jogDir == Direction_NONE_DIR {
		accel = -accel
	}

	delta := accel / GlobalConfig.PhysicsTickCount * float64(deltaTime)
	v := Clamp(0, player.maxSpeed, player.speed+delta)

	player.speed = v

	return v
}

func (player *Player) UpdateHP(delta int) {
	player.hp += delta
}

func (player *Player) GetHP() int {
	return player.hp
}

func (player *Player) UpdateDirection() Direction {
	if player.jogDir == Direction_NONE_DIR {
		return player.currDir
	}

	if player.currDir == Direction_NONE_DIR {
		player.currDir = player.jogDir
		return player.jogDir
	}

	diff := GetDiffDirection(player.currDir, player.jogDir)
	dir := player.currDir

	if diff == 0 {
		return player.currDir
	} else if diff == 4 {
		if rand.Int()%2 == 0 {
			dir = DirectionToClockwise(player.currDir)
		} else {
			dir = DirectionToAntiClockwise(player.currDir)
		}
	} else if diff < 4 {
		dir = DirectionToClockwise(player.currDir)
	} else {
		dir = DirectionToAntiClockwise(player.currDir)
	}

	player.currDir = dir

	return dir
}

func (player *Player) NextPoint(nextSpeed float64, nextDirection Direction, deltaTime int64) Point {
	midSpeed := (player.speed + nextSpeed) / float64(2)
	straightVal := midSpeed * float64(deltaTime) / float64(GlobalConfig.PhysicsTickCount)
	diagonalVal := straightVal * math.Sqrt2 / float64(2)

	if nextDirection == Direction_NORTH {
		return Point{player.circle.point.x,
			player.circle.point.y + straightVal}
	} else if nextDirection == Direction_NORTH_EAST {
		return Point{player.circle.point.x + diagonalVal,
			player.circle.point.y + diagonalVal}
	} else if nextDirection == Direction_EAST {
		return Point{player.circle.point.x + straightVal,
			player.circle.point.y}
	} else if nextDirection == Direction_EAST_SOUTH {
		return Point{player.circle.point.x + diagonalVal,
			player.circle.point.y - diagonalVal}
	} else if nextDirection == Direction_SOUTH {
		return Point{player.circle.point.x,
			player.circle.point.y - straightVal}
	} else if nextDirection == Direction_SOUTH_WEST {
		return Point{player.circle.point.x - diagonalVal,
			player.circle.point.y - diagonalVal}
	} else if nextDirection == Direction_WEST {
		return Point{player.circle.point.x - straightVal,
			player.circle.point.y}
	} else if nextDirection == Direction_WEST_NORTH {
		return Point{player.circle.point.x - diagonalVal,
			player.circle.point.y + diagonalVal}
	}

	return player.circle.point
}

func (player *Player) UpdatePoint(point Point) {
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

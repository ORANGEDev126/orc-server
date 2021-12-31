package orc

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
)

var globalPlayGround PlayGround

type MoveJogChan struct {
	id  uint64
	dir Direction
}

type PlayGround struct {
	players    map[uint64]*Player
	register   chan *Player
	unregister chan *Session
	stop       chan bool
	moveJog    chan MoveJogChan
}

func StartGlobal() {
	globalPlayGround = PlayGround{
		players:    map[uint64]*Player{},
		register:   make(chan *Player),
		unregister: make(chan *Session),
		stop:       make(chan bool),
		moveJog:    make(chan MoveJogChan),
	}

	go globalPlayGround.eventLoop()
}

func RegisterGlobal(session *Session) {
	player := NewPlayer(session)
	player.accel = GlobalConfig.Accel
	player.maxSpeed = GlobalConfig.MaxSpeed

	globalPlayGround.register <- player
}

func UnregisterGlobal(session *Session) {
	globalPlayGround.unregister <- session
}

func (ground *PlayGround) eventLoop() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(GlobalConfig.Tickcount))
	defer ticker.Stop()

	for {
		select {
		case player := <-ground.register:
			ground.registerPlayer(player)
		case session := <-ground.unregister:
			ground.unregisterPlayer(session)
		case <-ticker.C:
			ground.render()
		case <-ground.stop:
			fmt.Println("close play ground event loop")
			return
		case moveJog := <-ground.moveJog:
			ground.moveJogPlayer(moveJog)
		}
	}
}

func (ground *PlayGround) render() {
	moveNoti := &MovePlayerNotiMessage{}

	for id, player := range ground.players {
		player.speed = calcSpeed(player)
		if player.speed == 0 {
			continue
		}

		player.currDir = calcDir(player)

		fmt.Printf("current speed : %f, current dir : %v, current jog dir : %v\n", player.speed, player.currDir, player.jogDir)

		movePlayer(player)
		msg := &MovePlayerNotiMessage_MovePlayer{
			Id:  int64(id),
			X:   player.x,
			Y:   player.y,
			Dir: player.currDir,
		}

		moveNoti.Players = append(moveNoti.Players, msg)
	}

	if len(moveNoti.Players) != 0 {
		ground.broadcast(Protocol_MOVE_PLAYER_NOTI, moveNoti)
	}
}

func clamp(min, max, curr float32) float32 {
	if curr < min {
		return min
	}

	if curr > max {
		return max
	}

	return curr
}

func calcSpeed(player *Player) float32 {
	accel := player.accel
	if player.jogDir == Direction_NONE_DIR {
		accel = -accel
	}

	v := player.speed + accel
	return clamp(0, player.maxSpeed, v)
}

func calcDir(player *Player) Direction {
	if player.jogDir == Direction_NONE_DIR {
		return player.currDir
	}

	if player.currDir == Direction_NONE_DIR {
		return player.jogDir
	}

	diff := getDiffDirection(player.currDir, player.jogDir)
	if diff == 0 {
		return player.currDir
	} else if diff == 4 {
		if rand.Int()%2 == 0 {
			return directionToClockwise(player.currDir)
		} else {
			return directionToAntiClockwise(player.currDir)
		}
	} else if diff < 4 {
		return directionToClockwise(player.currDir)
	} else {
		return directionToAntiClockwise(player.currDir)
	}
}

func getDiffDirection(curr, jog Direction) int {
	count := 0
	for curr != jog {
		curr = directionToClockwise(curr)
		count++
	}

	return count
}

func directionToClockwise(dir Direction) Direction {
	to := int(dir) + 1
	if to == 9 {
		return Direction_NORTH
	}

	return Direction(to)
}

func directionToAntiClockwise(dir Direction) Direction {
	to := int(dir) - 1
	if to == 0 {
		return Direction_WEST_NORTH
	}

	return Direction(to)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}

func movePlayer(player *Player) {
	diagonalVal := player.speed * math.Sqrt2 / float32(2)

	if player.currDir == Direction_NORTH {
		player.y += player.speed
	} else if player.currDir == Direction_NORTH_EAST {
		player.x += diagonalVal
		player.y += diagonalVal
	} else if player.currDir == Direction_EAST {
		player.x += player.speed
	} else if player.currDir == Direction_EAST_SOUTH {
		player.x += diagonalVal
		player.y -= diagonalVal
	} else if player.currDir == Direction_SOUTH {
		player.y -= player.speed
	} else if player.currDir == Direction_SOUTH_WEST {
		player.x -= diagonalVal
		player.y -= diagonalVal
	} else if player.currDir == Direction_WEST {
		player.x -= player.speed
	} else if player.currDir == Direction_WEST_NORTH {
		player.x -= diagonalVal
		player.y += diagonalVal
	}
}

func (ground *PlayGround) moveJogPlayer(moveJog MoveJogChan) {
	player, ok := ground.players[moveJog.id]
	if !ok {
		fmt.Printf("cannot find player on move jog request id : %d\n", moveJog.id)
		return
	}

	player.jogDir = moveJog.dir
}

func (ground *PlayGround) registerPlayer(enterPlayer *Player) {
	enterMsg := &EnterPlayerNotiMessage{
		Player: enterPlayer.ToPlayerMessage(),
	}

	ground.broadcast(Protocol_ENTER_PLAYER_NOTI, enterMsg)
	ground.players[enterPlayer.GetId()] = enterPlayer
	welcomeMsg := &WelcomePlayerNotiMessage{}
	welcomeMsg.MyId = int64(enterPlayer.GetId())

	for _, player := range ground.players {
		welcomeMsg.Players = append(welcomeMsg.Players, player.ToPlayerMessage())
	}

	enterPlayer.SendMessage(Protocol_WELCOME_PLAYER_NOTI, welcomeMsg)
	fmt.Println("register to play ground")
}

func (ground *PlayGround) unregisterPlayer(session *Session) {
	var leavePlayer *Player

	for id, player := range ground.players {
		if id == session.GetId() {
			leavePlayer = player
			delete(ground.players, id)
			fmt.Println("unregister to play ground")
			break
		}
	}

	if leavePlayer == nil {
		return
	}

	leaveMsg := &LeavePlayerNotiMessage{
		Id: int64(leavePlayer.GetId()),
	}

	ground.broadcast(Protocol_LEAVE_PLAYER_NOTI, leaveMsg)
}

func (ground *PlayGround) broadcast(id Protocol, msg proto.Message) {
	for _, player := range ground.players {
		player.SendMessage(id, msg)
	}
}

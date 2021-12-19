package orc

import (
	"fmt"
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

func init() {
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

}

func (ground *PlayGround) moveJogPlayer(moveJog MoveJogChan) {
	player, ok := ground.players[moveJog.id]
	if !ok {
		fmt.Println("cannot find player on move jog request id : %d", moveJog.id)
		return
	}

	player.currDir = moveJog.dir
}

func (ground *PlayGround) registerPlayer(enterPlayer *Player) {
	welcomeMsg := WelcomePlayerNotiMessage{}

	for _, player := range ground.players {
		welcomeMsg.Players = append(welcomeMsg.Players, player.ToPlayerMessage())
	}

	enterPlayer.SendMessage(Protocol_WELCOME_PLAYER, welcomeMsg)

	enterMsg := EnterPlayerNotiMessage{
		Player: enterPlayer.ToPlayerMessage(),
	}

	ground.broadcast(Protocol_ENTER_PLAYER, enterMsg)

	ground.players[enterPlayer.GetId()] = enterPlayer

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

	leaveMsg := LeavePlayerNotiMessage{
		Id: int64(leavePlayer.GetId()),
	}

	ground.broadcast(Protocol_LEAVE_PLAYER, leaveMsg)
}

func (ground *PlayGround) broadcast(id Protocol, msg proto.Message) {
	for _, player := range ground.players {
		player.SendMessage(id, msg)
	}
}

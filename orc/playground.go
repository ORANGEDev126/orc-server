package orc

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
)

var globalPlayGround PlayGround

type MoveJogChan struct {
	id  uint64
	dir Direction
}

type ShootProjectileChan struct {
	playerId uint64
	angle    int
}

type AttackChan struct {
	playerId uint64
}

type PlayGround struct {
	players     map[uint64]*Player
	projectiles []*Projectile
	register    chan *Player
	unregister  chan *Session
	stop        chan bool
	moveJog     chan MoveJogChan
	shoot       chan ShootProjectileChan
	attackChan  chan AttackChan
}

func StartGlobal() {
	globalPlayGround = PlayGround{
		players:     map[uint64]*Player{},
		projectiles: []*Projectile{},
		register:    make(chan *Player),
		unregister:  make(chan *Session),
		stop:        make(chan bool),
		moveJog:     make(chan MoveJogChan),
		shoot:       make(chan ShootProjectileChan),
		attackChan:  make(chan AttackChan),
	}

	go globalPlayGround.eventLoop()
}

func TestSpawn() {
	go func() {
		time.Sleep(5 * time.Second)
		player := NewPlayer(NewSession(nil))
		globalPlayGround.register <- player
		time.Sleep(1 * time.Second)

		for {
			projectile := ShootProjectileChan{
				playerId: player.GetId(),
				angle:    rand.Intn(360),
			}
			globalPlayGround.shoot <- projectile
			time.Sleep(300 * time.Millisecond)
		}
	}()
}

func RegisterGlobal(session *Session) {
	player := NewPlayer(session)

	if len(globalPlayGround.players) == 0 {
		TestSpawn()
	}

	globalPlayGround.register <- player
}

func UnregisterGlobal(session *Session) {
	globalPlayGround.unregister <- session
}

func (ground *PlayGround) eventLoop() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(GlobalConfig.FrameTickCount))
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
		case projectile := <-ground.shoot:
			ground.shootProjectile(projectile.playerId, projectile.angle)
		case attack := <-ground.attackChan:
			ground.attack(attack.playerId)
		}
	}
}

func (ground *PlayGround) render() {
	moveNoti := &MoveObjectNotiMessage{}

	for i := 0; i < len(ground.projectiles); {
		nextPoint := ground.projectiles[i].NextPos()
		projectileR := ground.projectiles[i].circle.radius
		isCollision := false

		for _, player := range ground.players {
			if IsCollision(Circle{nextPoint, projectileR}, player.circle) {
				player.ProjectileAttacked(ground.projectiles[i].angle)
				ground.notifyProjectileAttack(ground.projectiles[i], player)
				isCollision = true
				break
			}
		}

		if isCollision {
			ground.projectiles = append(ground.projectiles[:i], ground.projectiles[i+1:]...)
		} else {
			ground.projectiles[i].Move(nextPoint)
			if nextPoint.x > 50 || nextPoint.y > 50 || nextPoint.x < -50 || nextPoint.y < -50 {
				ground.notifyRemoveObject(ground.projectiles[i].Id())
				ground.projectiles = append(ground.projectiles[:i], ground.projectiles[i+1:]...)
			} else {
				moveNoti.Objects = append(moveNoti.Objects,
					ground.projectiles[i].ToMoveObjectMessage())
				i++
			}
		}
	}

	for id, player := range ground.players {
		nextSpeed := player.UpdateNextSpeed()
		if nextSpeed == 0 {
			continue
		}

		nextDirection := player.UpdateNextDirection()
		nextPoint := player.NextPoint(nextSpeed, nextDirection)
		isCollision := false

		for otherId, otherPlayer := range ground.players {
			if id == otherId {
				continue
			}

			if !IsCollision(player.circle, otherPlayer.circle) &&
				IsCollision(Circle{nextPoint, player.circle.radius}, otherPlayer.circle) {
				fmt.Println("collision true players")
				isCollision = true
				break
			}
		}

		if isCollision {
			continue
		}

		for i := 0; i < len(ground.projectiles); {
			if IsCollision(Circle{nextPoint, player.circle.radius}, ground.projectiles[i].circle) {
				player.ProjectileAttacked(ground.projectiles[i].angle)
				ground.notifyProjectileAttack(ground.projectiles[i], player)
				ground.projectiles = append(ground.projectiles[:i], ground.projectiles[i+1:]...)
			} else {
				i++
			}
		}

		player.Move(nextSpeed, nextDirection, nextPoint)
		fmt.Printf("current speed : %f, current dir : %v, current jog dir : %v\n", player.speed, player.currDir, player.jogDir)
		moveNoti.Objects = append(moveNoti.Objects, player.ToMoveObjectMessage())
	}

	if len(moveNoti.Objects) != 0 {
		ground.broadcast(Notification_MOVE_OBJECT_NOTI, moveNoti)
	}
}

func (ground *PlayGround) notifyProjectileAttack(projectile *Projectile, player *Player) {
	msg := &ProjectileAttackNotiMessage{
		PlayerId:     int64(player.GetId()),
		ProjectileId: int64(projectile.id),
	}

	ground.broadcast(Notification_PROJECTILE_ATTACK_NOTI, msg)
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

	ground.broadcast(Notification_ENTER_PLAYER_NOTI, enterMsg)
	ground.players[enterPlayer.GetId()] = enterPlayer
	welcomeMsg := &WelcomePlayerNotiMessage{}
	welcomeMsg.MyId = int64(enterPlayer.GetId())

	for _, player := range ground.players {
		welcomeMsg.Players = append(welcomeMsg.Players, player.ToPlayerMessage())
	}

	enterPlayer.SendMessage(Notification_WELCOME_PLAYER_NOTI, welcomeMsg)
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
	ground.notifyRemoveObject(leavePlayer.GetId())
}

func (ground *PlayGround) notifyRemoveObject(id uint64) {
	leaveMsg := &LeaveObjectNotiMessage{
		Id: int64(id),
	}
	ground.broadcast(Notification_LEAVE_OBJECT_NOTI, leaveMsg)
}

func (ground *PlayGround) broadcast(id Notification, msg proto.Message) {
	for _, player := range ground.players {
		player.SendMessage(id, msg)
	}
}

func (ground *PlayGround) shootProjectile(id uint64, angle int) {
	if angle < 0 || angle > 360 {
		fmt.Println("wrong angle when shoot projectile", angle)
		return
	}

	player, ok := ground.players[id]
	if !ok {
		fmt.Println("cannot find player when shoot projectile id", id)
		return
	}

	point := GetPosAngle(player.circle.point, player.circle.radius+0.1, angle)
	projectile := NewProjectile(point, angle)

	ground.projectiles = append(ground.projectiles, projectile)
	ground.notifyProjectileEnter(projectile)
}

func (ground *PlayGround) attack(id uint64) {

}

func (ground *PlayGround) notifyProjectileEnter(projectile *Projectile) {
	msg := &EnterProjectileNotiMessage{
		Projectile: &ProjectileMessage{
			Id:    int64(projectile.id),
			X:     projectile.circle.point.x,
			Y:     projectile.circle.point.y,
			Angle: int32(projectile.angle),
		},
	}

	ground.broadcast(Notification_ENTER_PROJECTILE_NOTI, msg)
}

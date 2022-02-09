package orc

import "math/rand"

type Projectile struct {
	id     uint64
	circle Circle
	speed  float64
	angle  int
}

func NewProjectile(point Point, angle int) *Projectile {
	return &Projectile{
		id:     uint64(rand.Int63()),
		circle: Circle{point, GlobalConfig.ProjectileRadius},
		speed:  GlobalConfig.ProjectileSpeed,
		angle:  angle,
	}
}

func (projectile *Projectile) Id() uint64 {
	return projectile.id
}

func (projectile *Projectile) NextPos() Point {
	return GetPosAngle(projectile.circle.point,
		projectile.speed,
		projectile.angle)
}

func (projectile *Projectile) Move(point Point) {
	projectile.circle.point.x = point.x
	projectile.circle.point.y = point.y
}

func (projectile *Projectile) ToMoveObjectMessage() *MoveObjectNotiMessage_Object {
	return &MoveObjectNotiMessage_Object{
		Id: int64(projectile.id),
		X:  projectile.circle.point.x,
		Y:  projectile.circle.point.y,
	}
}

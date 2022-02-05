package orc

import "testing"

func TestNextDirection(t *testing.T) {
	mockPlayer := Player{
		currDir: Direction_EAST,
		jogDir:  Direction_EAST,
	}

	mockPlayer.jogDir = Direction_EAST_SOUTH
	if mockPlayer.UpdateNextDirection() != Direction_EAST_SOUTH {
		t.Error("case 1 fail")
	}

	mockPlayer.currDir = Direction_EAST
	mockPlayer.jogDir = Direction_NORTH_EAST
	if mockPlayer.UpdateNextDirection() != Direction_NORTH_EAST {
		t.Error("case 2 fail")
	}

	mockPlayer.currDir = Direction_EAST
	mockPlayer.jogDir = Direction_WEST_NORTH
	if mockPlayer.UpdateNextDirection() != Direction_NORTH_EAST {
		t.Error("case3 fail")
	}
}

func TestNextSpeed(t *testing.T) {
	GlobalConfig.FrameTickCount = 10
	GlobalConfig.PhysicsTickCount = 100
	GlobalConfig.MaxSpeed = 10
	mockPlayer := &Player{
		accel:    0.5,
		jogDir:   Direction_EAST,
		speed:    0,
		maxSpeed: 10,
	}

	if next := mockPlayer.UpdateNextSpeed(); next != 0.05 {
		t.Error("case1 fail", next)
	}

	mockPlayer.jogDir = Direction_NONE_DIR
	if next := mockPlayer.UpdateNextSpeed(); next != 0 {
		t.Error("case2 fail", next)
	}

	mockPlayer.speed = 0.1
	if next := mockPlayer.UpdateNextSpeed(); next != 0.05 {
		t.Error("case3 fail", next)
	}
}

func TestNextPoint(t *testing.T) {

}

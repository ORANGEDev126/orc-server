package orc

import "testing"

func TestNextDirection(t *testing.T) {
	mockPlayer := Player{
		currDir: Direction_EAST,
		jogDir:  Direction_EAST,
	}

	mockPlayer.jogDir = Direction_EAST_SOUTH
	if mockPlayer.NextDirection() != Direction_EAST_SOUTH {
		t.Error("case 1 fail")
	}

	mockPlayer.jogDir = Direction_NORTH_EAST
	if mockPlayer.NextDirection() != Direction_NORTH_EAST {
		t.Error("case 2 fail")
	}

	mockPlayer.jogDir = Direction_WEST_NORTH
	if mockPlayer.NextDirection() != Direction_NORTH_EAST {
		t.Error("case3 fail")
	}
}

func TestNextSpeed(t *testing.T) {

}

package orc

import "testing"

func TestDirectionToClockwise(t *testing.T) {
	if directionToClockwise(Direction_NORTH) != Direction_NORTH_EAST {
		t.Error("case 1 fail")
	}

	if directionToClockwise(Direction_WEST_NORTH) != Direction_NORTH {
		t.Error("case 2 fail")
	}

	if directionToAntiClockwise(Direction_NORTH) != Direction_WEST_NORTH {
		t.Error("case 3 fail")
	}

	if directionToAntiClockwise(Direction_SOUTH) != Direction_EAST_SOUTH {
		t.Error("case 4 fail")
	}
}

func TestCalcDir(t *testing.T) {
	mockPlayer := Player{
		currDir: Direction_EAST,
		jogDir:  Direction_EAST,
	}

	mockPlayer.jogDir = Direction_EAST_SOUTH
	if calcDir(&mockPlayer) != Direction_EAST_SOUTH {
		t.Error("case 1 fail")
	}

	mockPlayer.jogDir = Direction_NORTH_EAST
	if calcDir(&mockPlayer) != Direction_NORTH_EAST {
		t.Error("case 2 fail")
	}

	mockPlayer.jogDir = Direction_WEST_NORTH
	if calcDir(&mockPlayer) != Direction_NORTH_EAST {
		t.Error("case3 fail")
	}
}

package orc

import "testing"

func TestDirectionToClockwise(t *testing.T) {
	if DirectionToClockwise(Direction_NORTH) != Direction_NORTH_EAST {
		t.Error("case 1 fail")
	}

	if DirectionToClockwise(Direction_WEST_NORTH) != Direction_NORTH {
		t.Error("case 2 fail")
	}

	if DirectionToAntiClockwise(Direction_NORTH) != Direction_WEST_NORTH {
		t.Error("case 3 fail")
	}

	if DirectionToAntiClockwise(Direction_SOUTH) != Direction_EAST_SOUTH {
		t.Error("case 4 fail")
	}
}

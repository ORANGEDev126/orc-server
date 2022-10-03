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

func TestVectorToAngle(t *testing.T) {
	if VectorToAngle(Point{1, 1}) != 45 {
		t.Error("point (1, 1) fail")
	}

	if VectorToAngle(Point{0, 0}) != 0 {
		t.Error("point (0, 0) fail")
	}

	if VectorToAngle(Point{0, 10}) != 90 {
		t.Error("point (0, 10) fail")
	}

	result := VectorToAngle(Point{0, -10})

	if result != 270 {
		t.Error("point (0, -10) fail result : ", result)
	}
}

func TestIsAttackSuccess(t *testing.T) {

}

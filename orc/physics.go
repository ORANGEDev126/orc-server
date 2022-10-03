package orc

import "math"

type Circle struct {
	point  Point
	radius float64
}

type Point struct {
	x float64
	y float64
}

func (p *Point) Minus(other Point) Point {
	return Point{p.x - other.x, p.y - other.y}
}

func VectorToAngle(point Point) int {
	result := RadianToAngle(math.Atan(point.y / point.x))
	if result < 0 {
		return result + 360
	}

	return result
}

func AngleToRadian(angle int) float64 {
	return float64(angle) * 3.141592 / 180
}

func RadianToAngle(radian float64) int {
	return int(radian * 180 / 3.141592)
}

func GetPosAngle(point Point, dist float64, angle int) Point {
	radian := float64(angle) * 3.141592 / 180
	return Point{point.x + dist*math.Cos(radian),
		point.y + dist*math.Sin(radian)}
}

func IsCollision(circle, otherCircle Circle) bool {
	return Distance(circle.point, otherCircle.point) < circle.radius+otherCircle.radius
}

func IsAttackSuccess(attacker Point, opponent Circle, attackDistance float64, attackerDir Direction, attackRange int) bool {
	return Distance(attacker, opponent.point) < attackDistance+opponent.radius &&
		IsInRange(opponent.point, attackerDir, attackRange)
}

func GetPointFromDir(dir Direction) Point {
	switch dir {
	case Direction_NORTH:
		return Point{0, 1}
	case Direction_NORTH_EAST:
		return Point{math.Sqrt(2) / 2, math.Sqrt(2) / 2}
	case Direction_EAST:
		return Point{1, 0}
	case Direction_EAST_SOUTH:
		return Point{math.Sqrt(2) / 2, -math.Sqrt(2) / 2}
	case Direction_SOUTH:
		return Point{0, -1}
	case Direction_SOUTH_WEST:
		return Point{-math.Sqrt(2) / 2, -math.Sqrt(2) / 2}
	case Direction_WEST:
		return Point{-1, 0}
	case Direction_WEST_NORTH:
		return Point{-math.Sqrt(2) / 2, math.Sqrt(2) / 2}
	default:
		return Point{0, 0}
	}
}

func GetAngleFromDir(dir Direction) int {
	switch dir {
	case Direction_NORTH:
		return 90
	case Direction_NORTH_EAST:
		return 45
	case Direction_EAST:
		return 0
	case Direction_EAST_SOUTH:
		return 315
	case Direction_SOUTH:
		return 270
	case Direction_SOUTH_WEST:
		return 225
	case Direction_WEST:
		return 180
	case Direction_WEST_NORTH:
		return 135
	default:
		return 0
	}
}

func GetAngleFromPoint(point Point) int {
	t := math.Atan(point.y / point.x)
	return int(t * 180 / 3.141592)
}

func GetAngleRangeFromDir(dir Direction, inputRange int) (int, int) {
	r := GetAngleFromDir(dir)
	a := r - inputRange
	if a < 0 {
		a += 360
	}

	b := r + inputRange
	if b < a {
		b += 360
	}

	return a, b
}

func IsInRange(opponentPoint Point, attackerDir Direction, attackRange int) bool {
	a, b := GetAngleRangeFromDir(attackerDir, attackRange)
	c := GetAngleFromPoint(opponentPoint)

	if a < c && c < b {
		return true
	}

	c += 360

	if a < c && c < b {
		return true
	}

	return false
}

func Distance(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) +
		math.Pow(a.y-b.y, 2))
}

func Clamp(min, max, curr float64) float64 {
	if curr < min {
		return min
	}

	if curr > max {
		return max
	}

	return curr
}

func Abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}

func GetDiffDirection(curr, jog Direction) int {
	count := 0
	for curr != jog {
		curr = DirectionToClockwise(curr)
		count++
	}

	return count
}

func DirectionToClockwise(dir Direction) Direction {
	to := int(dir) + 1
	if to == 9 {
		return Direction_NORTH
	}

	return Direction(to)
}

func DirectionToAntiClockwise(dir Direction) Direction {
	to := int(dir) - 1
	if to == 0 {
		return Direction_WEST_NORTH
	}

	return Direction(to)
}

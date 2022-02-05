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

func GetPosAngle(point Point, dist float64, angle int) Point {
	radian := float64(angle) * 3.141592 / 180
	return Point{point.x + dist*math.Cos(radian),
		point.y + dist*math.Sin(radian)}
}

func IsCollision(circle, otherCircle Circle) bool {
	dis := math.Sqrt(math.Pow(circle.point.x-otherCircle.point.x, 2) +
		math.Pow(circle.point.y-otherCircle.point.y, 2))
	return dis < circle.radius+otherCircle.radius
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

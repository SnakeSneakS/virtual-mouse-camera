package model

import "math"

// HandLandmark point of finger from mediapipe
// check https://google.github.io/mediapipe/solutions/hands.html
// X: position x in image. 0 when hand is left from camera (right from you)
// Y: position y in image. 0 when hand is up from camera (up from you)
// Z: depth
type HandLandmark = [21]Position

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// add p2 to p1
func (p1 Position) Add(p2 Position) Position {
	return Position{
		X: p1.X + p2.X,
		Y: p1.Y + p2.Y,
		Z: p1.Z + p2.Z,
	}
}

func (p1 Position) Sub(p2 Position) Position {
	return Position{
		X: p1.X - p2.X,
		Y: p1.Y - p2.Y,
		Z: p1.Z - p2.Z,
	}
}

// multiply each element of p1 & p2
func (p Position) Multiply(p2 Position) Position {
	return Position{
		X: p.X * p2.X,
		Y: p.Y * p2.Y,
		Z: p.Z * p2.Z,
	}
}

func (p Position) Distance(p2 Position) float64 {
	return math.Sqrt((p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y))

}

func NewPosition(x, y, z float64) Position {
	return Position{
		X: x,
		Y: y,
		Z: z,
	}
}

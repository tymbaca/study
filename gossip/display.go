package main

import "math"

// Vector2 struct represents a 2D vector with x and y coordinates
type Vector2 struct {
	X, Y float32
}

// CircleLayout arranges `n` points evenly spaced in a circle
// and returns their positions as a slice of Vector2
func CircleLayout(n int, radius float64, offsetX, offsetY float32) []Vector2 {
	if n <= 0 {
		return nil // No points to place
	}

	// Slice to hold the calculated positions
	positions := make([]Vector2, n)

	for i := 0; i < n; i++ {
		// Calculate angle for each point
		angle := 2 * math.Pi * float64(i) / float64(n)
		// Calculate x and y coordinates
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		positions[i] = Vector2{X: float32(x) + offsetX, Y: float32(y) + offsetY}
	}

	return positions
}

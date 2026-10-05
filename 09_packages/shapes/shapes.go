// Package shapes provides basic 2D shape calculations.
// This is a real, importable package — not package main.
package shapes

import "math"

// Circle is exported (capitalized) so other packages can use it.
type Circle struct {
	Radius float64
}

// Area is an exported method — capitalized = public.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// helperNote is unexported (lowercase) — only visible inside this package.
const helperNote = "internal use only"

// Describe uses the unexported constant, demonstrating package-private state.
func Describe(c Circle) string {
	return "circle: " + helperNote
}

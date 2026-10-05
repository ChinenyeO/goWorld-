package main

import (
	"fmt"
	"math"
)

// An interface defines a set of methods. Any type that implements those
// methods satisfies the interface automatically — no "implements" keyword.
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Any function that takes a Shape works with Circle, Rectangle, or any
// future type that implements Area()/Perimeter() — this is polymorphism.
func describe(s Shape) {
	fmt.Printf("  area=%.2f perimeter=%.2f\n", s.Area(), s.Perimeter())
}

// The empty interface "any" (alias for interface{}) can hold a value of
// any type — used sparingly, when you genuinely don't know the type ahead.
func printAnything(v any) {
	fmt.Println("value:", v)
}

func main() {
	shapes := []Shape{
		Circle{Radius: 2},
		Rectangle{Width: 3, Height: 4},
	}

	for _, s := range shapes {
		describe(s)

		// Type switch: ask "which concrete type is this?" at runtime.
		switch v := s.(type) {
		case Circle:
			fmt.Println("  it's a circle with radius", v.Radius)
		case Rectangle:
			fmt.Println("  it's a rectangle", v.Width, "x", v.Height)
		}
	}

	printAnything(42)
	printAnything("hello")
	printAnything(Circle{Radius: 1})

	// Type assertion: extract the concrete type from an interface value.
	var s Shape = Circle{Radius: 5}
	if c, ok := s.(Circle); ok {
		fmt.Println("asserted circle radius:", c.Radius)
	}

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Define a `Triangle` struct (Base, Height) above and implement
//      Area()/Perimeter() so it satisfies Shape (Perimeter can just
//      return 0 if you don't want to compute it properly yet).
//   2. Add it to the `shapes` slice in main and confirm describe()
//      works on it without any changes to describe() itself.
//   3. Add a case for Triangle in the type switch.
//   4. Write a small interface of your own, e.g.
//      `type Greeter interface { Greet() string }`, implement it on
//      a type, and call it through the interface.
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

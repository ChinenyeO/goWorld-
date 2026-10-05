package main

import "fmt"

type Rectangle struct {
	Width, Height float64
}

// A method with a VALUE receiver — gets a copy of the Rectangle.
// Use value receivers for read-only operations.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// A method with a POINTER receiver — can mutate the original struct.
// Use pointer receivers when the method needs to modify the receiver,
// or when the struct is large (avoids copying).
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	rect := Rectangle{Width: 3, Height: 4}

	fmt.Println("Area:", rect.Area())
	fmt.Println("Perimeter:", rect.Perimeter())

	// Go automatically takes the address for you here: it's really
	// (&rect).Scale(2), because Scale needs a pointer receiver.
	rect.Scale(2)
	fmt.Println("After scaling:", rect)
	fmt.Println("New area:", rect.Area())

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Add a value-receiver method `IsSquare() bool` on Rectangle above.
//   2. Add a pointer-receiver method `Grow(amount float64)` that adds
//      `amount` to both Width and Height.
//   3. Define a new type (e.g. `type Circle struct { Radius float64 }`)
//      with an Area() method, and call it from here.
//   4. Try changing Scale to a value receiver — call it, print rect
//      afterward, and see that the change no longer sticks. Why?
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

package main

import (
	"fmt"

	// Local packages are imported by their full module path (from go.mod)
	// plus the folder path — not a relative path like "./shapes".
	"example/goWorld/09_packages/shapes"
)

func main() {
	c := shapes.Circle{Radius: 3}

	// Only exported (capitalized) names are visible from outside the package.
	fmt.Println("area:", c.Area())
	fmt.Println(shapes.Describe(c))

	// shapes.helperNote would NOT compile here — it's unexported.

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here (edit shapes/shapes.go too):
//   1. In shapes/shapes.go, add a new exported type `Square struct { Side float64 }`
//      with an Area() method, then use it here as shapes.Square{...}.
//   2. Add an exported function `shapes.Perimeter(c Circle) float64` in
//      the shapes package and call it from here.
//   3. Try referencing shapes.helperNote from this file — read the
//      compiler error to see unexported visibility enforced.
//   4. Create a second subpackage (e.g. 09_packages/greetings) with its
//      own exported Hello() function, and import it here alongside shapes.
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

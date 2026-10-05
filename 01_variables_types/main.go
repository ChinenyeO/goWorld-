package main

import "fmt"
import "strconv"

func main() {
	// var declares a variable with an explicit type.
	var name string = "Gopher"

	// Go can infer the type from the value, so the type is optional.
	var age = 5

	// ':=' is shorthand for "var x = value" — only usable inside functions.
	city := "Lagos"

	// Go has fixed-width numeric types. int/float64 are the usual defaults.
	var score int = 42
	var price float64 = 19.99
	var isAwesome bool = true

	// Multiple variables in one line.
	var x, y int = 1, 2

	// Uninitialized variables get a "zero value" instead of being undefined.
	var zeroInt int
	var zeroStr string
	var zeroBool bool

	// Constants are fixed at compile time and can't be reassigned.
	const Pi = 3.14159

	fmt.Println(name, age, city, score, price, isAwesome, x, y)
	fmt.Println("zero values ->", zeroInt, zeroStr, zeroBool)
	fmt.Println("Pi:", Pi)

	// Type conversion is explicit in Go — no implicit int->float64 etc.
	var i int = 10
	var f float64 = float64(i)
	fmt.Println("converted:", f)

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Declare a variable for your own name and age using ':='.
//   2. Declare a rune (a single character, e.g. 'A') and print it —
//      notice fmt.Println shows a number, not a letter. Why?
//   3. Convert an int to a string using strconv.Itoa (you'll need to
//      add "strconv" to the import block above).
//   4. Declare a constant for your favorite number and use it in a
//      calculation.
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	
	name := "Chinenye"
	age := 27

	var aRune string  = "A"
	fmt.Println(aRune, name)

	strconv.Itoa(age)
	fmt.Println(age)

	const cons = 1738.123

	fmt.Println(cons + 99.22)




}

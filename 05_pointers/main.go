package main

import "fmt"

// Takes a copy — changes here don't affect the caller's variable.
func incByValue(n int) {
	n++
}

// Takes a pointer — dereferencing with *n lets us modify the original.
func incByPointer(n *int) {
	*n++
}

type Counter struct {
	Count int
}

func main() {
	x := 10

	incByValue(x)
	fmt.Println("after incByValue:", x) // unchanged

	// & takes the address of a variable, producing a pointer.
	incByPointer(&x)
	fmt.Println("after incByPointer:", x) // changed

	// A pointer variable holds a memory address.
	var p *int = &x
	fmt.Println("p points to:", *p) // * dereferences: "the value at this address"

	// Modifying through the pointer changes the original.
	*p = 100
	fmt.Println("x is now:", x)

	// new(T) allocates a zero-valued T and returns a pointer to it.
	c := new(Counter)
	c.Count = 1 // Go auto-dereferences for field access — no need for (*c).Count
	fmt.Println("counter:", *c)

	// nil is the zero value for pointers — "points to nothing".
	var np *int
	fmt.Println("nil pointer:", np == nil)

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Write a function `double(n *int)` that doubles the value it
//      points to, then call it on a local variable.
//   2. Write a function `swap(a, b *int)` that swaps two ints via
//      pointers, and prove it works by printing before/after.
//   3. Create a pointer to a struct (e.g. Counter) and modify a field
//      through the pointer using dot notation (no need for (*p).Field).
//   4. Deliberately dereference a nil pointer and read the panic
//      message it produces (wrap it in recover() if you want the
//      program to keep running).
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

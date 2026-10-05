package main

import "fmt"

// A struct groups related fields into one type.
type Person struct {
	Name string
	Age  int
}

// Structs can be composed (embedding) instead of using inheritance.
type Employee struct {
	Person   // embedded struct — Employee "has" all Person's fields promoted
	Company string
	Salary  float64
}

func main() {
	// Struct literal with named fields (preferred — order-independent).
	p1 := Person{Name: "Ada", Age: 30}

	// Positional literal — must match field order exactly.
	p2 := Person{"Grace", 28}

	// Zero-value struct: fields get their type's zero value.
	var p3 Person
	p3.Name = "Linus"
	p3.Age = 55

	fmt.Println(p1, p2, p3)

	e := Employee{
		Person:  Person{Name: "Margaret", Age: 40},
		Company: "NASA",
		Salary:  95000,
	}
	// Promoted fields: access Name/Age directly, no need for e.Person.Name.
	fmt.Println(e.Name, e.Age, e.Company, e.Salary)

	// Struct equality: comparable if all fields are comparable.
	fmt.Println("p1 == p2:", p1 == p2)

	// Anonymous struct — a one-off type with no name, handy for small scopes.
	point := struct{ X, Y int }{X: 1, Y: 2}
	fmt.Println(point)

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Define your own struct above main, e.g. `type Book struct { ... }`
//      with fields like Title, Author, Pages.
//   2. Create a couple of Book values and print them.
//   3. Try embedding: make a `type Ebook struct { Book; FileSizeMB int }`
//      and access a promoted field directly.
//   4. Compare two structs with `==` — what happens if one field is a
//      slice? (Hint: try it and read the compiler error.)
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

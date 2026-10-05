package main

import "fmt"

// Basic function: name, params (with types), return type.
func add(a int, b int) int {
	return a + b
}

// Consecutive params of the same type can share one type annotation.
func multiply(a, b int) int {
	return a * b
}

// Multiple return values are common in Go — often (result, error).
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide %d by zero", a)
	}
	return a / b, nil
}

// Named return values: "quotient" and "remainder" are declared up front
// and returned automatically by a bare "return".
func divmod(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return
}

// Variadic functions accept a variable number of arguments.
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func largestNum(nums ...int) int {

	largestNum := nums[0]

	for _, num := range nums[1:]{
		if num > largestNum{
			largestNum = num
		}
	}

	return largestNum

}

// Functions are values: they can be passed as arguments (higher-order funcs).
func applyTwice(f func(int) int, x int) int {
	return f(f(x))
}

func isEven(n int) bool {
	return n % 2 == 0

}

func main() {
	fmt.Println("add:", add(2, 3))
	fmt.Println("multiply:", multiply(4, 5))

	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("divide:", result)
	}

	// This one triggers the error path.
	_, err = divide(10, 0)
	fmt.Println("divide by zero error:", err)

	q, r := divmod(17, 5)
	fmt.Println("divmod:", q, r)

	fmt.Println("sum:", sum(1, 2, 3, 4, 5))

	// Anonymous function passed inline.
	double := func(n int) int { return n * 2 }
	fmt.Println("applyTwice:", applyTwice(double, 3))

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Write a function `isEven(n int) bool` above main and call it.
//   2. Write a variadic function `max(nums ...int) int` that returns
//      the largest number.
//   3. Write a function that returns two values: the square and cube
//      of an int, and destructure both with `:=`.
//   4. Write a higher-order function like applyTwice, but that takes
//      a func(int, int) int (e.g. apply it to add or multiply).
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	fmt.Println(isEven(1000000008))
	fmt.Println(largestNum(100 , 1, 2,3, 1000))
}

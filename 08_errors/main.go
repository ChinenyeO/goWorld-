package main

import (
	"errors"
	"fmt"
)

// The 'error' type is just an interface with one method: Error() string.
// Returning an error as the last value is Go's standard way to signal failure.

// A sentinel error — a package-level value callers can compare against.
var ErrNotFound = errors.New("item not found")

// A custom error type carries structured data, not just a string.
type ValidationError struct {
	Field string
	Msg   string
}

// Implementing Error() string makes ValidationError satisfy the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Msg)
}

func findItem(id int) (string, error) {
	items := map[int]string{1: "apple", 2: "banana"}
	item, ok := items[id]
	if !ok {
		return "", ErrNotFound
	}
	return item, nil
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Msg: "must not be negative"}
	}
	return nil
}

// errors.Wrap-style: fmt.Errorf with %w wraps an error, preserving the chain
// so errors.Is / errors.As can still find the original underneath.
func loadUser(id int) error {
	_, err := findItem(id)
	if err != nil {
		return fmt.Errorf("loadUser %d: %w", id, err)
	}
	return nil
}

func main() {
	item, err := findItem(1)
	fmt.Println(item, err)

	_, err = findItem(99)
	if err != nil {
		fmt.Println("error:", err)
		// errors.Is checks the error chain for a specific sentinel value.
		fmt.Println("is ErrNotFound:", errors.Is(err, ErrNotFound))
	}

	err = validateAge(-5)
	if err != nil {
		fmt.Println("error:", err)
		// errors.As checks whether an error in the chain matches a type,
		// and if so, extracts it into the target variable.
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Println("bad field:", ve.Field)
		}
	}

	err = loadUser(99)
	fmt.Println("wrapped error:", err)
	fmt.Println("still matches ErrNotFound:", errors.Is(err, ErrNotFound))

	// panic/recover: for truly exceptional, unrecoverable-in-place situations
	// — NOT a substitute for normal error returns.
	safeDivide(10, 0)

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Write a function `parsePositive(n int) (int, error)` that
//      returns an error if n <= 0, and call it with both a valid and
//      an invalid value.
//   2. Define your own sentinel error, e.g. `var ErrTooLarge = errors.New(...)`,
//      and return it from a function; check it with errors.Is.
//   3. Define your own custom error struct (like ValidationError) for
//      a different scenario, implement Error() string on it, and use
//      errors.As to extract it.
//   4. Wrap an error from #1 with fmt.Errorf and %w, then confirm
//      errors.Is still finds the original underneath.
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

func safeDivide(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic:", r)
		}
	}()
	fmt.Println(a / b) // dividing by zero panics for ints
}

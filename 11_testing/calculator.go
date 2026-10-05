package main
import "errors"
// EXERCISE: Divide currently returns one value (int), and dividing by
// zero panics. Change it to return TWO values: the result AND an error.
//
// Syntax hint (unrelated example): a function that returns a value and
// an error lists both return types in parentheses:
//
//	func FindUser(id int) (string, error) {
//	    // ...
//	}
//
// Hint: once a function returns two values, every return statement in it
// must return two values as well.
func Divide(a int, b int) (int,error) {
	if b == 0{
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

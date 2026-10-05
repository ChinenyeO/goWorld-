package main

import "fmt"

func main() {
	// --- Arrays vs Slices ---
	// Arrays have a fixed size baked into their type. Rarely used directly.
	var arr [3]int = [3]int{1, 2, 3}
	fmt.Println("array:", arr)

	// Slices are the flexible, resizable "view" over an array that you'll
	// actually use day to day.
	nums := []int{1, 2, 3}
	fmt.Println("slice:", nums, "len:", len(nums), "cap:", cap(nums))

	// append grows the slice, reassigning is required (may return a new
	// underlying array if capacity was exceeded).
	nums = append(nums, 4, 5)
	fmt.Println("after append:", nums)

	// Slicing syntax: s[low:high] — half-open range, excludes 'high'.
	fmt.Println("slice[1:3]:", nums[1:3])

	// make(T, len, cap) pre-allocates a slice.
	buffer := make([]int, 0, 10)
	fmt.Println("empty buffer:", buffer, "cap:", cap(buffer))

	// range iterates index, value.
	for i, v := range nums {
		fmt.Printf("  nums[%d] = %d\n", i, v)
	}

	// --- Maps ---
	// map[KeyType]ValueType — Go's hash map / dictionary.
	ages := map[string]int{
		"Ada":   30,
		"Grace": 28,
	}
	ages["Linus"] = 55 // add/update a key

	// The "comma ok" idiom: check whether a key exists.
	age, ok := ages["Ada"]
	fmt.Println("Ada's age:", age, "found:", ok)

	missingAge, ok := ages["Nobody"]
	fmt.Println("missing age:", missingAge, "found:", ok) // zero value if absent

	delete(ages, "Grace")

	for name, a := range ages {
		fmt.Printf("  %s is %d\n", name, a)
	}
	fmt.Println("map len:", len(ages))

	yourTurn()
}

// ---------------------------------------------------------------------
// YOUR TURN
// Try this here:
//   1. Build a []string of your favorite foods and append two more.
//   2. Slice it to get just the first two items.
//   3. Build a map[string]bool of things you've done today and check
//      membership with the comma-ok idiom.
//   4. Write a loop that sums only the even numbers in a []int slice.
//   5. Make a slice of structs (e.g. []Person from topic 3) and range
//      over it, printing each one.
// ---------------------------------------------------------------------
func yourTurn() {
	fmt.Println("\n--- your turn ---")
	// write your code here
}

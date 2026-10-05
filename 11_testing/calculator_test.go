package main

import "testing"

func TestDivide(t *testing.T) {
	// 1. Define test cases

	tests := []struct{
		a int
		b int
		expected int
		name string
		wantErr bool
	}{
		{8,2,4, "Positive Even", false},
		{10,2,5, "Positive Even Big", false},
		{9,3,3, "Positive Odd",false},
		{0,5,0, "Dividing with zero",false},
		{-10,2,-5, "Divide with negative",false},
		{10, 0, 0, "Divide by zero", true},

	}


	// 2. Loop through test cases

	for _,test := range tests {
		t.Run(test.name, func(t *testing.T) {
		// HINT: after you change Divide to return two values, this line
		// won't compile anymore. When a function returns two values, Go
		// code commonly captures both like this:
		//
		//	value, err := someFunction()
		result,err := Divide(test.a, test.b)
		if test.wantErr && err == nil {
    	// fail
			t.Errorf("wantErr is true AND err is nil")

		}

		if !test.wantErr && err != nil {
    	// fail
			t.Errorf("wantErr is false AND err is not nil")

		}
		if result != test.expected{
			t.Errorf("Expected %d got %d", test.expected, result)
		}
        
    	})

		

	}

	// 3. Compare actual result with expected result
}

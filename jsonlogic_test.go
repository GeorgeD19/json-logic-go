package jsonlogic

// import (
// 	"testing"
// )

// func TestSoftEqual(t *testing.T) {
// 	test1a := 1
// 	test1b := 1
// 	result := false

// 	result = Operations["=="].run(test1a, test1b).(bool)
// 	if result != true {
// 		t.Fatalf("%d is equal to %d!", test1a, test1b)
// 	}
// }

// func TestLess(t *testing.T) {
// 	test1a := 1
// 	test1b := 2
// 	result := false

// 	result = Operations["<"].run(test1a, test1b).(bool)
// 	if result != true {
// 		t.Fatalf("%d is less than %d!", test1a, test1b)
// 	}

// 	test2a := 2
// 	test2b := 1

// 	result = Operations["<"].run(test2a, test2b).(bool)
// 	if result != false {
// 		t.Fatalf("%d is not less than %d!", test2a, test2b)
// 	}
// }

// func TestPercentage(t *testing.T) {
// 	test1a := 20
// 	test1b := 50

// 	result := Operations["%"].run(test1a, test1b).(float64)
// 	if result != 40.0 {
// 		t.Fatalf("%d should be 40 percent of %d! We got %d instead.", test1a, test1b, result)
// 	}
// }

// func TestRemoveOperation(t *testing.T) {
// 	test1 := "=="

// 	test1a := 1
// 	test1b := 1
// 	result := false

// 	// Throw true
// 	result = Operations["=="].run(test1a, test1b).(bool)

// 	RemoveOperation(test1)

// 	// Throw err
// 	result = Operations["=="].run(test1a, test1b).(bool)

// }

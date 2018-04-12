package jsonlogic

import "testing"

func TestLess(t *testing.T) {
	test1a := 1
	test1b := 2
	result := false

	result = Operations["<"].run(test1a, test1b).(bool)
	if result != true {
		t.Fatalf("%d is less than %d!", test1a, test1b)
	}

	test2a := 2
	test2b := 1

	result = Operations["<"].run(test2a, test2b).(bool)
	if result != false {
		t.Fatalf("%d is not less than %d!", test2a, test2b)
	}
}

func TestPercentage(t *testing.T) {
	test1a := 20
	test1b := 50

	result := Operations["%"].run(test1a, test1b).(float64)
	if result != 40.0 {
		t.Fatalf("%d should be 40 percent of %d!", test1a, test1b)
	}
}

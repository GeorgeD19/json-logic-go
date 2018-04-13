package jsonlogic

import (
	"testing"
)

func TestSoftEqualTrue(t *testing.T) {
	logic := `{"==" : [ 10, "10" ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestHardEqualTrue(t *testing.T) {
	logic := `{"===" : [ "10", "10" ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestHardEqualFalse(t *testing.T) {
	logic := `{"===" : [ 10, "10" ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != false {
		t.Fatal("Logic should return false")
	}
}

func TestNotSoftEqualTrue(t *testing.T) {
	logic := `{"!=" : [ "100", 110 ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestNotSoftEqualFalse(t *testing.T) {
	logic := `{"!=" : [ "110", 110 ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != false {
		t.Fatal("Logic should return false")
	}
}

func TestNotHardEqualTrue(t *testing.T) {
	logic := `{"!==" : [ "110", 110 ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestNotHardEqualFalse(t *testing.T) {
	logic := `{"!==" : [ 110, 110 ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != false {
		t.Fatal("Logic should return false")
	}
}

func TestLessTrue(t *testing.T) {
	logic := `{"<" : [ 100, 110 ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestLessFalse(t *testing.T) {
	logic := `{"<" : [ 110, 110 ]}`
	data := `{}`

	// Should throw true
	result, _ := Apply(logic, data)

	if result == true {
		t.Fatal("Logic should return false")
	}
}

func TestLessTrueVar(t *testing.T) {
	logic := `{"<" : [ { "var" : "temp" }, 110 ]}`
	data := `{ "temp" : 100 }`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestLessEqualVarTrue(t *testing.T) {
	logic := `{"<=" : [ { "var" : "temp" }, 110 ]}`
	data := `{ "temp" : 110 }`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestAndTrue(t *testing.T) {
	logic := `{ "and" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{ "temp" : 100, "pie" : { "filling" : "apple" } }`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

func TestOrTrue(t *testing.T) {
	logic := `{ "or" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{ "temp" : 120, "pie" : { "filling" : "apple" } }`

	// Should throw true
	result, _ := Apply(logic, data)

	if result != true {
		t.Fatal("Logic should return true")
	}
}

// Tests
// SoftEqualTest TODO Move this to test
// func SoftEqualTest() {
// Throws false
// fmt.Println(strconv.FormatBool(SoftEqual(1, 2)))
// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual(2, 2)))

// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual(0, false)))

// // Throws false
// fmt.Println(strconv.FormatBool(SoftEqual(true, false)))
// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual(false, false)))

// // Throws false
// fmt.Println(strconv.FormatBool(SoftEqual("1", "2")))
// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual("1", "1")))

// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual("1", 1)))
// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual("1", true)))

// // Throws false
// fmt.Println(strconv.FormatBool(SoftEqual("1", nil)))
// // Throws true
// fmt.Println(strconv.FormatBool(SoftEqual(nil, nil)))
// }

// HardEqualTest checks the functionality of the HardEqual function
// func HardEqualTest() {
// Throws false
// fmt.Println(strconv.FormatBool(HardEqual(1, 2)))
// // Throws true
// fmt.Println(strconv.FormatBool(HardEqual(2, 2)))

// // Throws false
// fmt.Println(strconv.FormatBool(HardEqual(true, false)))
// // Throws true
// fmt.Println(strconv.FormatBool(HardEqual(false, false)))

// // Throws false
// fmt.Println(strconv.FormatBool(HardEqual("1", "2")))
// // Throws true
// fmt.Println(strconv.FormatBool(HardEqual("1", "1")))

// // Throws false
// fmt.Println(strconv.FormatBool(HardEqual("1", 1)))
// // Throws false
// fmt.Println(strconv.FormatBool(HardEqual("1", true)))
// }

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

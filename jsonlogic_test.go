package jsonlogic

import (
	"testing"

	"github.com/spf13/cast"
)

func TestMaxTrue(t *testing.T) {
	logic := `{"max": [
		4,
		3,
		5
		]} `
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToFloat64(result) != 5.0 {
		t.Fatal("Logic should return 5.0")
	}
}

func TestMinTrue(t *testing.T) {
	logic := `{"min": [
		4,
		3,
		5
		]} `
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToFloat64(result) != 3.0 {
		t.Fatal("Logic should return 3.0")
	}
}

func TestMultiplyTrue(t *testing.T) {
	logic := `{"*": [
		2,
		2
		]} `
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToFloat64(result) != 4.0 {
		t.Fatal("Logic should return 4.0")
	}
}

func TestDivideTrue(t *testing.T) {
	logic := `{"/": [
		2,
		2
		]} `
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToFloat64(result) != 1.0 {
		t.Fatal("Logic should return 1.0")
	}
}

func TestPercentageSoftEquals(t *testing.T) {
	logic := `{"==": [{"%" : [20,50]}, 40]} `
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestPercentageTrue(t *testing.T) {
	logic := `{"%" : [20,50]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToInt(result) != 40 {
		t.Fatal("Logic should return 40")
	}
}

func TestTruthyTrue(t *testing.T) {
	logic := `{"!!" : []}`
	data := `{"a":1,"b":2}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestNotTruthyTrue(t *testing.T) {
	logic := `{"!" : []}`
	data := `{"a":1,"b":2}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != false {
		t.Fatal("Logic should return false")
	}
}

func TestVarTrue(t *testing.T) {
	logic := `{"var" : "a"}`
	data := `{"a":1,"b":2}`

	// Should return 1
	result, _ := Apply(logic, data)

	if cast.ToInt(result) != 1 {
		t.Fatal("Logic should return 1")
	}
}

func TestVarFalse(t *testing.T) {
	logic := `{"var" : "c"}`
	data := `{"a":1,"b":2}`

	// Should return 1
	result, _ := Apply(logic, data)

	if result != nil {
		t.Fatal("Logic should return nil")
	}
}

func TestVarDefault(t *testing.T) {
	logic := `{"var":["z", 26]}`
	data := `{"a":1,"b":2}`

	// Should return 26
	result, _ := Apply(logic, data)

	if cast.ToInt(result) != 26 {
		t.Fatal("Logic should return 26")
	}
}

func TestVarNest(t *testing.T) {
	logic := `{"var" : "champ.name"}`
	data := `{
		"champ" : {
		  "name" : "Fezzig",
		  "height" : 223
		},
		"challenger" : {
		  "name" : "Dread Pirate Roberts",
		  "height" : 183
		}
	  }`

	// Should return Fezzig
	result, _ := Apply(logic, data)

	if cast.ToString(result) != "Fezzig" {
		t.Fatal("Logic should return Fezzig")
	}
}

func TestVarArray(t *testing.T) {
	logic := `{"var":1}`
	data := `["zero", "one", "two"]`

	// Should return Fezzig
	result, _ := Apply(logic, data)

	if cast.ToString(result) != "one" {
		t.Fatal("Logic should return one")
	}

}

func TestSoftEqualTrue(t *testing.T) {
	logic := `{"==" : [ 10, "10" ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestHardEqualTrue(t *testing.T) {
	logic := `{"===" : [ "10", "10" ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestHardEqualFalse(t *testing.T) {
	logic := `{"===" : [ 10, "10" ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != false {
		t.Fatal("Logic should return false")
	}
}

func TestNotSoftEqualTrue(t *testing.T) {
	logic := `{"!=" : [ "100", 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestNotSoftEqualFalse(t *testing.T) {
	logic := `{"!=" : [ "110", 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != false {
		t.Fatal("Logic should return false")
	}
}

func TestNotHardEqualTrue(t *testing.T) {
	logic := `{"!==" : [ "110", 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestNotHardEqualFalse(t *testing.T) {
	logic := `{"!==" : [ 110, 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != false {
		t.Fatal("Logic should return false")
	}
}

func TestLessTrue(t *testing.T) {
	logic := `{"<" : [ 100, 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestLessFalse(t *testing.T) {
	logic := `{"<" : [ 110, 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) == true {
		t.Fatal("Logic should return false")
	}
}

func TestLessTrueVar(t *testing.T) {
	logic := `{"<" : [ { "var" : "temp" }, 110 ]}`
	data := `{ "temp" : 100 }`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestLessEqualVarTrue(t *testing.T) {
	logic := `{"<=" : [ { "var" : "temp" }, 110 ]}`
	data := `{ "temp" : 110 }`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

func TestAndTrue(t *testing.T) {
	logic := `{ "and" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{ "temp" : 100, "pie" : { "filling" : "apple" } }`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

// func TestCat() {

// }

func TestMinusTrue(t *testing.T) {
	logic := `{"-":[
		1,
		1
	  ]}`
	data := `{}`

	// Should return True string
	result, _ := Apply(logic, data)

	if cast.ToFloat64(result) != 0.0 {
		t.Fatal("Logic should return 0.0")
	}
}

func TestPlusTrue(t *testing.T) {
	logic := `{"+":[
		1,
		1
	  ]}`
	data := `{}`

	// Should return True string
	result, _ := Apply(logic, data)

	if cast.ToFloat64(result) != 2.0 {
		t.Fatal("Logic should return 2.0")
	}
}

func TestIfTrue(t *testing.T) {
	logic := `{"if":[
		{"==":["b", "b"]},
		"True",
		"False"
	  ]}`
	data := `{"a":"apple", "b":"banana"}`

	// Should return True string
	result, _ := Apply(logic, data)

	if cast.ToString(result) != "True" {
		t.Fatal("Logic should return True")
	}
}

func TestOrTrue(t *testing.T) {
	logic := `{ "or" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{ "temp" : 120, "pie" : { "filling" : "apple" } }`

	// Should return true
	result, _ := Apply(logic, data)

	if cast.ToBool(result) != true {
		t.Fatal("Logic should return true")
	}
}

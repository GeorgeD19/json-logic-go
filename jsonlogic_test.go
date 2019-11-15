package jsonlogic

import (
	"reflect"
	"testing"

	"github.com/spf13/cast"
)

func TestVar(t *testing.T) {
	rule := `{ "var" : ["a"] }`
	data := `{ "a":1, "b":2 }`

	result, _ := Apply(rule, data)

	if cast.ToInt(result) != 1 {
		t.Fatalf("rule should return 1, instead returned %s", result)
	}
}

func TestVarSugar(t *testing.T) {
	rule := `{"var":"a"}`
	data := `{"a":1,"b":2}`

	result, _ := Apply(rule, data)

	if cast.ToInt(result) != 1 {
		t.Fatalf("rule should return 1, instead returned %s", result)
	}
}

func TestVarFallback(t *testing.T) {
	rule := `{"var":["z", 26]}`
	data := `{"a":1,"b":2}`

	result, _ := Apply(rule, data)

	if cast.ToInt(result) != 26 {
		t.Fatalf("rule should return 26, instead returned %s", result)
	}
}

func TestVarDotNotation(t *testing.T) {
	rule := `{"var" : "champ.name"}`
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

	result, _ := Apply(rule, data)

	if cast.ToString(result) != "Fezzig" {
		t.Fatalf("rule should return Fezzig, instead returned %s", result)
	}
}

func TestVarNumericIndex(t *testing.T) {
	rule := `{"var":1}`
	data := `["zero", "one", "two"]`

	result, _ := Apply(rule, data)

	if cast.ToString(result) != "one" {
		t.Fatalf("rule should return one, instead returned %s", result)
	}
}

func TestVarComplex(t *testing.T) {
	rule := `{ "and" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{
		"temp" : 100,
		"pie" : { "filling" : "apple" }
	  }`

	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

func TestVarEmpty(t *testing.T) {
	rule := `{ "cat" : [
		"Hello, ",
		{"var":""}
	] }`
	data := `"Dolly"`

	result, _ := Apply(rule, data)

	if cast.ToString(result) != "Hello, Dolly" {
		t.Fatalf("rule should return Hello, Dolly, instead returned %s", result)
	}
}

// Between

// Between exclusive
func TestBetweenExclusiveLess(t *testing.T) {
	rule := `{"<" : [1, 2, 3]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

func TestBetweenExclusiveLessString(t *testing.T) {
	rule := `{"<" : [1, 2, "3"]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

func TestBetweenExclusiveLessNot(t *testing.T) {
	rule := `{"<" : [1, 1, 3]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != false {
		t.Fatalf("rule should return false, instead returned %s", result)
	}
}

func TestBetweenExclusiveLessMore(t *testing.T) {
	rule := `{"<" : [1, 4, 3]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != false {
		t.Fatalf("rule should return false, instead returned %s", result)
	}
}

func TestBetweenInclusiveLess(t *testing.T) {
	rule := `{"<=" : [1, 2, 3]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

func TestBetweenInclusiveLessNot(t *testing.T) {
	rule := `{"<=" : [1, 1, 3]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

func TestBetweenInclusiveLessMore(t *testing.T) {
	rule := `{"<=" : [1, 4, 3]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != false {
		t.Fatalf("rule should return false, instead returned %s", result)
	}
}

func TestBetweenExclusiveLessData(t *testing.T) {
	rule := `{ "<": [0, {"var":"temp"}, 100]}`
	data := `{"temp" : 37}`

	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

// Arithmetic

// Addition

func TestAdd(t *testing.T) {
	rule := `{"+":[4,2]}`

	result, _ := Run(rule)

	if cast.ToInt(result) != 6 {
		t.Fatalf("rule should return 6, instead returned %s", result)
	}
}

func TestMinus(t *testing.T) {
	rule := `{"-":[4,2]}`

	result, _ := Run(rule)

	if cast.ToInt(result) != 2 {
		t.Fatalf("rule should return 2, instead returned %s", result)
	}
}

func TestMultiply(t *testing.T) {
	rule := `{"*":[4,2]}`

	result, _ := Run(rule)

	if cast.ToInt(result) != 8 {
		t.Fatalf("rule should return 8, instead returned %s", result)
	}
}

func TestDivide(t *testing.T) {
	rule := `{"/":[4,2]}`

	result, _ := Run(rule)

	if cast.ToInt(result) != 2 {
		t.Fatalf("rule should return 2, instead returned %s", result)
	}
}

func TestAddArgs(t *testing.T) {
	rule := `{"+":[2,2,2,2,2]}`

	result, _ := Run(rule)

	if cast.ToInt(result) != 10 {
		t.Fatalf("rule should return 10, instead returned %s", result)
	}
}

func TestMultiplyArgs(t *testing.T) {
	rule := `{"*":[2,2,2,2,2]}`

	result, _ := Run(rule)

	if cast.ToInt(result) != 32 {
		t.Fatalf("rule should return 32, instead returned %s", result)
	}
}

func TestMinusArgPos(t *testing.T) {
	rule := `{"-": 2 }`

	result, _ := Run(rule)

	if cast.ToInt(result) != -2 {
		t.Fatalf("rule should return -2, instead returned %s", result)
	}
}

func TestMinusArgNeg(t *testing.T) {
	rule := `{"-": -2 }`

	result, _ := Run(rule)

	if cast.ToInt(result) != 2 {
		t.Fatalf("rule should return 2, instead returned %s", result)
	}
}

func TestAddCast(t *testing.T) {
	rule := `{"+": "3.14" }`

	result, _ := Run(rule)

	if cast.ToFloat64(result) != 3.14 {
		t.Fatalf("rule should return 3.14, instead returned %s", result)
	}
}

// String Operations

// In

func TestIn(t *testing.T) {
	rule := `{"in":["Spring", "Springfield"]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

func TestInArray(t *testing.T) {
	rule := `{"in":[ "Ringo", ["John", "Paul", "George", "Ringo"] ]}`

	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatalf("rule should return true, instead returned %s", result)
	}
}

// Cat

// %
// func TestModolo(t *testing.T) {
// 	rule := `{"%": [101,2]}`

// 	result, _ := Run(rule)

// 	if cast.ToInt(result) != 1 {
// 		t.Fatalf("rule should return 1, instead returned %s", result)
// 	}
// }

// substr

func TestSubstrPosition(t *testing.T) {
	rule := `{"substr": ["jsonlogic", 4]}`

	result, _ := Run(rule)

	if cast.ToString(result) != "logic" {
		t.Fatalf("rule should return logic, instead returned %s", result)
	}
}

func TestSubstrPositionNeg(t *testing.T) {
	rule := `{"substr": ["jsonlogic", -5]}`

	result, _ := Run(rule)

	if cast.ToString(result) != "logic" {
		t.Fatalf("rule should return logic, instead returned %s", result)
	}
}

func TestSubstrPositionLength(t *testing.T) {
	rule := `{"substr": ["jsonlogic", 1, 3]}`

	result, _ := Run(rule)

	if cast.ToString(result) != "son" {
		t.Fatalf("rule should return son, instead returned %s", result)
	}
}

func TestSubstrPositionLengthNeg(t *testing.T) {
	rule := `{"substr": ["jsonlogic", 4, -2]}`

	result, _ := Run(rule)

	if cast.ToString(result) != "log" {
		t.Fatalf("rule should return log, instead returned %s", result)
	}
}

// Merge

func TestMerge(t *testing.T) {
	rule := `{"merge":[ [1,2], [3,4] ]}`

	result, _ := Run(rule)
	target := []int{1, 2, 3, 4}

	if reflect.DeepEqual(result, target) {
		t.Fatalf("rule should return [1,2,3,4], instead returned %s", result)
	}
}

func TestMergeMixed(t *testing.T) {
	rule := `{"merge":[ 1, 2, [3,4] ]}`

	result, _ := Run(rule)
	target := []int{1, 2, 3, 4}

	if reflect.DeepEqual(result, target) {
		t.Fatalf("rule should return [1,2,3,4], instead returned %s", result)
	}
}

func TestMergeStringMixed(t *testing.T) {
	rule := `{"missing" :
		{ "merge" : [
		  "vin",
		  {"if": [{"var":"financing"}, ["apr", "term"], [] ]}
		]}
	  }`
	data := `{"financing":true}`

	result, _ := Apply(rule, data)
	target := []string{"vin", "apr", "term"}

	if reflect.DeepEqual(result, target) {
		t.Fatalf("rule should return [vin, api, term], instead returned %s", result)
	}
}

func TestMergeStringMixedMissing(t *testing.T) {
	rule := `{"missing" :
		{ "merge" : [
		  "vin",
		  {"if": [{"var":"financing"}, ["apr", "term"], [] ]}
		]}
	  }`
	data := `{"financing":false}`

	result, _ := Apply(rule, data)
	target := []string{"vin"}

	if reflect.DeepEqual(result, target) {
		t.Fatalf("rule should return [vin], instead returned %s", result)
	}
}

// Missing

func TestMissing(t *testing.T) {
	rule := `{"missing":["a", "b"]}`
	data := `{"a":"apple", "c":"carrot"}`

	result, _ := Apply(rule, data)
	target := []string{"b"}

	if reflect.DeepEqual(result, target) {
		t.Fatalf("rule should return [b], instead returned %s", result)
	}
}

// TODO Clean up tests to match against http://jsonlogic.com/operations.html

func TestMaxTrue(t *testing.T) {
	rule := `{"max": [
		4,
		3,
		5
		]} `

	// Should return true
	result, _ := Run(rule)

	if cast.ToFloat64(result) != 5.0 {
		t.Fatal("rule should return 5.0")
	}
}

func TestMinTrue(t *testing.T) {
	rule := `{"min": [
		4,
		3,
		5
		]} `

	// Should return true
	result, _ := Run(rule)

	if cast.ToFloat64(result) != 3.0 {
		t.Fatal("rule should return 3.0")
	}
}

func TestMultiplyTrue(t *testing.T) {
	rule := `{"*": [
		2,
		2
		]} `

	// Should return true
	result, _ := Run(rule)

	if cast.ToFloat64(result) != 4.0 {
		t.Fatal("rule should return 4.0")
	}
}

func TestDivideTrue(t *testing.T) {
	rule := `{"/": [
		2,
		2
		]} `

	// Should return true
	result, _ := Run(rule)

	if cast.ToFloat64(result) != 1.0 {
		t.Fatal("rule should return 1.0")
	}
}

func TestPercentageSoftEquals(t *testing.T) {
	rule := `{"==": [{"%" : [20,50]}, 40]} `

	// Should return true
	result, _ := Run(rule)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestPercentageTrue(t *testing.T) {
	rule := `{"%" : [20,50]}`

	// Should return true
	result, _ := Run(rule)

	if cast.ToInt(result) != 40 {
		t.Fatal("rule should return 40")
	}
}

func TestTruthyTrue(t *testing.T) {
	rule := `{"!!" : []}`
	data := `{"a":1,"b":2}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestNotTruthyTrue(t *testing.T) {
	rule := `{"!" : []}`
	data := `{"a":1,"b":2}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != false {
		t.Fatal("rule should return false")
	}
}

func TestVarTrue(t *testing.T) {
	rule := `{"var" : "a"}`
	data := `{"a":1,"b":2}`

	// Should return 1
	result, _ := Apply(rule, data)

	if cast.ToInt(result) != 1 {
		t.Fatal("rule should return 1")
	}
}

func TestVarFalse(t *testing.T) {
	rule := `{"var" : "c"}`
	data := `{"a":1,"b":2}`

	// Should return 1
	result, _ := Apply(rule, data)

	if result != nil {
		t.Fatal("rule should return nil")
	}
}

func TestVarDefault(t *testing.T) {
	rule := `{"var":["z", 26]}`
	data := `{"a":1,"b":2}`

	// Should return 26
	result, _ := Apply(rule, data)

	if cast.ToInt(result) != 26 {
		t.Fatal("rule should return 26")
	}
}

func TestVarNest(t *testing.T) {
	rule := `{"var" : "champ.name"}`
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
	result, _ := Apply(rule, data)

	if cast.ToString(result) != "Fezzig" {
		t.Fatal("rule should return Fezzig")
	}
}

func TestVarArray(t *testing.T) {
	rule := `{"var":1}`
	data := `["zero", "one", "two"]`

	// Should return Fezzig
	result, _ := Apply(rule, data)

	if cast.ToString(result) != "one" {
		t.Fatal("rule should return one")
	}

}

func TestSoftEqualTrue(t *testing.T) {
	rule := `{"==" : [ 10, "10" ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestHardEqualTrue(t *testing.T) {
	rule := `{"===" : [ "10", "10" ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestHardEqualFalse(t *testing.T) {
	rule := `{"===" : [ 10, "10" ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != false {
		t.Fatal("rule should return false")
	}
}

func TestNotSoftEqualTrue(t *testing.T) {
	rule := `{"!=" : [ "100", 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestNotSoftEqualFalse(t *testing.T) {
	rule := `{"!=" : [ "110", 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != false {
		t.Fatal("rule should return false")
	}
}

func TestNotHardEqualTrue(t *testing.T) {
	rule := `{"!==" : [ "110", 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestNotHardEqualFalse(t *testing.T) {
	rule := `{"!==" : [ 110, 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != false {
		t.Fatal("rule should return false")
	}
}

func TestLessTrue(t *testing.T) {
	rule := `{"<" : [ 100, 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestLessFalse(t *testing.T) {
	rule := `{"<" : [ 110, 110 ]}`
	data := `{}`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) == true {
		t.Fatal("rule should return false")
	}
}

func TestLessTrueVar(t *testing.T) {
	rule := `{"<" : [ { "var" : "temp" }, 110 ]}`
	data := `{ "temp" : 100 }`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestLessEqualVarTrue(t *testing.T) {
	rule := `{"<=" : [ { "var" : "temp" }, 110 ]}`
	data := `{ "temp" : 110 }`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestAndTrue(t *testing.T) {
	rule := `{ "and" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{ "temp" : 100, "pie" : { "filling" : "apple" } }`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

func TestMinusTrue(t *testing.T) {
	rule := `{"-":[
		1,
		1
	  ]}`
	data := `{}`

	// Should return True string
	result, _ := Apply(rule, data)

	if cast.ToFloat64(result) != 0.0 {
		t.Fatal("rule should return 0.0")
	}
}

func TestPlusTrue(t *testing.T) {
	rule := `{"+":[
		1,
		1
	  ]}`
	data := `{}`

	// Should return True string
	result, _ := Apply(rule, data)

	if cast.ToFloat64(result) != 2.0 {
		t.Fatal("rule should return 2.0")
	}
}

func TestIfTrue(t *testing.T) {
	rule := `{
		"if": [
			{
				"==": [{
					"var": "b"
				}, {
					"var": "b"
				}]
			},
			"True",
			"False"
		]
	}`
	data := `{"a":"apple", "b":"banana"}`

	// Should return True string
	result, _ := Apply(rule, data)

	if cast.ToString(result) != "True" {
		t.Fatalf("rule should return True, instead returned %s", result)
	}
}

func TestOrTrue(t *testing.T) {
	rule := `{ "or" : [
		{"<" : [ { "var" : "temp" }, 110 ]},
		{"==" : [ { "var" : "pie.filling" }, "apple" ] }
	  ] }`
	data := `{ "temp" : 120, "pie" : { "filling" : "apple" } }`

	// Should return true
	result, _ := Apply(rule, data)

	if cast.ToBool(result) != true {
		t.Fatal("rule should return true")
	}
}

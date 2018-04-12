package jsonlogic

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dariubs/percent"

	"math"

	"github.com/buger/jsonparser"
)

// var Operations = map[string]interface{}{
// 	"==":  SoftEqual,
// 	"===": HardEqual,
// 	"!=":  NotSoftEqual,
// 	"!==": NotHardEqual,
// 	">":   More,
// 	">=":  MoreEqual,
// 	"<":   Less,
// 	"<=":  LessEqual,
// 	"!":   NotTruthy,
// 	"!!":  Truthy,
// 	// "%": lambda a, b: a % b,
// 	// "and": lambda *args: reduce(lambda total, arg: total and arg, args, True),
// 	// "or": lambda *args: reduce(lambda total, arg: total or arg, args, False),
// 	// "?:": lambda a, b, c: b if a else c,
// 	// "if": if_,
// 	// "log": lambda a: logger.info(a) or a,
// 	// "in": lambda a, b: a in b if "__contains__" in dir(b) else False,
// 	// "cat": lambda *args: "".join(str(arg) for arg in args),
// 	// "+": plus,
// 	// "*": lambda *args: reduce(lambda total, arg: total * float(arg), args, 1),
// 	// "-": minus,
// 	// "/": lambda a, b=None: a if b is None else float(a) / float(b),
// 	// "min": lambda *args: min(args),
// 	// "max": lambda *args: max(args),
// 	// "merge": merge,
// 	// "count": lambda *args: sum(1 if a else 0 for a in args),
// }

// Errors
var (
	InvalidOperationError = errors.New("Invalid Operation %s")
)

var Operations = make(map[string]Operation)

type Operation interface {
	run(a ...interface{}) interface{}
}

func AddOperation(Name string, Callable Operation) {
	Operations[Name] = Callable
}

func init() {
	AddOperation("==", new(SoftEqual))
	AddOperation("===", new(HardEqual))
	AddOperation("!=", new(NotSoftEqual))
	AddOperation("!==", new(NotHardEqual))
	AddOperation(">", new(More))
	AddOperation(">=", new(MoreEqual))
	AddOperation("<", new(Less))
	AddOperation("<=", new(LessEqual))
	AddOperation("!", new(NotTruthy))
	AddOperation("!!", new(Truthy))
	AddOperation("%", new(Percentage))
	AddOperation("and", new(And))
}

// func GetOperator() {

// }

// func GetValues() {

// }

// func IsLogic() {

// }

// func UsesData(logic string) {

// }

// And implements the 'and' conditional requiring all bubbled up bools to be true
type And struct{}

func (o And) run(a ...interface{}) interface{} {
	result := true
	// So we would run through the array of objects and execute it and if any values are false we stop and return false
	results, _, err := ParseArray(a[0].(string), a[1].(string))
	if err != nil {
		return false
	}

	for _, res := range results {
		if res == false {
			result = false
		}
	}

	return result
}

type Or struct{}

func (o Or) run(a ...interface{}) interface{} {
	// So we would run through the array of objects and if any operations return true then we return true
	return false
}

// Less implements the '<' operator with JS-style type coertion. Returns bool.
type Less struct{}

func (o Less) run(a ...interface{}) interface{} {
	fmt.Println("Less")
	fmt.Println(GetType(a))
	if ToFloat(a[0]) < ToFloat(a[1]) {
		return true
	}
	return false
}

// LessEqual implements the '<=' operator with JS-style type coertion. Returns bool.
type LessEqual struct{}

func (o LessEqual) run(a ...interface{}) interface{} {
	if ToFloat(a[0]) <= ToFloat(a[1]) {
		return true
	}
	return false
}

// More implements the '>' operator with JS-style type coertion. Returns bool.
type More struct{}

func (o More) run(a ...interface{}) interface{} {
	return Operations["<"].run(a[1], a[0])
}

// MoreEqual implements the '>=' operator with JS-style type coertion. Returns bool.
type MoreEqual struct{}

func (o MoreEqual) run(a ...interface{}) interface{} {
	return Operations["<"].run(a[1], a[0]).(bool) || Operations["=="].run(a[0], a[1]).(bool)
}

// SoftEqual implements the '==' operator, which does type JS-style coertion. Returns bool.
type SoftEqual struct{}

func (o SoftEqual) run(a ...interface{}) interface{} {
	fmt.Println("SoftEqual")
	if ToString(a[0]) == ToString(a[1]) {
		return true
	}
	return false
}

// NotSoftEqual implements the '!=' operator, which does type JS-style coertion. Returns bool.
type NotSoftEqual struct{}

func (o NotSoftEqual) run(a ...interface{}) interface{} {
	return !Operations["!="].run(a[0], a[1]).(bool)
}

// HardEqual Implements the '===' operator, which does type JS-style coertion. Returns bool.
type HardEqual struct{}

func (o HardEqual) run(a ...interface{}) interface{} {

	aType := GetType(a[0])
	bType := GetType(a[1])

	if aType != bType {
		return false
	}

	switch v := aType; v {
	case 1:
		return a[0].(int) == a[1].(int)
	case 2:
		return a[0].(float64) == a[1].(float64)
	case 3:
		return a[0].(string) == a[1].(string)
	case 4:
		return a[0].(bool) == a[1].(bool)
	default:
		return false
	}
}

// NotHardEqual implements the '!==' operator, which does type JS-style coertion. Returns bool.
type NotHardEqual struct{}

func (o NotHardEqual) run(a ...interface{}) interface{} {
	results := !Operations["!=="].run(a[0], a[1]).(bool)
	return results
}

// Percentage implements the '%' operator, which does type JS-style coertion. Returns float64.
type Percentage struct{}

func (o Percentage) run(a ...interface{}) interface{} {
	return percent.PercentOf(ToInt(a[0]), ToInt(a[1]))
}

// Truthy implements the '!!' operator with JS-style type coertion. Returns bool.
type Truthy struct{}

func (o Truthy) run(a ...interface{}) interface{} {
	aVal := ToString(a[0])
	if aVal == "0" {
		return true
	}
	return ToBool(a[0])
}

// NotTruthy implements the '!' operator with JS-style type coertion. Returns bool.
type NotTruthy struct{}

func (o NotTruthy) run(a ...interface{}) interface{} {
	return !Operations["!!"].run(a[0], a[1]).(bool)
}

// Var implements the 'var' operator, which does type JS-style coertion.
type Var struct{}

func (o Var) run(a ...interface{}) interface{} {
	fmt.Println("Var")
	fmt.Println(a)
	// We can use this to get the data value
	// jsonparser.Get(data, a[0])

	return true
}

// func ParseOperation(logic string) (operation []byte, value []byte], err error) {
// An operation could return anything
func ParseObject(logic string, data string) (res interface{}, err error) {
	result := false
	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if operation, ok := Operations[string(key)]; ok {
			result = operation.run(string(value), data).(bool)
		} else {
			return InvalidOperationError
		}
		return nil
	})

	if err != nil {
		return false, err
	}
	return false, nil
}

// func GetOperation(Type string) (o Operation, err error) {
// 	if _, ok := Operations[Type]; ok {
// 		return Operations[Type], nil
// 	} else {
// 		return nil, InvalidOperationError
// 	}
// }

// ParseArray runs all the operations in an array and returns an array of bools
func ParseArray(logic string, data string) (res []bool, off int, e error) {
	result := []bool{}
	offset, err := jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// fmt.Printf("Value: '%s'\n Type: %s\n Offset: %s\n", string(value), dataType, string(offset))
		switch dataType {
		case jsonparser.Object:
			ParseObject(string(value), data)
			break
		case jsonparser.Array:
			ParseArray(string(value), data)
			break
		}
	})
	if err != nil {
		return result, 0, err
	}
	return result, offset, nil
}

// Apply is the entry function to parse logic and optional data
func Apply(logic string, data string) (res bool, errs error) {

	// Ensure data is object
	if data == `` {
		data = `{}`
	}

	// Must be an object to kick off process
	result, err := ParseObject(logic, data)
	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

// GetType returns an int to map against type
func GetType(a interface{}) int {
	switch a.(type) {
	case int:
		return 1
	case float64:
		return 2
	case string:
		return 3
	case bool:
		return 4
	default:
		return 0
	}
}

// ToString converts recognised types to string
func ToString(a interface{}) string {
	aType := GetType(a)

	switch v := aType; v {
	case 1:
		return strconv.Itoa(a.(int))
	case 2:
		return strconv.FormatFloat(a.(float64), 'E', -1, 64)
	case 3:
		return a.(string)
	case 4:
		if a.(bool) {
			return ToString(1)
		}
		return ToString(0)
	default:
		return ""
	}
}

// ToInt converts recognised types to int
func ToInt(a ...interface{}) int {
	aType := GetType(a)

	switch v := aType; v {
	case 1:
		return a[0].(int)
	case 2:
		return int(a[0].(float64))
	case 3:
		fmt.Println(a)
		// TODO Check if string int or float rather than just float
		result, err := ParseObject(a[0].(string), a[1].(string))
		fmt.Println(result)
		fmt.Println(err)

		f, err := strconv.ParseFloat(a[0].(string), 64)
		if err != nil {
			return 0
		}
		return int(f)
	case 4:
		if a[0].(bool) {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// ToFloat converts recognised types to float64
// TODO If we pass in an object var then we need to get the data extracted first before conversion
func ToFloat(a ...interface{}) float64 {
	aType := GetType(a)

	fmt.Println(a)

	switch v := aType; v {
	case 1:
		return float64(a[0].(int))
	case 2:
		return a[0].(float64)
	case 3:
		// We should check if it is a var object

		result, err := ParseObject(a[0].(string), a[1].(string))
		fmt.Println(result)
		fmt.Println(err)

		f, err := strconv.ParseFloat(a[0].(string), 64)
		if err != nil {
			return 0.0
		}
		return f
	case 4:
		if a[0].(bool) {
			return 1.0
		}
		return 0.0
	default:
		return 0.0
	}
}

// ToBool converts recognised types to bool
func ToBool(a interface{}) bool {
	aType := GetType(a)

	switch v := aType; v {
	case 1:
		if a.(int) >= 1 {
			return true
		}
		return false
	case 2:
		if math.Round(a.(float64)) >= 1.0 {
			return true
		}
		return false
	case 3:
		f, err := strconv.ParseFloat(a.(string), 64)
		if err != nil {
			return false
		}
		if math.Round(f) >= 1.0 {
			return true
		}
		return false
	case 4:
		return a.(bool)
	default:
		return false
	}
}

func TestOperationCall(Type string, a ...interface{}) {
	fmt.Println(Operations[Type].run(a[0]))
}

func ToBoolTest() {
	fmt.Println(ToBool(0.4))
	fmt.Println(ToBool(0.5))
}

// Tests
// SoftEqualTest TODO Move this to test
func SoftEqualTest() {
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
}

// HardEqualTest checks the functionality of the HardEqual function
func HardEqualTest() {
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
}

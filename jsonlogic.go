package jsonlogic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dariubs/percent"
	"github.com/spf13/cast"

	"github.com/buger/jsonparser"
)

// Errors
var (
	ErrInvalidOperation = errors.New("Invalid Operation %s")
)

// Operations contains all possible operations that can be performed
var Operations = make(map[string]Operation)

// Operation interface that allows operations to be registered in a list
type Operation interface {
	run(a ...interface{}) interface{}
}

// AddOperation adds possible operation to Operations library
func AddOperation(name string, callable Operation) {
	Operations[name] = callable
}

// RemoveOperation removes possible operation from Operations library if it exists
func RemoveOperation(name string) {
	_, ok := Operations[name]
	if ok {
		delete(Operations, name)
	}
}

// RunOperation is to ensure that any operation ran doesn't crash the system if it doesn't exist
func RunOperation(name string, params ...interface{}) (res interface{}, err error) {
	_, ok := Operations[name]
	if ok {
		return Operations[name].run(params), nil
	}
	return nil, ErrInvalidOperation
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
	AddOperation("or", new(Or))
	AddOperation("var", new(Var))

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

}

// SoftEqual type is entry point for parser
type SoftEqual struct{}

func (o SoftEqual) run(a ...interface{}) interface{} {
	values := GetValues(cast.ToString(a[0]), cast.ToString(a[1]))
	return SoftEqualOperation(cast.ToString(values[0]), cast.ToString(values[1]))
}

// SoftEqualOperation implements the '==' operator, which does type JS-style coertion. Returns bool.
func SoftEqualOperation(a string, b string) bool {
	if a == b {
		return true
	}
	return false
}

// HardEqual type is entry point for parser
type HardEqual struct{}

func (o HardEqual) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return HardEqualOperation(values[0], values[1])
}

// HardEqualOperation Implements the '===' operator, which does type JS-style coertion. Returns bool.
func HardEqualOperation(a interface{}, b interface{}) bool {
	if GetType(a) != GetType(b) {
		return false
	}

	if a == b {
		return true
	}

	return false
}

// NotSoftEqual type is entry point for parser
type NotSoftEqual struct{}

func (o NotSoftEqual) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return NotSoftEqualOperation(cast.ToString(values[0]), cast.ToString(values[1]))
}

// NotSoftEqualOperation implements the '!=' operator, which does type JS-style coertion. Returns bool.
func NotSoftEqualOperation(a string, b string) bool {
	return !SoftEqualOperation(a, b)
}

// NotHardEqual type is entry point for parser
type NotHardEqual struct{}

func (o NotHardEqual) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return NotHardEqualOperation(values[0], values[1])
}

// NotHardEqualOperation implements the '!==' operator, which does type JS-style coertion. Returns bool.
func NotHardEqualOperation(a interface{}, b interface{}) bool {
	return !HardEqualOperation(a, b)
}

// More type is entry point for parser
type More struct{}

func (o More) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return MoreOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// MoreOperation implements the '>' operator with JS-style type coertion. Returns bool.
func MoreOperation(a float64, b float64) bool {
	return LessEqualOperation(b, a)
}

// MoreEqual type is entry point for parser
type MoreEqual struct{}

func (o MoreEqual) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return MoreEqualOperation(cast.ToString(values[0]), cast.ToString(values[1]))
}

// MoreEqualOperation implements the '>=' operator with JS-style type coertion. Returns bool.
func MoreEqualOperation(a string, b string) bool {
	return LessOperation(cast.ToFloat64(b), cast.ToFloat64(a)) || SoftEqualOperation(a, b)
}

// Less type is entry point for parser
type Less struct{}

func (o Less) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return LessOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// LessOperation implements the '<' operator with JS-style type coertion. Returns bool.
func LessOperation(a float64, b float64) bool {
	if a < b {
		return true
	}
	return false
}

// LessEqual implements the '<=' operator with JS-style type coertion. Returns bool.
type LessEqual struct{}

func (o LessEqual) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return LessEqualOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

func LessEqualOperation(a float64, b float64) bool {
	if a <= b {
		return true
	}
	return false
}

// And type is entry point for parser
type And struct{}

func (o And) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return AndOperation(values)
}

// AndOperation implements the 'and' conditional requiring all bubbled up bools to be true
func AndOperation(values []interface{}) bool {
	result := true
	for _, res := range values {
		if res == false {
			result = false
		}
	}
	return result
}

// Or type is entry point for parser
type Or struct{}

func (o Or) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return OrOperation(values)
}

// OrOperation implements the 'or' conditional requiring at least one of the bubbled up bools to be true
func OrOperation(values []interface{}) bool {
	result := false
	for _, res := range values {
		if res == true {
			result = true
		}
	}
	return result
}

// GetValues will return values of any kind
func GetValues(logic string, data string) (results []interface{}) {
	_, err := jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			res, _ := ParseObject(string(value), data)
			results = append(results, res)
			break
		// case jsonparser.Array:
		// 	fmt.Println("GetValueArray")
		// 	results = append(results, value)
		// 	break
		case jsonparser.String:
			results = append(results, cast.ToString(value))
			break
		case jsonparser.Number:
			fmt.Println("GetValueNumber")
			results = append(results, cast.ToFloat64(cast.ToString(value)))
			break
		case jsonparser.Boolean:
			fmt.Println("GetValueBoolean")
			results = append(results, cast.ToBool(value))
			break
		case jsonparser.Null:
			fmt.Println("GetValueNull")
			results = append(results, value)
			break
		}
	})
	if err != nil {
		return nil
	}
	return results
}

// GetType returns an int to map against type so we can see if we are dealing with a specific type of data or an object operation
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

		// It could be an object or array
		fmt.Println("Don't know what this is")
		return 0
	}
}

// ParseArray is basically ParseValue runs all the operations in an array and returns an array of values whether that be bool or whatever
func ParseArray(logic string, data string) (res []bool, off int, e error) {
	result := []bool{}
	offset, err := jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// fmt.Printf("Value: '%s'\n Type: %s\n Offset: %s\n", string(value), dataType, string(offset))

		switch dataType {
		case jsonparser.Object:
			objectResult, _ := ParseObject(string(value), data)
			result[offset] = objectResult.(bool)
			break
		case jsonparser.Array:
			// Would we really have an array within an array? probably not
			// result, err := ParseArray(string(value), data)
			break
		}
	})
	if err != nil {
		return nil, 0, err
	}
	return result, offset, nil
}

// Percentage implements the '%' operator, which does type JS-style coertion. Returns float64.
type Percentage struct{}

func (o Percentage) run(a ...interface{}) interface{} {
	// a[0] is the value passed in
	// a[1] is any data passed in so it can trickle down to any var objects

	return percent.PercentOf(cast.ToInt(a[0]), cast.ToInt(a[1]))
}

func PercentageOperation(a int, b int) float64 {
	return percent.PercentOf(a, b)
}

// Truthy implements the '!!' operator with JS-style type coertion. Returns bool.
type Truthy struct{}

func (o Truthy) run(a ...interface{}) interface{} {
	aVal := cast.ToString(a[0])
	if aVal == "0" {
		return true
	}
	return cast.ToBool(a[0])
}

// NotTruthy implements the '!' operator with JS-style type coertion. Returns bool.
type NotTruthy struct{}

func (o NotTruthy) run(a ...interface{}) interface{} {
	return !Operations["!!"].run(a[0], a[1]).(bool)
}

// Var implements the 'var' operator, which does type JS-style coertion.
type Var struct{}

func (o Var) run(a ...interface{}) interface{} {
	fmt.Println("VarOperation")
	fmt.Println(a[0].(string))
	return OperationVar(a[0].(string), a[1].(string))
}

func OperationVar(logic string, data string) interface{} {
	key := strings.Split(logic, ".")
	// value, dataType, offset, err := jsonparser.Get([]byte(a[1].(string)), variable...)
	value, dataType, _, _ := jsonparser.Get([]byte(data), key...)
	translatedValue := TranslateType(value, dataType)
	return translatedValue
}

func TranslateType(data []byte, dataType jsonparser.ValueType) interface{} {
	switch dataType {
	case jsonparser.String:
		return string(data)
		break
	case jsonparser.Number:
		numberString := cast.ToString(data)
		numberFloat := cast.ToFloat64(numberString)
		return numberFloat
		break
	case jsonparser.Boolean:
		return string(data)
		break
	case jsonparser.Null:
		return string(data)
		break
	default:
		return nil
	}
	return nil

}

// An operation could return anything
func ParseObject(logic string, data string) (res interface{}, err error) {
	// result := interface{}
	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		// fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)

		// res, err = RunOperation(string(key), string(value), data)
		// if err != nil {
		// 	return ErrInvalidOperation
		// }
		// return nil

		if operation, ok := Operations[string(key)]; ok {
			res = operation.run(string(value), data)
		} else {
			return ErrInvalidOperation
		}
		return nil
	})

	if err != nil {
		return false, err
	}
	return res, nil
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

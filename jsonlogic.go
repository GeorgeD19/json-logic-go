package jsonlogic

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
// func RunOperation(name string, params ...interface{}) (res interface{}, err error) {
// 	_, ok := Operations[name]
// 	if ok {
// 		return Operations[name].run(params), nil
// 	}
// 	return nil, ErrInvalidOperation
// }

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
	AddOperation("?", new(If))
	AddOperation("if", new(If))
	AddOperation("log", new(Log))
	AddOperation("+", new(Plus))
	AddOperation("-", new(Minus))
	AddOperation("*", new(Multiply))
	AddOperation("/", new(Divide))
	AddOperation("min", new(Min))
	AddOperation("max", new(Max))

	// 	// "in": lambda a, b: a in b if "__contains__" in dir(b) else False,
	// 	// "cat": lambda *args: "".join(str(arg) for arg in args),
	// 	// "merge": merge,
	// 	// "count": lambda *args: sum(1 if a else 0 for a in args),

}

func RunOperation(key string, logic string, data string) (res interface{}) {
	values := GetValues(logic, data)
	switch key {
	case "==":
		res = SoftEqualOperation(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "===":
		res = HardEqualOperation(values[0], values[1])
		break
	case "!=":
		res = NotSoftEqualOperation(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "!==":
		res = NotHardEqualOperation(values[0], values[1])
		break
	case ">":
		res = MoreOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case ">=":
		res = MoreEqualOperation(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "<":
		res = LessOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "<=":
		res = LessEqualOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "!":
		res = NotTruthyOperation(values)
		break
	case "!!":
		res = TruthyOperation(values)
		break
	case "%":
		res = PercentageOperation(cast.ToInt(values[0]), cast.ToInt(values[1]))
		break
	case "and":
		res = AndOperation(values)
		break
	case "or":
		res = OrOperation(values)
		break
	case "var":
		var fallback interface{}
		if len(values) > 1 {
			fallback = values[1]
		} else {
			fallback = nil
		}

		res = VarOperation(cast.ToString(values[0]), fallback, data)
		break
	case "?":
	case "if":
		res = IfOperation(cast.ToBool(values[0]), values[1], values[2])
		break
	case "log":

		res = LogOperation(cast.ToString(values[0]))
		break
	case "+":
		res = PlusOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "-":
		res = MinusOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "*":
		res = MultiplyOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "/":
		res = DivideOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "min":
		res = MinOperation(values)
		break
	case "max":
		res = MaxOperation(values)
		break
	}
	return res
}

// ParseObject entry point
func ParseObject(logic string, data string) (res interface{}, err error) {
	// result := interface{}
	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		res = RunOperation(string(key), string(value), data)
		// if operation, ok := Operations[string(key)]; ok {
		// 	res = operation.run(string(value), data)
		// } else {
		// 	return ErrInvalidOperation
		// }
		return nil
	})

	if err != nil {
		return false, err
	}
	return res, nil
}

// Max type is entry point for parser
type Max struct{}

func (o Max) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return MaxOperation(values)
}

// MaxOperation implements the 'Max' conditional returning the Maximum value from an array of values.
func MaxOperation(values []interface{}) (max float64) {
	if len(values) == 0 {
		return 0
	}

	max = cast.ToFloat64(values[0])
	for _, v := range values {
		val := cast.ToFloat64(v)
		if val > max {
			max = val
		}
	}
	return max
}

// Min type is entry point for parser
type Min struct{}

func (o Min) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return MinOperation(values)
}

// MinOperation implements the 'min' conditional returning the minimum value from an array of values.
func MinOperation(values []interface{}) (min float64) {
	if len(values) == 0 {
		return 0
	}

	min = cast.ToFloat64(values[0])
	for _, v := range values {
		val := cast.ToFloat64(v)
		if val < min {
			min = val
		}
	}

	return min
}

// Log type is entry point for parser
type Log struct{}

func (o Log) run(a ...interface{}) interface{} {
	values := GetValues(cast.ToString(a[0]), cast.ToString(a[1]))
	return LogOperation(cast.ToString(values[0]))
}

// LogOperation implements the 'log' operator, which prints a log inside termianl.
func LogOperation(a string) interface{} {
	fmt.Println(a)
	return nil
}

// Plus type is entry point for parser
type Plus struct{}

func (o Plus) run(a ...interface{}) interface{} {
	values := GetValues(cast.ToString(a[0]), cast.ToString(a[1]))
	return PlusOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// PlusOperation implements the '+' operator, which does type JS-style coertion.
func PlusOperation(a float64, b float64) interface{} {
	return a + b
}

// Minus type is entry point for parser
type Minus struct{}

func (o Minus) run(a ...interface{}) interface{} {
	values := GetValues(cast.ToString(a[0]), cast.ToString(a[1]))
	return MinusOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// MinusOperation implements the '-' operator, which does type JS-style coertion.
func MinusOperation(a float64, b float64) interface{} {
	return a - b
}

// Multiply type is entry point for parser
type Multiply struct{}

func (o Multiply) run(a ...interface{}) interface{} {
	values := GetValues(cast.ToString(a[0]), cast.ToString(a[1]))
	return MultiplyOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// MultiplyOperation implements the '-' operator, which does type JS-style coertion.
func MultiplyOperation(a float64, b float64) interface{} {
	return a * b
}

// Divide type is entry point for parser
type Divide struct{}

func (o Divide) run(a ...interface{}) interface{} {
	values := GetValues(cast.ToString(a[0]), cast.ToString(a[1]))
	return DivideOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// DivideOperation implements the '-' operator, which does type JS-style coertion.
func DivideOperation(a float64, b float64) interface{} {
	return a / b
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
func HardEqualOperation(a ...interface{}) bool {
	if GetType(a[0]) != GetType(a[1]) {
		return false
	}

	if a[0] == a[1] {
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
func NotHardEqualOperation(a ...interface{}) bool {
	return !HardEqualOperation(a[0], a[1])
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

// LessEqual type is entry point for parser
type LessEqual struct{}

func (o LessEqual) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return LessEqualOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
}

// LessEqualOperation implements the '<=' operator with JS-style type coertion. Returns bool.
func LessEqualOperation(a float64, b float64) bool {
	if a <= b {
		return true
	}
	return false
}

// NotTruthy type is entry point for parser
type NotTruthy struct{}

func (o NotTruthy) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return NotTruthyOperation(values)
}

// NotTruthyOperation implements the '!' operator with JS-style type coertion. Returns bool.
func NotTruthyOperation(a interface{}) bool {
	return !TruthyOperation(a)
}

// Truthy type is entry point for parser
type Truthy struct{}

func (o Truthy) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return TruthyOperation(values)
}

// TruthyOperation implements the '!!' operator with JS-style type coertion. Returns bool.
func TruthyOperation(a interface{}) bool {
	valid, length := isArray(a)
	if valid == true && length == 0 {
		return true
	}

	return cast.ToBool(a)
}

// Percentage type is entry point for parser
type Percentage struct{}

func (o Percentage) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return percent.PercentOf(cast.ToInt(values[0]), cast.ToInt(values[1]))
}

// PercentageOperation implements the '%' operator, which does type JS-style coertion. Returns float64.
func PercentageOperation(a int, b int) float64 {
	return percent.PercentOf(a, b)
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

// If type is entry point for parser
type If struct{}

func (o If) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))
	return IfOperation(cast.ToBool(values[0]), values[1], values[2])
}

// IfOperation implements the 'if' conditional where if the first value is true, the second value is returned, otherwise the third
func IfOperation(conditional bool, success interface{}, fail interface{}) interface{} {
	if conditional {
		return success
	}
	return fail
}

// GetValues will return values of any kind
func GetValues(logic string, data string) (results []interface{}) {

	_, err := jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			res, _ := ParseObject(string(value), data)
			results = append(results, res)
			break
		case jsonparser.String:
			results = append(results, cast.ToString(value))
			break
		case jsonparser.Number:
			results = append(results, cast.ToFloat64(cast.ToString(value)))
			break
		case jsonparser.Boolean:
			results = append(results, cast.ToBool(value))
			break
		case jsonparser.Null:
			results = append(results, value)
			break
		}
	})
	if err != nil {
		// Is this a string? or is it nil
		if logic != "" {
			if _, err := strconv.Atoi(logic); err == nil {
				value, dataType, _, _ := jsonparser.Get([]byte(data), "["+logic+"]")
				if len(value) > 0 {
					results = append(results, logic)
					switch dataType {
					case jsonparser.Object:
						res, _ := ParseObject(string(value), data)
						results = append(results, res)
						break
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
				} else {
					results = append(results, logic)
				}

			} else {
				results = append(results, logic)
			}
		} else {
			return nil
		}
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
		switch dataType {
		case jsonparser.Object:
			objectResult, _ := ParseObject(string(value), data)
			result[offset] = objectResult.(bool)
			break
		}
	})
	if err != nil {
		return nil, 0, err
	}
	return result, offset, nil
}

// Var type is entry point for parser
type Var struct{}

func (o Var) run(a ...interface{}) interface{} {
	values := GetValues(a[0].(string), a[1].(string))

	var fallback interface{}
	if len(values) > 1 {
		fallback = values[1]
	} else {
		fallback = nil
	}

	return VarOperation(cast.ToString(values[0]), fallback, cast.ToString(a[1]))
}

// VarOperation implements the 'var' operator, which grabs value from passed data and has a fallback
func VarOperation(logic string, fallback interface{}, data string) interface{} {
	key := strings.Split(logic, ".")
	dataValue, dataType, _, _ := jsonparser.Get([]byte(data), key...)
	value := TranslateType(dataValue, dataType)
	if value == nil {
		value = fallback
	}
	return value
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

// Apply is the entry function to parse logic and optional data
func Apply(logic string, data string) (res interface{}, errs error) {

	// Ensure data is object
	if data == `` {
		data = `{}`
	}

	// Must be an object to kick off process
	result, err := ParseObject(logic, data)
	if err != nil {
		return false, err
	}

	return result, nil
}

func isArray(args interface{}) (valid bool, length int) {
	val := reflect.ValueOf(args)

	if val.Kind() == reflect.Array {
		return true, val.Len()
	} else if val.Kind() == reflect.Slice {
		return true, val.Len()
	} else {
		return false, 0
	}
}

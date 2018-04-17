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

func RunOperation(key string, logic string, data string) (res interface{}) {
	values := GetValues(logic, data)
	switch key {
	case "==":
		res = SoftEqual(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "===":
		res = HardEqual(values[0], values[1])
		break
	case "!=":
		res = NotSoftEqual(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "!==":
		res = NotHardEqual(values[0], values[1])
		break
	case ">":
		res = More(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case ">=":
		res = MoreEqual(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "<":
		res = Less(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "<=":
		res = LessEqual(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "!":
		res = NotTruthy(values)
		break
	case "!!":
		res = Truthy(values)
		break
	case "%":
		res = Percentage(cast.ToInt(values[0]), cast.ToInt(values[1]))
		break
	case "and":
		res = And(values)
		break
	case "or":
		res = Or(values)
		break
	case "var":
		var fallback interface{}
		if len(values) > 1 {
			fallback = values[1]
		} else {
			fallback = nil
		}

		res = Var(cast.ToString(values[0]), fallback, data)
		break
	case "?":
	case "if":
		res = If(cast.ToBool(values[0]), values[1], values[2])
		break
	case "log":
		res = Log(cast.ToString(values[0]))
		break
	case "+":
		res = Plus(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "-":
		res = Minus(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "*":
		res = Multiply(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "/":
		res = Divide(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "min":
		res = Min(values)
		break
	case "max":
		res = Max(values)
		break
	}
	return res
}

// Max implements the 'Max' conditional returning the Maximum value from an array of values.
func Max(values []interface{}) (max float64) {
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

// Min implements the 'min' conditional returning the minimum value from an array of values.
func Min(values []interface{}) (min float64) {
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

// Log implements the 'log' operator, which prints a log inside termianl.
func Log(a string) interface{} {
	fmt.Println(a)
	return nil
}

// Plus implements the '+' operator, which does type JS-style coertion.
func Plus(a float64, b float64) interface{} {
	return a + b
}

// Minus implements the '-' operator, which does type JS-style coertion.
func Minus(a float64, b float64) interface{} {
	return a - b
}

// Multiply implements the '-' operator, which does type JS-style coertion.
func Multiply(a float64, b float64) interface{} {
	return a * b
}

// Divide implements the '-' operator, which does type JS-style coertion.
func Divide(a float64, b float64) interface{} {
	return a / b
}

// SoftEqual implements the '==' operator, which does type JS-style coertion. Returns bool.
func SoftEqual(a string, b string) bool {
	if a == b {
		return true
	}
	return false
}

// HardEqual Implements the '===' operator, which does type JS-style coertion. Returns bool.
func HardEqual(a ...interface{}) bool {
	if GetType(a[0]) != GetType(a[1]) {
		return false
	}

	if a[0] == a[1] {
		return true
	}

	return false
}

// NotSoftEqual implements the '!=' operator, which does type JS-style coertion. Returns bool.
func NotSoftEqual(a string, b string) bool {
	return !SoftEqual(a, b)
}

// NotHardEqual implements the '!==' operator, which does type JS-style coertion. Returns bool.
func NotHardEqual(a ...interface{}) bool {
	return !HardEqual(a[0], a[1])
}

// More implements the '>' operator with JS-style type coertion. Returns bool.
func More(a float64, b float64) bool {
	return LessEqual(b, a)
}

// MoreEqual implements the '>=' operator with JS-style type coertion. Returns bool.
func MoreEqual(a string, b string) bool {
	return Less(cast.ToFloat64(b), cast.ToFloat64(a)) || SoftEqual(a, b)
}

// Less implements the '<' operator with JS-style type coertion. Returns bool.
func Less(a float64, b float64) bool {
	if a < b {
		return true
	}
	return false
}

// LessEqual implements the '<=' operator with JS-style type coertion. Returns bool.
func LessEqual(a float64, b float64) bool {
	if a <= b {
		return true
	}
	return false
}

// NotTruthy implements the '!' operator with JS-style type coertion. Returns bool.
func NotTruthy(a interface{}) bool {
	return !Truthy(a)
}

// Truthy implements the '!!' operator with JS-style type coertion. Returns bool.
func Truthy(a interface{}) bool {
	valid, length := isArray(a)
	if valid == true && length == 0 {
		return true
	}

	return cast.ToBool(a)
}

// Percentage implements the '%' operator, which does type JS-style coertion. Returns float64.
func Percentage(a int, b int) float64 {
	return percent.PercentOf(a, b)
}

// And implements the 'and' conditional requiring all bubbled up bools to be true
func And(values []interface{}) bool {
	result := true
	for _, res := range values {
		if res == false {
			result = false
		}
	}
	return result
}

// Or implements the 'or' conditional requiring at least one of the bubbled up bools to be true
func Or(values []interface{}) bool {
	result := false
	for _, res := range values {
		if res == true {
			result = true
		}
	}
	return result
}

// If implements the 'if' conditional where if the first value is true, the second value is returned, otherwise the third
func If(conditional bool, success interface{}, fail interface{}) interface{} {
	if conditional {
		return success
	}
	return fail
}

// Var implements the 'var' operator, which grabs value from passed data and has a fallback
func Var(logic string, fallback interface{}, data string) interface{} {
	key := strings.Split(logic, ".")
	dataValue, dataType, _, _ := jsonparser.Get([]byte(data), key...)
	value := TranslateType(dataValue, dataType)
	if value == nil {
		value = fallback
	}
	return value
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

// ParseObject entry point
func ParseObject(logic string, data string) (res interface{}, err error) {
	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		res = RunOperation(string(key), string(value), data)
		return nil
	})

	if err != nil {
		return false, err
	}
	return res, nil
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

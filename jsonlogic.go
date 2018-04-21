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

// Run is an alias to Apply without data
func Run(rule string) (res interface{}, errs error) {
	return Apply(rule, ``)
}

// Apply is the entry function to parse rule and optional data
func Apply(rule string, data string) (res interface{}, errs error) {

	// Ensure data is object
	if data == `` {
		data = `{}`
	}

	// Must be an object to start process
	result, err := ParseOperator(rule, data)
	if err != nil {
		return false, err
	}

	return result, nil
}

// ParseOperator takes in the json rule and data and attempts to parse
func ParseOperator(rule string, data string) (result interface{}, err error) {
	err = jsonparser.ObjectEach([]byte(rule), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		result = RunOperator(string(key), string(value), data)
		return nil
	})

	if err != nil {
		return false, err
	}

	return result, nil
}

// GetValues will attempt to recursively resolve all values for a given operator
func GetValues(rule string, data string) (results []interface{}) {

	// Jsonrule rule is always one key, with an array of values
	_, err := jsonparser.ArrayEach([]byte(rule), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			res, _ := ParseOperator(string(value), data)
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

	// Jsonrule may also support syntactic sugar
	if err != nil {
		if rule != "" {
			if _, err := strconv.Atoi(rule); err == nil {

				// If string then we can attempt to retrieve the value from the data
				value, dataType, _, _ := jsonparser.Get([]byte(data), "["+rule+"]")
				if len(value) > 0 {
					results = append(results, rule)
					switch dataType {
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
				} else {
					// No data was found so we just append the rule and move on
					results = append(results, rule)
				}

			} else {
				// Is an integer so we assume it's value
				results = append(results, rule)
			}
		} else {
			return nil
		}
	}

	return results
}

// RunOperator determines what function to run against the passed rule and data
func RunOperator(key string, rule string, data string) (result interface{}) {
	values := GetValues(rule, data)
	switch key {
	case "==":
		result = SoftEqual(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "===":
		result = HardEqual(values[0], values[1])
		break
	case "!=":
		result = NotSoftEqual(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "!==":
		result = NotHardEqual(values[0], values[1])
		break
	case ">":
		result = More(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case ">=":
		result = MoreEqual(cast.ToString(values[0]), cast.ToString(values[1]))
		break
	case "<":
		result = Less(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "<=":
		result = LessEqual(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "!":
		result = NotTruthy(values)
		break
	case "!!":
		result = Truthy(values)
		break
	case "%":
		result = Percentage(cast.ToInt(values[0]), cast.ToInt(values[1]))
		break
	case "and":
		result = And(values)
		break
	case "or":
		result = Or(values)
		break
	case "var":
		var fallback interface{}
		if len(values) > 1 {
			fallback = values[1]
		} else {
			fallback = nil
		}

		result = Var(cast.ToString(values[0]), fallback, data)
		break
	case "?":
	case "if":
		result = If(cast.ToBool(values[0]), values[1], values[2])
		break
	case "log":
		result = Log(cast.ToString(values[0]))
		break
	case "+":
		result = Plus(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "-":
		result = Minus(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "*":
		result = Multiply(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "/":
		result = Divide(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		break
	case "min":
		result = Min(values)
		break
	case "max":
		result = Max(values)
		break
	}
	return result
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

// And implements the 'and' conditional requiring all bubbled up bools to be true.
func And(values []interface{}) bool {
	result := true
	for _, res := range values {
		if res == false {
			result = false
		}
	}
	return result
}

// Or implements the 'or' conditional requiring at least one of the bubbled up bools to be true.
func Or(values []interface{}) bool {
	result := false
	for _, res := range values {
		if res == true {
			result = true
		}
	}
	return result
}

// If implements the 'if' conditional where if the first value is true, the second value is returned, otherwise the third.
func If(conditional bool, success interface{}, fail interface{}) interface{} {
	if conditional {
		return success
	}
	return fail
}

// Var implements the 'var' operator, which grabs value from passed data and has a fallback.
func Var(rule string, fallback interface{}, data string) interface{} {
	key := strings.Split(rule, ".")
	dataValue, dataType, _, _ := jsonparser.Get([]byte(data), key...)
	value := TranslateType(dataValue, dataType)
	if value == nil {
		value = fallback
	}
	return value
}

// GetType returns an int to map against type so we can see if we are dealing with a specific type of data or an object operation.
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

// TranslateType Takes the returned dataType from jsonparser along with it's returned []byte data and returns the casted value.
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

// isArray is a simple function to determine if passed args is of type array.
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

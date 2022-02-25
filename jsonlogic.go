package jsonlogic

import (
	"bytes"
	"encoding/json"
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
	ErrInvalidOperation = "invalid operation: %s"
)

// Operators holds any operators
var Operators = make(map[string]func(rule string, data string) (result interface{}))

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

	// Unicode &
	data = strings.ReplaceAll(data, `\u0026`, `&`)

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
		switch dataType {
		case jsonparser.String:
			result = RunOperator(string(key), "\""+string(value)+"\"", data)
		default:
			result = RunOperator(string(key), string(value), data)
		}
		return nil
	})

	if err != nil {
		return false, fmt.Errorf(ErrInvalidOperation, err)
	}

	return result, nil
}

// GetValues will attempt to recursively resolve all values for a given operator
func GetValues(rule string, data string) (results []interface{}) {

	ruleValue, dataType, _, _ := jsonparser.Get([]byte(rule))
	switch dataType {
	case jsonparser.Object:
		res, _ := ParseOperator(string(ruleValue), data)
		results = append(results, res)
	case jsonparser.Array:
		jsonparser.ArrayEach([]byte(ruleValue), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			switch dataType {
			case jsonparser.Array:
				m := make([]interface{}, 0)
				json.Unmarshal(value, &m)
				results = append(results, m)
			case jsonparser.Object:
				res, _ := ParseOperator(string(value), data)
				results = append(results, res)
			case jsonparser.String:
				results = append(results, cast.ToString(value))
			case jsonparser.Number:
				results = append(results, cast.ToFloat64(cast.ToString(value)))
			case jsonparser.Boolean:
				results = append(results, cast.ToBool(string(value)))
			case jsonparser.Null:
				results = append(results, value)
			}
		})
	case jsonparser.Number:
		results = append(results, cast.ToFloat64(string(ruleValue)))
	case jsonparser.String:
		// Remove the quotes we added so we could detect string type
		rule = rule[1 : len(rule)-1]
		value, dataType, _, _ := jsonparser.Get([]byte(data), rule)
		if len(value) > 0 {
			results = append(results, rule)
			switch dataType {
			case jsonparser.String:
				results = append(results, cast.ToString(value))
			case jsonparser.Number:
				results = append(results, cast.ToFloat64(cast.ToString(value)))
			case jsonparser.Boolean:
				results = append(results, cast.ToBool(value))
			case jsonparser.Null:
				results = append(results, value)
			}
		} else {
			// No data was found so we just append the rule and move on
			results = append(results, rule)
		}
	default:
		return nil
	}

	return results
}

// AddOperator allows for custom operators to be used
func AddOperator(key string, cb func(rule string, data string) (result interface{})) {
	Operators[key] = cb
}

// RunOperator determines what function to run against the passed rule and data
func RunOperator(key string, rule string, data string) (result interface{}) {

	values := GetValues(rule, data)
	switch key {
	// Accessing Data
	case "var":
		var fallback interface{}
		if len(values) > 1 {
			fallback = values[1]
		} else {
			fallback = nil
		}

		result = Var(values[0], fallback, data)
	// TODO missing
	case "missing":
		result = Missing(values, data)
	case "missing_some":

		break
	// TODO missing_some
	// Logic and Boolean Operations
	case "?":
	case "if":
		// TOFIX basically the "success" value is showing false when it should be showing true
		// result = If(values[0], values[1], values[2])
		result = If(values)
	case "==":
		result = SoftEqual(cast.ToString(values[0]), cast.ToString(values[1]))
	case "===":
		result = HardEqual(values[0], values[1])
	case "!=":
		result = NotSoftEqual(cast.ToString(values[0]), cast.ToString(values[1]))
	case "!==":
		result = NotHardEqual(values[0], values[1])
	case "!":
		result = NotTruthy(values)
	case "!!":
		result = Truthy(values)
	case "or":
		result = Or(values)
	case "and":
		result = And(values)
		// Numeric Operations
	case ">":
		result = More(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
	case ">=":
		result = MoreEqual(cast.ToString(values[0]), cast.ToString(values[1]))
	case "<":
		// Test for exclusive between
		result = false
		if len(values) > 2 && IsNumeric(values[0]) && IsNumeric(values[1]) && IsNumeric(values[2]) {
			result = LessBetween(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]), cast.ToFloat64(values[2]))
		} else if IsNumeric(values[0]) && IsNumeric(values[1]) {
			result = Less(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		}
	case "<=":
		// Test for inclusive between
		result = false
		if len(values) > 2 && IsNumeric(values[0]) && IsNumeric(values[1]) && IsNumeric(values[2]) {
			result = LessEqualBetween(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]), cast.ToFloat64(values[2]))
		} else if IsNumeric(values[0]) && IsNumeric(values[1]) {
			result = LessEqual(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
		}
	case "max":
		result = Max(values)
	case "min":
		result = Min(values)
	case "+":
		result = Plus(values)
	case "-":
		result = Minus(values)
	case "*":
		result = Multiply(values)
	case "/":
		result = Divide(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
	case "%":
		result = Percentage(cast.ToInt(values[0]), cast.ToInt(values[1]))
		// String Operations
	case "cat":
		result = Cat(values)
	case "in":
		result = In(values)
	case "substr":
		if len(values) > 2 {
			result = Substr(cast.ToString(values[0]), cast.ToInt(values[1]), cast.ToInt(values[2]))
		} else {
			result = Substr(cast.ToString(values[0]), cast.ToInt(values[1]), 0)
		}
	case "merge":
		result = Merge(values)
		// TODO All, None and Some http://jsonlogic.com/operations.html#all-none-and-some
	case "all":

		break
	case "some":

		break
	case "none":

		break
		// TODO Map, Reduce and Filter http://jsonlogic.com/operations.html#map-reduce-and-filter
	case "map":

		break
	case "reduce":

		break
	case "filter":

		break
		// Miscellaneous
	case "log":
		result = Log(cast.ToString(values[0]))
	}

	// Check against any custom operators
	for index, operation := range Operators {
		if key == index {
			result = operation(rule, data)
		}
	}

	return result
}

func IsNumeric(s interface{}) bool {
	_, err := strconv.ParseFloat(cast.ToString(s), 64)
	return err == nil
}

func Missing(a []interface{}, data string) interface{} {
	result := make([]interface{}, 0)

	for i := 0; i < len(a); i++ {
		_, dataType, _, _ := jsonparser.Get([]byte(data), cast.ToString(a[i]))
		if dataType == jsonparser.NotExist {
			result = append(result, a[i])
		}
	}

	return result
}

func Merge(a []interface{}) interface{} {
	result := make([]interface{}, 0)

	for i := 0; i < len(a); i++ {
		array, _ := isArray(a[i])
		if array {
			item := a[i].([]interface{})
			for x := 0; x < len(item); x++ {
				result = append(result, item[x])
			}
		} else {
			result = append(result, a[i])
		}
	}

	return result
}

func Substr(a string, position int, length int) string {
	start := 0
	end := len(a)

	if position < 0 {
		start = end + position
	} else {
		start = position
	}

	if length < 0 {
		end = end + length
	} else if length > 0 {
		end = position + length
	}

	return a[start:end]
}

func In(a []interface{}) bool {
	array, _ := isArray(a[1])
	if array {
		items := a[1].([]interface{})
		result := false
		for i := 0; i < len(items); i++ {
			if strings.Contains(cast.ToString(items[i]), cast.ToString(a[0])) && a[0] != nil {
				result = true
			}
		}
		return result
	}
	return strings.Contains(cast.ToString(a[1]), cast.ToString(a[0]))
}

// Cat implements the 'cat' conditional returning all the values merged together.
func Cat(values []interface{}) string {
	buffer := new(bytes.Buffer)
	for _, v := range values {
		buffer.WriteString(cast.ToString(v))
	}
	return buffer.String()
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
func Plus(a []interface{}) interface{} {
	result := 0.0

	for _, v := range a {
		result = result + cast.ToFloat64(v)
	}

	return result
}

// Minus implements the '-' operator, which does type JS-style coertion.
func Minus(a []interface{}) interface{} {
	result := cast.ToFloat64(a[0])

	if len(a) < 2 {
		result = -1 * cast.ToFloat64(a[0])
	} else {
		for i, v := range a {
			if i != 0 {
				result = result - cast.ToFloat64(v)
			}
		}
	}

	return result
}

// Multiply implements the '-' operator, which does type JS-style coertion.
func Multiply(a []interface{}) interface{} {
	result := 1.0

	for _, v := range a {
		result = result * cast.ToFloat64(v)
	}

	return result
}

// Divide implements the '-' operator, which does type JS-style coertion.
func Divide(a float64, b float64) interface{} {
	return a / b
}

// SoftEqual implements the '==' operator, which does type JS-style coertion.
func SoftEqual(a string, b string) bool {
	return a == b
}

// HardEqual Implements the '===' operator, which does type JS-style coertion.
func HardEqual(a ...interface{}) bool {
	if GetType(a[0]) != GetType(a[1]) {
		return false
	}

	if a[0] == a[1] {
		return true
	}

	return false
}

// NotSoftEqual implements the '!=' operator, which does type JS-style coertion.
func NotSoftEqual(a string, b string) bool {
	return !SoftEqual(a, b)
}

// NotHardEqual implements the '!==' operator, which does type JS-style coertion.
func NotHardEqual(a ...interface{}) bool {
	return !HardEqual(a[0], a[1])
}

// More implements the '>' operator with JS-style type coertion.
func More(a float64, b float64) bool {
	return LessEqual(b, a)
}

// MoreEqual implements the '>=' operator with JS-style type coertion.
func MoreEqual(a string, b string) bool {
	return Less(cast.ToFloat64(b), cast.ToFloat64(a)) || SoftEqual(a, b)
}

// Less implements the '<' operator however checks against 3 values to test that one value is between but not equal to two others.
func LessBetween(a float64, b float64, c float64) bool {
	leftCheck := Less(a, b)
	rightCheck := Less(b, c)

	if leftCheck && rightCheck {
		return true
	}
	return false
}

// Less implements the '<' operator with JS-style type coertion.
func Less(a float64, b float64) bool {
	return a < b
}

// Less implements the '<' operator however checks against 3 values to test that one value is between two others.
func LessEqualBetween(a float64, b float64, c float64) bool {
	leftCheck := LessEqual(a, b)
	rightCheck := LessEqual(b, c)

	if leftCheck && rightCheck {
		return true
	}
	return false
}

// LessEqual implements the '<=' operator with JS-style type coertion.
func LessEqual(a float64, b float64) bool {
	return a <= b
}

// NotTruthy implements the '!' operator with JS-style type coertion.
func NotTruthy(a interface{}) bool {
	return !Truthy(a)
}

// Truthy implements the '!!' operator with JS-style type coertion.
func Truthy(a interface{}) bool {
	valid, length := isArray(a)
	if valid && length == 0 {
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
// func If(conditional interface{}, success interface{}, fail interface{}) interface{} {
func If(conditions []interface{}) interface{} {
	var result interface{}

	lastElement := conditions[len(conditions)-1]
	isTrue := false

	for i := 0; i < len(conditions); i++ {

		if (i + 1) < len(conditions) {
			value := conditions[i+1]
			if cast.ToBool(conditions[i]) {
				result = value
				isTrue = true
			}
		}

		i++
	}

	if isTrue {
		return result
	}

	return lastElement
}

// Var implements the 'var' operator, which grabs value from passed data and has a fallback.
func Var(rules interface{}, fallback interface{}, data string) (value interface{}) {
	ruleType := GetType(rules)
	rule := ""
	switch ruleType {
	case 1:
	case 2:
		rule = "[" + cast.ToString(rules) + "]"
	default:
		rule = cast.ToString(rules)
	}

	if cast.ToString(rules) == "" {
		dataValue, dataType, _, _ := jsonparser.Get([]byte(data))
		if dataType != jsonparser.NotExist {
			value = TranslateType(dataValue, dataType)
		}
	} else {
		key := strings.Split(rule, ".")
		dataValue, dataType, _, _ := jsonparser.Get([]byte(data), key...)
		value = TranslateType(dataValue, dataType)
		if value == nil {
			value = fallback
		}
	}

	if value == "" {
		value = data
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
	case jsonparser.Number:
		numberString := cast.ToString(data)
		numberFloat := cast.ToFloat64(numberString)
		return numberFloat
	case jsonparser.Boolean:
		return string(data)
	case jsonparser.Null:
		return string(data)
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

# json-logic-go

This parser accepts [JsonLogic](http://jsonlogic.com) rules and executes them in GO.

The JsonLogic format is designed to allow you to share rules (logic) between front-end and back-end code (regardless of language difference), even to store logic along with a record in a database.  JsonLogic is documented extensively at [JsonLogic.com](http://jsonlogic.com), including examples of every [supported operation](http://jsonlogic.com/operations.html) and a place to [try out rules in your browser](http://jsonlogic.com/play.html).

The same format can also be executed in JavaScript by the library [json-logic-js](https://github.com/jwadhams/json-logic-js/)

The same format can also be executed in PHP by the library [json-logic-php](https://github.com/jwadhams/json-logic-php/)

## Examples

### A note about types

This is a GO interpreter of a format designed to be transmitted and stored as JSON.  So it makes sense to conceptualize the rules in JSON.

Expressed in JSON, a JsonLogic rule is always one key, with an array of values.

```GO
logic = `{"==":["apples", "apples"]}`
result, err := jsonlogic.Apply(logic, ``)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)
// true
```

### Simple
```GO
logic = `{"==":[1, 1]}`
result, err := jsonlogic.Apply(logic, ``)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)
// true
```

This is a simple test, equivalent to `1 == 1`.  A few things about the format:

  1. The operator is always in the "key" position. There is only one key per JsonLogic rule.
  1. The values are typically an array.
  1. Each value can be a string, number, boolean, array, or null

### Compound
Here we're beginning to nest rules. 

```GO
logic = `{"and": [
		{ ">" => [3,1] },
		{ "<" => [1,3] }
	] }`
result, err := jsonlogic.Apply(logic, ``)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)
// true
```
    
### Data-Driven

Obviously these rules aren't very interesting if they can only take static literal data. Typically `JsonLogic::apply` will be called with a rule object and a data object. You can use the `var` operator to get attributes of the data object:

```GO
logic = `{ "var" => ["a"] }`
data = `{ "a" => 1, "b" => 2 }`
result, err := jsonlogic.Apply(logic, data)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)
// 1
```

If you like, we support [syntactic sugar](https://en.wikipedia.org/wiki/Syntactic_sugar) on unary operators to skip the array around values:

You can also use the `var` operator to access an array by numeric index:

```GO
logic = `{ "var": [1] }`
data = `{ "apple", "banana", "carrot" }`
result, err := jsonlogic.Apply(logic, data)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)
// "banana"
```

Here's a complex rule that mixes literals and data. The pie isn't ready to eat unless it's cooler than 110 degrees, *and* filled with apples.

```GO
logic = `{ "and": [
	{ "<": [ { "var": "temp" }, 110 ] },
	{ "==": [ { "var": "pie.filling" }, "apple" ] }
] }`
data = `{ "temp": 100, "pie": { "filling": "apple" } }`
result, err := jsonlogic.Apply(logic, data)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)
// true
```
    
## Installation

```
go get github.com/GeorgeD19/json-logic-go
```
package main

import (
	"fmt"

	jsonlogic "github.com/GeorgeD19/json-logic-go"
)

func main() {
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

	result, _ := jsonlogic.Apply(rule, data)

	fmt.Printf("rule returned %s", result)
}

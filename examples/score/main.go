package main

import (
	"fmt"

	jsonlogic "github.com/GeorgeD19/json-logic-go"
)

func main() {
	rule := `
	{
		"if": [
			{"<": [{"var":"question_1_score"}, 4] }, true,
			{"<": [{"var":"question_2_score"}, 4] }, true,
			{"<": [{"var":"question_3_score"}, 4] }, true,
			{"<": [{"var":"question_4_score"}, 4] }, true,
			{"<": [{"var":"question_5_score"}, 4] }, true,
			{"<": [{"var":"question_6_score"}, 4] }, true,
			{"<": [{"var":"question_7_score"}, 4] }, true,
			{"<": [{"var":"question_8_score"}, 4] }, true,
			{"<": [{"var":"question_9_score"}, 4] }, true,
			{"<": [{"var":"question_10_score"}, 4] }, true,
			{"<": [{"var":"question_11_score"}, 4] }, true,
			false
		]
	}	
	`
	data := `{
		"question_1_score": "4",
		"question_2_score": "4",
		"question_3_score": "4",
		"question_4_score": "4",
		"question_5_score": "4",
		"question_6_score": "4",
		"question_7_score": "4",
		"question_8_score": "4",
		"question_9_score": "4",
		"question_10_score": "4",
		"question_11_score": "3"
	}`

	result, _ := jsonlogic.Apply(rule, data)

	fmt.Printf("rule returned %s", result)
}

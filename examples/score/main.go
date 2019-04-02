package main

import (
	"fmt"

	"github.com/spf13/cast"

	jsonlogic "github.com/GeorgeD19/json-logic-go"
)

func main() {
	rule := `
	{
		"if": [{
			"and": [{
					"<": [{
						"var": "question_1_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_1_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_2_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_2_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_3_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_3_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_4_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_4_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_5_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_5_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_6_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_6_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_7_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_7_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_8_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_8_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_9_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_9_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_10_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_10_score"
						},
						"N/A"
					]
				}
			]
		}, true, {
			"and": [{
					"<": [{
						"var": "question_11_score"
					}, 3]
				},
				{
					"!=": [{
							"var": "question_11_score"
						},
						"N/A"
					]
				}
			]
		}, true, false]
	}
	`
	data := `{
		"question_1_score": "4",
		"question_2_score": "4",
		"question_3_score": "4",
		"question_4_score": "4",
		"question_5_score": "N/A",
		"question_6_score": "4",
		"question_7_score": "4",
		"question_8_score": "4",
		"question_9_score": "4",
		"question_10_score": "4",
		"question_11_score": "2"
	}`

	result, _ := jsonlogic.Apply(rule, data)

	fmt.Printf("rule returned %s", result)
	if cast.ToBool(result) == true {
		fmt.Println("Running [{\"create_ors\":{}},{\"email\":{\"to\":\"russel.kerr@securigroup.co.uk;allan.burnett@securigroup.co.uk;david.wilson@securigroup.co.uk;andrew.mcallister@securigroup.co.uk;peter.nokes@sglsecurity.co.uk;Biagio.Paciolla@securigroup.co.uk;susan.fitzpatrick@securigroup.co.uk;Jonathan.Greenlees@securigroup.co.uk;susanne.scott@securigroup.co.uk;incident.management@securigroup.co.uk\"}}]")
	}

	rule = `
	{
		"if": [{
			"<": [{
				"var": "question_10_score"
			}, 3]
		}, true, false]
	}
	`

	data = `{
		"question_1_score": "2",
		"question_2_score": "4",
		"question_3_score": "4",
		"question_4_score": "4",
		"question_5_score": "4",
		"question_6_score": "4",
		"question_7_score": "4",
		"question_8_score": "4",
		"question_9_score": "2",
		"question_10_score": "2",
		"question_11_score": "3"
	}`

	result, _ = jsonlogic.Apply(rule, data)

	fmt.Printf("rule returned %s", result)
	if cast.ToBool(result) == true {
		fmt.Println("Running [{\"email\":{\"to\":\"susan.fitzpatrick@securigroup.co.uk;jonathan.greenlees@securigroup.co.uk\"}}]")
	}

}

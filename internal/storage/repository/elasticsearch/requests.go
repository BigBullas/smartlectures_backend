package elasticsearch

import "fmt"

func buildSearchNoteReqBody(query string, userId int) map[string]interface{} {
	// return map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"multi_match": map[string]interface{}{
	// 			"query":     query,
	// 			"fields":    []string{"name", "body"},
	// 			"fuzziness": "AUTO",
	// 		},
	// 	},
	// 	"highlight": map[string]interface{}{
	// 		"fields": map[string]interface{}{
	// 			"name": map[string]interface{}{
	// 				"fragment_size": 50,
	// 			},
	// 			"body": map[string]interface{}{
	// 				"fragment_size": 100,
	// 			},
	// 		},
	// 	},
	// }

	fmt.Println("elastic------------",userId,"------------")
	fmt.Println("elastic2------------",query,"------------")

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					// map[string]interface{}{
					// 	"term": map[string]interface{}{
					// 		"userId": userId,
					// 	},
					// },
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":     query,
							"fields":    []string{"name", "body"},
							"fuzziness": "AUTO",
						},
					},
				},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"name": map[string]interface{}{
					"fragment_size": 50,
				},
				"body": map[string]interface{}{
					"fragment_size": 100,
				},
			},
		},
	}
}

func buildDeleteNotesByDirReqBody(dirId int) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"parentDir": dirId,
			},
		},
	}
}

func buildSearchSnippetReqBody(query string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"name": "*" + query + "*",
			},
		},
	}
}

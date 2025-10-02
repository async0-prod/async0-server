package utils

import "encoding/json"

func PrettyJSON(v interface{}) string {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(prettyJSON)
}

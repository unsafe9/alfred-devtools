package main

import "encoding/json"

func jsonFormatter(indent string) func(string) string {
	return func(input string) string {
		var parsed interface{}
		if err := json.Unmarshal([]byte(input), &parsed); err != nil {
			return "Invalid JSON input"
		}
		var formatted []byte
		if indent == "" {
			formatted, _ = json.Marshal(parsed)
		} else {
			formatted, _ = json.MarshalIndent(parsed, "", indent)
		}
		return string(formatted)
	}
}

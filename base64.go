package main

import "encoding/base64"

func encodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func decodeBase64(input string) string {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "Invalid base64 input"
	}
	return string(decoded)
}

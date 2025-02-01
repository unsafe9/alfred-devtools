package main

import (
	"encoding/base64"
	"strings"
)

func decodeJWT(input string) string {
	parts := strings.Split(input, ".")
	if len(parts) != 3 {
		return "Invalid JWT format"
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "Invalid JWT payload"
	}
	return string(payload)
}

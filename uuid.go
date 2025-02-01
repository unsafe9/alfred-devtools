package main

import "github.com/google/uuid"

func generateUUID(input string) string {
	return uuid.NewString()
}

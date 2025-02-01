package main

import (
	"encoding/hex"
	"hash"
)

func hashEncoder(hasher hash.Hash) func(string) string {
	return func(input string) string {
		hasher.Reset()
		hasher.Write([]byte(input))
		return hex.EncodeToString(hasher.Sum(nil))
	}
}

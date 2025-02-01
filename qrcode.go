package main

import (
	"github.com/skip2/go-qrcode"
	"os"
)

func qrcodeGenerator(level int) func(string) string {
	return func(input string) string {
		png, err := qrcode.Encode(input, qrcode.RecoveryLevel(level), level*256)
		if err != nil {
			return "Failed to generate QR code"
		}
		f, err := os.CreateTemp("", "alfred_devtools_qrcode_*.png")
		if err != nil {
			return "Failed to create QR code file"
		}
		defer f.Close()
		if _, err := f.Write(png); err != nil {
			return "Failed to write QR code file"
		}
		return f.Name()
	}
}

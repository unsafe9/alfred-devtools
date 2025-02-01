package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"strconv"
)

const (
	SubcommandTypeText = iota
	SubcommandTypeImage
)

type Command struct {
	Keyword     string
	Subtitle    string
	Subcommands []Subcommand
}

type Subcommand struct {
	Title   string
	Handler func(string) string
	Type    int
}

var commands = []Command{
	{
		Keyword:  "base64",
		Subtitle: "Encode/Decode Base64",
		Subcommands: []Subcommand{
			{
				Title:   "encode",
				Handler: encodeBase64,
			},
			{
				Title:   "decode",
				Handler: decodeBase64,
			},
		},
	},
	{
		Keyword:  "hash",
		Subtitle: "Hash Text",
		Subcommands: []Subcommand{
			{
				Title:   "md5",
				Handler: hashEncoder(md5.New()),
			},
			{
				Title:   "sha1",
				Handler: hashEncoder(sha1.New()),
			},
			{
				Title:   "sha256",
				Handler: hashEncoder(sha256.New()),
			},
		},
	},
	{
		Keyword:  "json",
		Subtitle: "Format JSON",
		Subcommands: []Subcommand{
			{
				Title:   "minify",
				Handler: jsonFormatter(""),
			},
			{
				Title:   "2-space",
				Handler: jsonFormatter("  "),
			},
			{
				Title:   "4-space",
				Handler: jsonFormatter("    "),
			},
			{
				Title:   "tabs",
				Handler: jsonFormatter("\t"),
			},
		},
	},
	{
		Keyword:  "uuid",
		Subtitle: "Generate UUID",
		Subcommands: []Subcommand{
			{
				Title:   "generate",
				Handler: generateUUID,
			},
		},
	},
	{
		Keyword:  "lorem_ipsum",
		Subtitle: "Generate lorem ipsum text",
		Subcommands: []Subcommand{
			{
				Title:   "words",
				Handler: loremIpsum.Words,
			},
			{
				Title:   "sentences",
				Handler: loremIpsum.Sentences,
			},
			{
				Title:   "paragraphs",
				Handler: loremIpsum.Paragraphs,
			},
		},
	},
	{
		Keyword:  "jwt",
		Subtitle: "Decode JWT",
		Subcommands: []Subcommand{
			{
				Title:   "decode",
				Handler: decodeJWT,
			},
		},
	},
	{
		Keyword:  "qrcode",
		Subtitle: "Generate QR code",
		Subcommands: []Subcommand{
			{
				Title:   "medium",
				Handler: qrcodeGenerator(1),
				Type:    SubcommandTypeImage,
			},
			{
				Title:   "high",
				Handler: qrcodeGenerator(2),
				Type:    SubcommandTypeImage,
			},
			{
				Title:   "highest",
				Handler: qrcodeGenerator(3),
				Type:    SubcommandTypeImage,
			},
		},
	},
}

func atoi(input string, def int) int {
	n, err := strconv.Atoi(input)
	if err != nil {
		return def
	}
	return n
}

package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/deanishe/awgo"
	"github.com/skip2/go-qrcode"
	"hash"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var wf *aw.Workflow

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
				Handler: generateLoremIpsumWords,
			},
			{
				Title:   "sentences",
				Handler: generateLoremIpsumSentences,
			},
			{
				Title:   "paragraphs",
				Handler: generateLoremIpsumParagraphs,
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

func main() {
	wf = aw.New()
	wf.Run(run)
}

func run() {
	query := ""
	if args := wf.Args(); len(args) > 0 {
		query = args[0]
	}
	parts := strings.SplitN(query, " ", 2)
	command := parts[0]
	found := false

	if command != "" {
		for _, cmd := range commands {
			if cmd.Keyword == command {
				found = true
				input := ""
				if len(parts) > 1 {
					input = parts[1]
				} else {
					input = getClipboardText()
				}

				for _, subcmd := range cmd.Subcommands {
					result := subcmd.Handler(input)
					item := wf.NewItem(subcmd.Title).
						Subtitle(result).
						Arg(result).
						Autocomplete(subcmd.Title).
						Valid(true)
					if subcmd.Type == SubcommandTypeImage {
						icon := &aw.Icon{Value: result, Type: aw.IconTypeImage}
						item.IsFile(true).Icon(icon)
					}
				}
				break
			}
		}
	}

	if !found {
		for _, cmd := range commands {
			wf.NewItem(cmd.Keyword).
				Subtitle(cmd.Subtitle).
				Arg(cmd.Keyword).
				Autocomplete(cmd.Keyword).
				Valid(true)
		}
	}

	wf.SendFeedback()
}

func getClipboardText() string {
	cmd := exec.Command("pbpaste")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

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

func hashEncoder(hasher hash.Hash) func(string) string {
	return func(input string) string {
		hasher.Reset()
		hasher.Write([]byte(input))
		return hex.EncodeToString(hasher.Sum(nil))
	}
}

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

func generateUUID(input string) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate UUID: %v", err))
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func generateLoremIpsumWords(input string) string {
	n, err := strconv.Atoi(input)
	if err != nil {
		n = 10
	}
	return LoremIpsum.Words(n)
}

func generateLoremIpsumSentences(input string) string {
	n, err := strconv.Atoi(input)
	if err != nil {
		n = 2
	}
	return LoremIpsum.Sentences(n)
}

func generateLoremIpsumParagraphs(input string) string {
	n, err := strconv.Atoi(input)
	if err != nil {
		n = 1
	}
	return LoremIpsum.Paragraphs(n)
}

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

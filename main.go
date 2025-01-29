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
	"hash"
	"os/exec"
	"strconv"
	"strings"
)

var wf *aw.Workflow

type Command struct {
	Keyword     string
	Subtitle    string
	Subcommands []Subcommand
}

type Subcommand struct {
	Title   string
	Handler func(string) string
}

var commands = []Command{
	{
		Keyword:  "base64",
		Subtitle: "Encode/Decode Base64",
		Subcommands: []Subcommand{
			{"encode", encodeBase64},
			{"decode", decodeBase64},
		},
	},
	{
		Keyword:  "hash",
		Subtitle: "Hash Text",
		Subcommands: []Subcommand{
			{"md5", hashEncoder(md5.New())},
			{"sha1", hashEncoder(sha1.New())},
			{"sha256", hashEncoder(sha256.New())},
		},
	},
	{
		Keyword:  "json",
		Subtitle: "Format JSON",
		Subcommands: []Subcommand{
			{"Minify", jsonFormatter("")},
			{"2-space", jsonFormatter("  ")},
			{"4-space", jsonFormatter("    ")},
			{"Tabs", jsonFormatter("\t")},
		},
	},
	{
		Keyword:  "uuid",
		Subtitle: "Generate UUID",
		Subcommands: []Subcommand{
			{"Generate", generateUUID},
		},
	},
	{
		Keyword:  "lorem_ipsum",
		Subtitle: "Generate lorem ipsum text",
		Subcommands: []Subcommand{
			{"Words", generateLoremIpsumWords},
			{"Sentences", generateLoremIpsumSentences},
			{"Paragraphs", generateLoremIpsumParagraphs},
		},
	},
	{
		Keyword:  "jwt",
		Subtitle: "Decode JWT",
		Subcommands: []Subcommand{
			{"Decode", decodeJWT},
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
					wf.NewItem(subcmd.Title).
						Subtitle(result).
						Arg(result).
						Autocomplete(subcmd.Title).
						Valid(true)
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

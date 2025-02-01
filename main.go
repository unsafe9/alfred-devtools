package main

import (
	"github.com/deanishe/awgo"
	"os/exec"
	"strings"
)

var wf *aw.Workflow

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

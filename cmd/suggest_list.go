package main

import (
	goterm "github.com/c-bata/go-prompt"
)

var suggests = []goterm.Suggest{
	{Text: "chatmode", Description: "Set terminal as a pure chat robot mode"},
	{Text: "cmdmode", Description: "Set terminal as a command mode, use natural language to communicate"},
	{Text: "taskmode", Description: "Set terminal as a task mode, use natural language to communicate to execute tasks"},
	{Text: "exit", Description: "Exit the terminal"},
}

func AddSuggest(text string, description string) {
	// Check if text not exist in suggests, then add it
	for _, suggest := range suggests {
		if suggest.Text == text {
			return
		}
	}
	suggests = append(suggests, goterm.Suggest{Text: text, Description: description})
}

func GetSuggestList() []goterm.Suggest {
	return suggests
}

func completer(in goterm.Document) []goterm.Suggest {
	if in.TextBeforeCursor() == "" {
		return []goterm.Suggest{}
	}

	return goterm.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
}

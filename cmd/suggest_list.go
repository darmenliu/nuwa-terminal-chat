package main

import (
	"os"
	"path/filepath"
	"strings"

	goterm "github.com/c-bata/go-prompt"
)

var suggests = []goterm.Suggest{
	{Text: "chatmode", Description: "Set terminal as a pure chat robot mode"},
	{Text: "cmdmode", Description: "Set terminal as a command mode, use natural language to communicate"},
	{Text: "taskmode", Description: "Set terminal as a task mode, use natural language to communicate to execute tasks"},
	{Text: "agentmode", Description: "Set terminal as an agent mode, use agent to do some automation work"},
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
	suggest := []goterm.Suggest{}

	// 检查是否是文件路径补全
	if strings.HasPrefix(in.Text, "./") || strings.HasPrefix(in.Text, "/") {
		dir := filepath.Dir(in.Text)
		files, err := os.ReadDir(dir)
		if err != nil {
			return suggest
		}

		for _, file := range files {
			name := file.Name()
			fullPath := filepath.Join(dir, name)
			if file.IsDir() {
				name += "/"
				fullPath += "/"
			}
			suggest = append(suggest, goterm.Suggest{
				Text:        fullPath,
				Description: name,
			})
		}
		return suggest
	}
	return goterm.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
}

package main

import (
	"context"
	"fmt"
	"os"

	"nuwa-engineer/pkg/llms/gemini"

	goterm "github.com/c-bata/go-prompt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func GenerateContent(ctx context.Context, prompt string) (string, error) {
	ctx = context.Background()
	model, err := gemini.NewGemini(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to create model: %w", err)
	}
	defer model.CloseBackend()

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("Failed to generate content: %w", err)
	}
	return resp, nil
}

func FailureExit() {
	os.Exit(1)
}

func executor(in string) {
	fmt.Println("You: " + in)
	if in == "" {
		return
	}
	ctx := context.Background()
	rsp, err := GenerateContent(ctx, in)
	if err != nil {
		fmt.Println("NUWA: " + err.Error())
		return
	}
	fmt.Println("NUWA: " + rsp)
}

func completer(in goterm.Document) []goterm.Suggest {
	s := []goterm.Suggest{
		{Text: "chatmode", Description: "Set terminal as a pure chat robot mode"},
		{Text: "cmdmode", Description: "Set terminal as a command mode, use natural language to communicate"},
		{Text: "taskmode", Description: "Set terminal as a task mode, use natural language to communicate to execute tasks"},
		{Text: "exit", Description: "Exit the terminal"},
	}
	return goterm.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {

	// Initialize a big text display with the letters "Nuwa" and "Engineer"
	// "P" is displayed in cyan and "Term" is displayed in light magenta
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Nuwa", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle(" Terminal", pterm.FgLightMagenta.ToStyle())).
		Render() // Render the big text to the terminal

	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	// // Create nuwa-engineer workspace
	// workspaceManager := NewDefaultWorkSpaceManager()
	// if workspaceManager.IsWorkspaceEixst() {
	// 	logger.Info("workspace already exist")
	// } else {
	// 	err := workspaceManager.CreateWorkspace()
	// 	if err != nil {
	// 		logger.Error("failed to create workspace,", logger.Args("err", err.Error()))
	// 		FailureExit()
	// 	}
	// }

	p := goterm.New(
		executor,
		completer,
		goterm.OptionPrefix(">>> "),
		goterm.OptionTitle("NUWA TERMINAL"),
	)
	p.Run()

}

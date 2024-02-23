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

const (
	ChatMode = "chatmode"
	CmdMode  = "cmdmode"
	TaskMode = "taskmode"
	Exit = "exit"
)

var CurrentMode string = ChatMode

func SetCurrentMode(in string) {
	CurrentMode = in
}

func GetSysPromptAccordingMode(current string) string {
	switch current {
	case ChatMode:
		return ""
	case CmdMode:
		return ""
	case TaskMode:
		return ""
	default:
		return ""
	}
}

func GetPromptAccordingToCurrentMode(current string, in string) string {
	sysPrompt := GetSysPromptAccordingMode(current)
	return sysPrompt + in
}

func GenerateContent(ctx context.Context, prompt string) (string, error) {
	model, err := gemini.NewGemini(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create model: %w", err)
	}
	defer model.CloseBackend()

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return resp, nil
}

func FailureExit() {
	os.Exit(1)
}

func executor(in string) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	fmt.Println("You: " + in)
	if in == "" {
		return
	}

	if in == Exit {
		logger.Info("NUWA TERMINAL: Goodbye!")
		os.Exit(0)
	}

	if (in == ChatMode) || (in == CmdMode) || (in == TaskMode) {
		SetCurrentMode(in)
		logger.Info("NUWA TERMINAL: Mode is " + CurrentMode)
		return
	}

	prompt := GetPromptAccordingToCurrentMode(CurrentMode, in)

	ctx := context.Background()
	rsp, err := GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"nuwa-engineer/pkg/cmdexe"
	"nuwa-engineer/pkg/llms"
	"nuwa-engineer/pkg/llms/gemini"
	"nuwa-engineer/pkg/llms/groq"
	"nuwa-engineer/pkg/llms/ollama"
	"nuwa-engineer/pkg/parser"
	"nuwa-engineer/pkg/prompts"

	goterm "github.com/c-bata/go-prompt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	ChatMode = "chatmode"
	CmdMode  = "cmdmode"
	TaskMode = "taskmode"
	Exit     = "exit"

	Catchdir   = ".nuwa-terminal"
	ScriptsDir = "scripts"
)

var CurrentMode string = ChatMode

func SetCurrentMode(in string) {
	CurrentMode = in
}

func GetSysPromptAccordingMode(current string) string {
	switch current {
	case ChatMode:
		return prompts.GetChatModePrompt()
	case CmdMode:
		return prompts.GetCmdModePrompt()
	case TaskMode:
		return prompts.GetTaskModePrompt()
	default:
		return ""
	}
}

func GetPromptAccordingToCurrentMode(current string, in string) string {
	sysPrompt := GetSysPromptAccordingMode(current)
	return sysPrompt + "\n" + in
}

func GenerateContent(ctx context.Context, prompt string) (string, error) {

	llmBackend := os.Getenv("LLM_BACKEND")
	modelName := os.Getenv("LLM_MODEL_NAME")
	modeTemperature, err := strconv.ParseFloat(os.Getenv("LLM_TEMPERATURE"), 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse LLM_TEMPERATURE: %w", err)
	}
	if llmBackend == "" {
		llmBackend = "gemini"
		modelName = "gemini-1.5-pro"
	}
	var model llms.Model
	switch llmBackend {
	case "gemini":
		model, err = gemini.NewGemini(ctx, modelName)
		if err != nil {
			return "", fmt.Errorf("failed to create model: %w", err)
		}

	case "ollama":
		model, err = ollama.NewOllama(ctx, modelName, modeTemperature)
		if err != nil {
			return "", fmt.Errorf("failed to create model: %w", err)
		}
	case "groq":
		model, err = groq.NewGroq(ctx, modelName, modeTemperature)
		if err != nil {
			return "", fmt.Errorf("failed to create model: %w", err)
		}
	default:
		return "", fmt.Errorf("unknown LLM backend: %s", llmBackend)
	}
	defer model.CloseBackend()

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return resp, nil
}

func Chat(ctx context.Context, message string, sysPrompt string) (string, error) {
	llmBackend := os.Getenv("LLM_BACKEND")
	modelName := os.Getenv("LLM_MODEL_NAME")
	modeTemperature, err := strconv.ParseFloat(os.Getenv("LLM_TEMPERATURE"), 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse LLM_TEMPERATURE: %w", err)
	}
	if llmBackend == "" {
		llmBackend = "gemini"
		modelName = "gemini-1.5-pro"
	}
	var model llms.Model
	switch llmBackend {
	case "gemini":
		model, err = gemini.NewGemini(ctx, modelName)
		if err != nil {
			return "", fmt.Errorf("failed to create model: %w", err)
		}
	case "ollama":
		model, err = ollama.NewOllama(ctx, modelName, modeTemperature)
		if err != nil {
			return "", fmt.Errorf("failed to create model: %w", err)
		}
	case "groq":
		model, err = groq.NewGroq(ctx, modelName, modeTemperature)
		if err != nil {
			return "", fmt.Errorf("failed to create model: %w", err)
		}
	default:
		return "", fmt.Errorf("unknown LLM backend: %s", llmBackend)
	}
	defer model.CloseBackend()
	model.SetSystemPrompt(ctx, sysPrompt)
	resp, err := model.Chat(ctx, message)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return resp, nil
}

func FailureExit() {
	os.Exit(1)
}

func ParseScript(response string) (filename, content string, err error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	parser := parser.NewGoCodeParser()
	sources, err := parser.ParseCode(response)
	if err != nil {
		logger.Error("Failed to parse code from LLM response, error:", logger.Args("err", err.Error()))
		return "", "", err
	}

	if len(sources) == 0 {
		logger.Error("No source files found in LLM response")
		return "", "", fmt.Errorf("no source files found")
	}

	sources[0].ParseFileName()
	sources[0].ParseFileContent()

	filename = sources[0].FileName
	content = sources[0].FileContent
	return filename, content, nil
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

	// Add Suggest
	AddSuggest(in, "")

	ctx := context.Background()
	if CurrentMode == ChatMode {
		sysPrompt := GetSysPromptAccordingMode(CurrentMode)

		rsp, err := Chat(ctx, in, sysPrompt)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
			return
		}
		fmt.Println("NUWA: " + rsp)
		return
	}

	rsp, err := GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return
	}
	fmt.Println("NUWA: " + rsp)

	if CurrentMode == CmdMode {
		cmd, err := parser.ParseCmdFromString(rsp)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to parse command,", logger.Args("err", err.Error()))
			return
		}

		if cmd == "" {
			logger.Info("NUWA TERMINAL: empty command")
			return
		}

		output, err := cmdexe.ExecCommandWithOutput(cmd)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to execute command,", logger.Args("err", err.Error(), "output", output))
			return
		}
		fmt.Println(output)
	} else if CurrentMode == TaskMode {
		fmt.Println(rsp)
		filename, content, err := ParseScript(rsp)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to parse script,", logger.Args("err", err.Error()))
			return
		}

		if filename == "" {
			logger.Info("NUWA TERMINAL: empty script")
			return
		}

		homedir := os.Getenv("HOME")
		scriptdir := filepath.Join(homedir, Catchdir, ScriptsDir)
		err = os.MkdirAll(scriptdir, os.ModePerm)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to create script directory,", logger.Args("err", err.Error()))
			return
		}

		scriptfile := filepath.Join(scriptdir, filename)
		err = os.WriteFile(scriptfile, []byte(content), os.ModePerm)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to write script file,", logger.Args("err", err.Error()))
			return
		}

		logger.Info("NUWA TERMINAL: script file saved to " + scriptfile)

		output, err := cmdexe.ExecScriptWithOutput(scriptfile)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to execute script,", logger.Args("err", err.Error()))
			return
		}

		fmt.Println(output)
		// remove the script
		err = os.Remove(scriptfile)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to remove script file,", logger.Args("err", err.Error()))
			return
		}

		logger.Info("NUWA TERMINAL: script file removed")
	}

}

func main() {

	// Initialize a big text display with the letters "Nuwa" and "Terminal"
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
	//      logger.Info("workspace already exist")
	// } else {
	// 	err := workspaceManager.CreateWorkspace()
	// 	if err != nil {
	// 		logger.Error("failed to create workspace,", logger.Args("err", err.Error()))
	// 		FailureExit()
	// 	}
	// }
	defer fmt.Println("Bye!")
	p := goterm.New(
		executor,
		completer,
		goterm.OptionPrefix(">>> "),
		goterm.OptionTitle("NUWA TERMINAL"),
	)
	p.Run()

}

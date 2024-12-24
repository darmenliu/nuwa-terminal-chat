package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/nuwa"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"

	goterm "github.com/c-bata/go-prompt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	ChatMode  = "chatmode"
	CmdMode   = "cmdmode"
	TaskMode  = "taskmode"
	AgentMode = "agentmode"
	Exit      = "exit"
	BashMode  = "bashmode"

	Catchdir   = ".nuwa-terminal"
	ScriptsDir = "scripts"

	// 快捷键常量
	ChatModeKey  = "@" // Ctrl+C
	CmdModeKey   = "#" // Ctrl+F
	TaskModeKey  = "$" // Ctrl+S
	AgentModeKey = "&" // Ctrl+A
	BashModeKey  = ">" // Ctrl+B
)

var CurrentMode string = ChatMode
var CurrentDir string = ""

var modeManager nuwa.NuwaModeManager = nil

// GetSysPromptAccordingMode returns the system prompt according to the current mode.
// It takes the current mode as input and returns the corresponding prompt string.
func GetSysPromptAccordingMode(current string) string {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	switch current {
	case ChatMode:
		return prompts.GetChatModePrompt()
	case CmdMode:
		return prompts.GetCmdModePrompt()
	case TaskMode:
		prompt, err := prompts.GetTaskModePrompt()
		if err != nil {
			logger.Error("Failed to get task mode prompt:", logger.Args("err", err.Error()))
			return ""
		}
		return prompt
	case AgentMode:
		return ""
	default:
		return ""
	}
}

func GetPromptAccordingToCurrentMode(current string, in string) string {
	sysPrompt := GetSysPromptAccordingMode(current)
	return sysPrompt + "\n" + in
}

func FailureExit() {
	os.Exit(1)
}

func handleScriptMode(ctx context.Context, in string) (bool, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	// check if it is a .nw script file
	if strings.HasSuffix(in, ".nw") {
		// handle "nw path/to/script.nw" format
		scriptPath := in
		if strings.HasPrefix(in, "nw ") {
			scriptPath = strings.TrimPrefix(in, "nw ")
			scriptPath = strings.TrimSpace(scriptPath)
		}

		// if it is a relative path, convert it to an absolute path
		if !filepath.IsAbs(scriptPath) {
			curDir, err := os.Getwd()
			if err != nil {
				logger.Error("NUWA TERMINAL: failed to get current directory,", logger.Args("err", err.Error()))
				return true, err
			}
			scriptPath = filepath.Join(curDir, scriptPath)
		}

		err := handleNuwaScript(ctx, scriptPath)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to execute script,", logger.Args("err", err.Error()))
			return true, err
		}
		return true, nil
	}
	return false, nil
}

// handleChatMode 处理聊天模式
func handleChatMode(ctx context.Context, input string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	sysPrompt := GetSysPromptAccordingMode(ChatMode)
	nuwa, err := nuwa.NewNuwaChat(ctx, sysPrompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create NuwaChat,", logger.Args("err", err.Error()))
		return err
	}
	return nuwa.Run(input)
}

// handleCmdMode 处理命令模式
func handleCmdMode(ctx context.Context, prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	rsp, err := llms.GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}
	fmt.Println("NUWA: " + rsp)

	cmd, err := parser.ParseCmdFromString(rsp)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to parse command,", logger.Args("err", err.Error()))
		return err
	}

	if cmd == "" {
		logger.Info("NUWA TERMINAL: empty command")
		return nil
	}

	output, err := cmdexe.ExecCommandWithOutput(cmd)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute command,", logger.Args("err", err.Error(), "output", output))
		return err
	}
	fmt.Println(output)

	return nil
}

// handleNuwaScript execute nuwa script according to the filepath
func handleNuwaScript(ctx context.Context, filepath string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	nuwa, err := nuwa.NewNuwaScript(ctx, filepath)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create NuwaScript,", logger.Args("err", err.Error()))
		return err
	}
	return nuwa.Run(filepath)
}

// handleTaskMode execute task according to the prompt
func handleTaskMode(ctx context.Context, prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	nuwa, err := nuwa.NewNuwaTask(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create NuwaTask,", logger.Args("err", err.Error()))
		return err
	}
	return nuwa.Run(prompt)
}

// handleAgentMode execute agent according to the input
func handleAgentMode(ctx context.Context, input string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	nuwa, err := nuwa.NewNuwaAgent(ctx, input)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create NuwaAgent,", logger.Args("err", err.Error()))
		return err
	}
	return nuwa.Run(input)
}

// handleBashMode execute bash command
func handleBashMode(input string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	// 直接执行命令，不经过 LLM
	output, err := cmdexe.ExecCommandWithOutput(input)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute command,", logger.Args("err", err.Error(), "output", output))
		return err
	}
	fmt.Println(output)
	return nil
}

// executor 主执行函数
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

	// 处理模式切换
	if (in == ChatMode) || (in == CmdMode) || (in == TaskMode) || (in == AgentMode) || (in == BashMode) {
		modeManager.SwitchMode(in)
		return
	}

	ctx := context.Background()
	if isScript, err := handleScriptMode(ctx, in); err != nil {
		logger.Error("NUWA TERMINAL: failed to handle script mode,", logger.Args("err", err.Error()))
		return
	} else if isScript {
		return
	}

	prompt := GetPromptAccordingToCurrentMode(CurrentMode, in)
	AddSuggest(in, "")

	// 根据当前模式处理输入
	var err error
	switch CurrentMode {
	case ChatMode:
		err = handleChatMode(ctx, in)
	case CmdMode:
		err = handleCmdMode(ctx, prompt)
		modeManager.CheckDirChanged()
	case TaskMode:
		err = handleTaskMode(ctx, prompt)
	case AgentMode:
		err = handleAgentMode(ctx, in)
	case BashMode:
		err = handleBashMode(in)
		modeManager.CheckDirChanged()
	}

	if err != nil {
		logger.Error("NUWA TERMINAL: Error executing command", logger.Args("mode", CurrentMode, "error", err.Error()))
	}
}

func main() {
	flags := ParseCmdParams()
	PrintHelp(flags)

	// 初始化大文本显示
	err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Nuwa", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle(" Terminal", pterm.FgLightMagenta.ToStyle())).
		Render()
	if err != nil {
		pterm.Error.Printf("Can not render the big text to the terminal: %v\n", err)
		os.Exit(1)
	}

	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	modeManager = nuwa.NewNuwaModeManager()
	// Get current directory path
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatal("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
	}
	// Set initial mode
	if flags.chatMode {
		modeManager.SetCurrentMode(ChatMode)
	} else if flags.cmdMode {
		modeManager.SetCurrentMode(CmdMode)
	} else if flags.taskMode {
		modeManager.SetCurrentMode(TaskMode)
	} else if flags.agentMode {
		modeManager.SetCurrentMode(AgentMode)
	}

	// If there is a query, process it directly and exit
	if flags.query != "" {
		executor(flags.query)
		return
	}

	// If it is interactive mode or no other mode is specified, enter interactive mode
	if flags.interactive || (!flags.chatMode && !flags.cmdMode && !flags.taskMode && !flags.agentMode && flags.query == "") {
		defer fmt.Println("Bye!")
		// Set initial mode
		modeManager.SetCurrentMode(ChatMode)
		modeManager.SetCurrentDir(currentDir)

		p := goterm.New(
			executor,
			completer,
			goterm.OptionPrefix(""),
			goterm.OptionLivePrefix(modeManager.GetLivePrefix),
			goterm.OptionTitle("NUWA TERMINAL"),
			goterm.OptionAddKeyBind(
				goterm.KeyBind{Key: goterm.ControlC, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(ChatMode) }},
				goterm.KeyBind{Key: goterm.ControlF, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(CmdMode) }},
				goterm.KeyBind{Key: goterm.ControlS, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(TaskMode) }},
				goterm.KeyBind{Key: goterm.ControlA, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(AgentMode) }},
				goterm.KeyBind{Key: goterm.ControlB, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(BashMode) }},
			),
		)
		p.Run()
	}
}

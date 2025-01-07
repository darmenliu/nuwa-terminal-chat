package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/nuwa"

	goterm "github.com/c-bata/go-prompt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	Exit      = "exit"
	Catchdir   = ".nuwa-terminal"
	ScriptsDir = "scripts"
)

var modeManager nuwa.NuwaModeManager = nil



func GetPromptAccordingToCurrentMode(in string) string {
	sysPrompt := modeManager.GetSysPrompt()
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
	sysPrompt := modeManager.GetSysPrompt()
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
	nuwa, err := nuwa.NewNuwaCmd(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create NuwaCmd,", logger.Args("err", err.Error()))
		return err
	}
	return nuwa.Run(prompt)
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


func newBashSession() {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	cmd := exec.Command("bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.Info("NUWA TERMINAL: starting new bash session")
	err := cmd.Run()
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to start new bash session,", logger.Args("err", err.Error()))
	}
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

	if in == nuwa.BashMode {
		newBashSession()
		return
	}

	// 处理模式切换
	if (in == nuwa.ChatMode) || (in == nuwa.CmdMode) || (in == nuwa.TaskMode) || (in == nuwa.AgentMode) {
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

	prompt := GetPromptAccordingToCurrentMode(in)
	AddSuggest(in, "")

	// 根据当前模式处理输入
	var err error
	switch modeManager.GetCurrentMode() {
	case nuwa.ChatMode:
		err = handleChatMode(ctx, in)
	case nuwa.CmdMode:
		err = handleCmdMode(ctx, prompt)
		modeManager.CheckDirChanged()
	case nuwa.TaskMode:
		err = handleTaskMode(ctx, prompt)
	case nuwa.AgentMode:
		err = handleAgentMode(ctx, in)
	}

	if err != nil {
		logger.Error("NUWA TERMINAL: Error executing command", logger.Args("mode", modeManager.GetCurrentMode(), "error", err.Error()))
	}
}

func setCurrentWorkMode(flags *CommandFlags) {
	workMode := nuwa.ChatMode
	if flags.chatMode {
		workMode = nuwa.ChatMode
	} else if flags.cmdMode {
		workMode = nuwa.CmdMode
	} else if flags.taskMode {
		workMode = nuwa.TaskMode
	} else if flags.agentMode {
		workMode = nuwa.AgentMode
	}
	modeManager.SetCurrentMode(workMode)
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
	setCurrentWorkMode(flags)

	// If there is a query, process it directly and exit
	if flags.query != "" {
		executor(flags.query)
		return
	}

	// If it is interactive mode or no other mode is specified, enter interactive mode
	if flags.interactive || (!flags.chatMode && !flags.cmdMode && !flags.taskMode && !flags.agentMode && flags.query == "") {
		defer fmt.Println("Bye!")
		// Set initial mode
		modeManager.SetCurrentMode(nuwa.ChatMode)
		modeManager.SetCurrentDir(currentDir)

		p := goterm.New(
			executor,
			completer,
			goterm.OptionPrefix(""),
			goterm.OptionLivePrefix(modeManager.GetLivePrefix),
			goterm.OptionTitle("NUWA TERMINAL"),
			goterm.OptionAddKeyBind(
				goterm.KeyBind{Key: goterm.ControlC, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(nuwa.ChatMode) }},
				goterm.KeyBind{Key: goterm.ControlF, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(nuwa.CmdMode) }},
				goterm.KeyBind{Key: goterm.ControlS, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(nuwa.TaskMode) }},
				goterm.KeyBind{Key: goterm.ControlA, Fn: func(b *goterm.Buffer) { modeManager.SwitchMode(nuwa.AgentMode) }},
			),
		)
		p.Run()
	}
}

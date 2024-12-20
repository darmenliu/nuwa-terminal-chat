package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/agents"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"

	goterm "github.com/c-bata/go-prompt"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	lcagents "github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	lcllms "github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
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

	ChatModePrefix  = "@"
	CmdModePrefix   = "#"
	TaskModePrefix  = "$"
	AgentModePrefix = "&"
	BashModePrefix  = ">"

	// 快捷键常量
	ChatModeKey  = "@" // Ctrl+C
	CmdModeKey   = "#" // Ctrl+F
	TaskModeKey  = "$" // Ctrl+S
	AgentModeKey = "&" // Ctrl+A
	BashModeKey  = ">" // Ctrl+B
)

var CurrentMode string = ChatMode
var CurrentDir string = ""

// Set current directory
func SetCurrentDir(in string) {
	CurrentDir = in
	LivePrefixState.LivePrefix = CurrentDir + getModePrefix(CurrentMode) + " "
	LivePrefixState.IsEnable = true
}

func CheckDirChanged(in string) bool {
	if in == CurrentDir {
		return false
	}
	SetCurrentDir(in)
	return true
}

// SetCurrentMode sets the current mode to the specified value.
func SetCurrentMode(in string) {
	CurrentMode = in
	LivePrefixState.LivePrefix = CurrentDir + getModePrefix(CurrentMode) + " "
	LivePrefixState.IsEnable = true
}

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

type NuwaChat struct {
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	SystemPrompt string
}

func NewNuwaChat(ctx context.Context, systemPrompt string) (*NuwaChat, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	content := []lcllms.MessageContent{
		lcllms.TextParts(lcllms.ChatMessageTypeSystem, systemPrompt),
	}

	return &NuwaChat{
		model:        model,
		chatHistory:  content,
		SystemPrompt: systemPrompt,
	}, nil
}

func (n *NuwaChat) Chat(ctx context.Context, message string) (string, error) {
	n.chatHistory = append(n.chatHistory, lcllms.TextParts(lcllms.ChatMessageTypeHuman, message))

	resp, err := n.model.GenerateContent(ctx, n.chatHistory)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}

func FailureExit() {
	os.Exit(1)
}

// ParseScript parses the code from the LLM response and returns the filename and content of the first source file.
// If there are no source files found or if there is an error parsing the code, an error is returned.
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

	sources[0].FileName = uuid.New().String() + ".sh"
	sources[0].ParseFileContent()

	filename = sources[0].FileName
	content = sources[0].FileContent
	return filename, content, nil
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

// handleModeSwitch 处理模式切换
func handleModeSwitch(mode string) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		curDir = CurrentDir
	}

	SetCurrentMode(mode)
	SetCurrentDir(curDir)

	logger.Info("NUWA TERMINAL: Mode is " + CurrentMode)
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
	nuwa, err := NewNuwaChat(ctx, sysPrompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create NuwaChat,", logger.Args("err", err.Error()))
		return err
	}

	rsp, err := nuwa.Chat(ctx, input)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}

	fmt.Println("NUWA: " + rsp)
	return nil
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

	// 检查当前目录是否改变
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		return err
	}

	if CheckDirChanged(curDir) {
		LivePrefixState.LivePrefix = CurrentDir + getModePrefix(CurrentMode) + " "
		LivePrefixState.IsEnable = true
	}

	return nil
}

// handleNuwaScript 处理 .nw 脚本文件
func handleNuwaScript(ctx context.Context, filepath string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	// check if the file exists and has execute permission
	info, err := os.Stat(filepath)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to access script file,", logger.Args("err", err.Error()))
		return err
	}

	if info.Mode()&0111 == 0 {
		return fmt.Errorf("script file %s is not executable", filepath)
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to read script file,", logger.Args("err", err.Error()))
		return err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "#!/bin/nuwa") {
		return fmt.Errorf("invalid nuwa script file: must start with #!/bin/nuwa")
	}

	scriptPrompt := strings.Join(lines[0:], "\n")

	prompt, err := prompts.GetScriptModePrompt()
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to get script mode prompt,", logger.Args("err", err.Error()))
		return err
	}
	prompt = prompt + "\n" + scriptPrompt
	// generate content
	rsp, err := llms.GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}
	fmt.Println("NUWA: " + rsp)

	if err := parseScriptAndExecute(rsp); err != nil {
		logger.Error("NUWA TERMINAL: failed to parse script and execute,", logger.Args("err", err.Error()))
		return err
	}

	return nil
}

// handleTaskMode 处理任务模式
func handleTaskMode(ctx context.Context, prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	rsp, err := llms.GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}
	fmt.Println("NUWA: " + rsp)

	if err := parseScriptAndExecute(rsp); err != nil {
		logger.Error("NUWA TERMINAL: failed to parse script and execute,", logger.Args("err", err.Error()))
		return err
	}

	return nil
}

func parseScriptAndExecute(rsp string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	filename, content, err := ParseScript(rsp)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to parse script,", logger.Args("err", err.Error()))
		return err
	}

	if filename == "" {
		logger.Info("NUWA TERMINAL: empty script")
		return nil
	}

	scriptfile, err := prepareScriptFile(filename, content)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to prepare script file,", logger.Args("err", err.Error()))
		return err
	}

	output, err := cmdexe.ExecScriptWithOutput(scriptfile)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute script,", logger.Args("err", err.Error()), logger.Args("output", output))
		return err
	}

	logger.Info("NUWA TERMINAL: script output", logger.Args("output", output))

	if err := os.Remove(scriptfile); err != nil {
		logger.Error("NUWA TERMINAL: failed to remove script file,", logger.Args("err", err.Error()))
		return err
	}
	logger.Info("NUWA TERMINAL: script file removed")
	return nil
}

// prepareScriptFile 准备脚本文件
func prepareScriptFile(filename, content string) (string, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	homedir := os.Getenv("HOME")
	scriptdir := filepath.Join(homedir, Catchdir, ScriptsDir)
	if err := os.MkdirAll(scriptdir, os.ModePerm); err != nil {
		logger.Error("NUWA TERMINAL: failed to create script directory,", logger.Args("err", err.Error()))
		return "", err
	}

	scriptfile := filepath.Join(scriptdir, filename)
	if err := os.WriteFile(scriptfile, []byte(content), os.ModePerm); err != nil {
		logger.Error("NUWA TERMINAL: failed to write script file,", logger.Args("err", err.Error()))
		return "", err
	}

	logger.Info("NUWA TERMINAL: script file saved to " + scriptfile)
	return scriptfile, nil
}

// handleAgentMode 处理代理模式
func handleAgentMode(ctx context.Context, input string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	llm, err := llms.GetLLMBackend(ctx)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to get LLM backend,", logger.Args("err", err.Error()))
		return err
	}

	agentTools := []tools.Tool{
		&agents.ScriptExecutor{},
	}

	agent := agents.NewTroubleshootingAgent(llm, agentTools, "output", nil)
	executor := lcagents.NewExecutor(agent)
	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to run agent,", logger.Args("err", err.Error()))
		return err
	}

	fmt.Println("NUWA: " + answer)
	return nil
}

// handleBashMode 处理 bash mode
func handleBashMode(input string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	// 直接执行命令，不经过 LLM
	output, err := cmdexe.ExecCommandWithOutput(input)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute command,", logger.Args("err", err.Error(), "output", output))
		return err
	}
	fmt.Println(output)

	// 检查当前目录是否改变
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		return err
	}

	if CheckDirChanged(curDir) {
		LivePrefixState.LivePrefix = CurrentDir + getModePrefix(CurrentMode) + " "
		LivePrefixState.IsEnable = true
	}

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
		handleModeSwitch(in)
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
	case TaskMode:
		err = handleTaskMode(ctx, prompt)
	case AgentMode:
		err = handleAgentMode(ctx, in)
	case BashMode:
		err = handleBashMode(in)
	}

	if err != nil {
		logger.Error("NUWA TERMINAL: Error executing command", logger.Args("mode", CurrentMode, "error", err.Error()))
	}
}

func getModePrefix(mode string) string {
	switch mode {
	case ChatMode:
		return ChatModePrefix
	case CmdMode:
		return CmdModePrefix
	case TaskMode:
		return TaskModePrefix
	case AgentMode:
		return AgentModePrefix
	case BashMode:
		return BashModePrefix
	default:
		return ChatModePrefix
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

	// Get current directory path
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatal("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
	}

	// Set initial mode
	if flags.chatMode {
		SetCurrentMode(ChatMode)
	} else if flags.cmdMode {
		SetCurrentMode(CmdMode)
	} else if flags.taskMode {
		SetCurrentMode(TaskMode)
	} else if flags.agentMode {
		SetCurrentMode(AgentMode)
	}

	// If there is a query, process it directly and exit
	if flags.query != "" {
		executor(flags.query)
		return
	}

	// If it is interactive mode or no other mode is specified, enter interactive mode
	if flags.interactive || (!flags.chatMode && !flags.cmdMode && !flags.taskMode && !flags.agentMode && flags.query == "") {
		defer fmt.Println("Bye!")

		// 设置初始 LivePrefix
		LivePrefixState.LivePrefix = currentDir + getModePrefix(CurrentMode) + " "
		LivePrefixState.IsEnable = true

		p := goterm.New(
			executor,
			completer,
			goterm.OptionPrefix(""),
			goterm.OptionLivePrefix(changeLivePrefix),
			goterm.OptionTitle("NUWA TERMINAL"),
			goterm.OptionAddKeyBind(
				goterm.KeyBind{Key: goterm.ControlC, Fn: func(b *goterm.Buffer) { handleModeSwitch(ChatMode) }},
				goterm.KeyBind{Key: goterm.ControlF, Fn: func(b *goterm.Buffer) { handleModeSwitch(CmdMode) }},
				goterm.KeyBind{Key: goterm.ControlS, Fn: func(b *goterm.Buffer) { handleModeSwitch(TaskMode) }},
				goterm.KeyBind{Key: goterm.ControlA, Fn: func(b *goterm.Buffer) { handleModeSwitch(AgentMode) }},
				goterm.KeyBind{Key: goterm.ControlB, Fn: func(b *goterm.Buffer) { handleModeSwitch(BashMode) }},
			),
		)
		p.Run()
	}
}

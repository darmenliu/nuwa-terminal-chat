package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/agents"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"

	goterm "github.com/c-bata/go-prompt"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	lcagents "github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	lcllms "github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	lcollama "github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

const (
	ChatMode  = "chatmode"
	CmdMode   = "cmdmode"
	TaskMode  = "taskmode"
	AgentMode = "agentmode"
	Exit      = "exit"

	Catchdir   = ".nuwa-terminal"
	ScriptsDir = "scripts"

	ChatModePrefix  = "@"
	CmdModePrefix   = "#"
	TaskModePrefix  = "$"
	AgentModePrefix = "&"

	// 快捷键常量
	ChatModeKey  = "@" // Ctrl+2
	CmdModeKey   = "#" // Ctrl+3
	TaskModeKey  = "$" // Ctrl+4
	AgentModeKey = "&" // Ctrl+7
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

func GetLLMBackend(ctx context.Context) (lcllms.Model, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	llmBackend := os.Getenv("LLM_BACKEND")
	modelName := os.Getenv("LLM_MODEL_NAME")
	apiKey := os.Getenv("LLM_API_KEY")

	if llmBackend == "" {
		llmBackend = "gemini"
		modelName = "gemini-1.5-pro"
	}

	if apiKey == "" {
		logger.Error("LLM_API_KEY is not set")
		return nil, errors.New("LLM_API_KEY is not set")
	}

	var model lcllms.Model
	var err error
	switch llmBackend {
	case "gemini":
		model, err = googleai.New(ctx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel(modelName))
	case "ollama":
		serverUrl := os.Getenv("OLLAMA_SERVER_URL")
		if serverUrl == "" {
			logger.Error("OLLAMA_SERVER_URL is not set")
			return nil, errors.New("OLLAMA_SERVER_URL is not set")
		}
		model, err = lcollama.New(lcollama.WithModel(modelName), lcollama.WithServerURL(serverUrl))
	case "groq":
		model, err = openai.New(
			openai.WithModel("llama3-8b-8192"),
			openai.WithBaseURL("https://api.groq.com/openai/v1"),
			openai.WithToken(apiKey),
		)
	case "deepseek":
		model, err = openai.New(
			openai.WithModel(modelName),
			openai.WithBaseURL("https://api.deepseek.com/beta"),
			openai.WithToken(apiKey),
		)
	case "claude":
		model, err = anthropic.New(
			anthropic.WithModel("claude-3-5-sonnet-20240620"),
			anthropic.WithToken(apiKey),
		)
	default:
		return nil, fmt.Errorf("unknown LLM backend: %s", llmBackend)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create %s client, error:", llmBackend), logger.Args("err", err.Error()))
		return nil, err
	}

	return model, nil
}

// GenerateContent generates content using the specified prompt.
// It takes a context.Context object and a prompt string as input.
// The function returns the generated content as a string and an error if any.
// The function first reads the values of environment variables LLM_BACKEND, LLM_MODEL_NAME, and LLM_TEMPERATURE.
// If LLM_BACKEND is empty, it defaults to "gemini" and LLM_MODEL_NAME defaults to "gemini-1.5-pro".
// The function then creates a model based on the value of LLM_BACKEND using the specified context, model name, and temperature.
// The model is closed using the CloseBackend method before returning the generated content.
// If an error occurs during model creation or content generation, an error is returned.
func GenerateContent(ctx context.Context, prompt string) (string, error) {

	model, err := GetLLMBackend(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get LLM backend: %w", err)
	}

	modeTemperature, err := strconv.ParseFloat(os.Getenv("LLM_TEMPERATURE"), 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse LLM_TEMPERATURE: %w", err)
	}

	resp, err := lcllms.GenerateFromSinglePrompt(ctx, model, prompt, lcllms.WithTemperature(modeTemperature))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return resp, nil
}

type NuwaChat struct {
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	SystemPrompt string
}

func NewNuwaChat(ctx context.Context, systemPrompt string) (*NuwaChat, error) {
	model, err := GetLLMBackend(ctx)
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

	rsp, err := GenerateContent(ctx, prompt)
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

// handleTaskMode 处理任务模式
func handleTaskMode(ctx context.Context, prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	rsp, err := GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}
	fmt.Println("NUWA: " + rsp)

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
		return err
	}

	output, err := cmdexe.ExecScriptWithOutput(scriptfile)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute script,", logger.Args("err", err.Error()), logger.Args("output", output))
		return err
	}

	fmt.Println(output)

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

	llm, err := GetLLMBackend(ctx)
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
	if (in == ChatMode) || (in == CmdMode) || (in == TaskMode) || (in == AgentMode) {
		handleModeSwitch(in)
		return
	}

	prompt := GetPromptAccordingToCurrentMode(CurrentMode, in)
	AddSuggest(in, "")
	ctx := context.Background()

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
	}

	if err != nil {
		logger.Error("NUWA TERMINAL: Error executing command", logger.Args("mode", CurrentMode, "error", err.Error()))
	}
}

// 添加新的命令行参数结构体
type CommandFlags struct {
	interactive bool
	chatMode    bool
	cmdMode     bool
	taskMode    bool
	agentMode   bool
	query       string
	help        bool
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
	default:
		return ChatModePrefix
	}
}

// 添加快捷键处理数
func handleKeyBinding(in goterm.Document) (goterm.Document, bool) {

	switch in.Text {
	case ChatModeKey: // Ctrl+2
		handleModeSwitch(ChatMode)
		return goterm.Document{}, true
	case CmdModeKey: // Ctrl+3
		handleModeSwitch(CmdMode)
		return goterm.Document{}, true
	case TaskModeKey: // Ctrl+4
		handleModeSwitch(TaskMode)
		return goterm.Document{}, true
	case AgentModeKey: // Ctrl+7
		handleModeSwitch(AgentMode)
		return goterm.Document{}, true
	}
	return in, false
}

func main() {
	// 定义命令行参数
	flags := CommandFlags{}
	flag.BoolVar(&flags.interactive, "i", false, "Interactive mode")
	flag.BoolVar(&flags.chatMode, "c", false, "Chat mode")
	flag.BoolVar(&flags.cmdMode, "m", false, "Command mode")
	flag.BoolVar(&flags.taskMode, "t", false, "Task mode")
	flag.BoolVar(&flags.agentMode, "a", false, "Agent mode")
	flag.StringVar(&flags.query, "q", "", "Query to process")
	flag.BoolVar(&flags.help, "h", false, "Show help message")
	flag.Parse()

	// 显示帮助信息
	if flags.help {
		fmt.Println("Nuwa Terminal - Your AI-powered terminal assistant")
		fmt.Println("\nUsage:")
		fmt.Println("  nuwa-terminal [flags] [query]")
		fmt.Println("\nFlags:")
		fmt.Println("  -i    Enter interactive mode, the nuwa will be like a bash environment，you can execute commands or tasks with natural language")
		fmt.Println("  -c    Chat mode, you can ask questions to Nuwa with natural language")
		fmt.Println("  -m    Command mode, you can execute commands with natural language")
		fmt.Println("  -t    Task mode, you can create a task with natural language，then nuwa will create a script to complete the task")
		fmt.Println("  -a    Agent mode, this is a experimental feature，you can ask Nuwa to help you execute more complex tasks, but the result may not be as expected")
		fmt.Println("  -q    User's input like a question, query or instruction")
		fmt.Println("  -h    Show this help message")
		fmt.Println("\nShortcuts (in interactive mode):")
		fmt.Println("  Ctrl+2    Switch to Chat mode")
		fmt.Println("  Ctrl+3    Switch to Command mode")
		fmt.Println("  Ctrl+4    Switch to Task mode")
		fmt.Println("  Ctrl+7    Switch to Agent mode")
		fmt.Println("\nExamples:")
		fmt.Println("  nuwa-terminal -c -q \"who are you?\"")
		fmt.Println("  nuwa-terminal -i")
		fmt.Println("  nuwa-terminal -m -q \"list all files\"")
		os.Exit(0)
	}

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
			),
		)
		p.Run()
	}
}

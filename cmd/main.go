package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/agents"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/system"

	goterm "github.com/c-bata/go-prompt"
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
)

var CurrentMode string = ChatMode
var CurrentDir string = ""

// Set current directory
func SetCurrentDir(in string) {
	CurrentDir = in
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

	sources[0].ParseFileName()
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

	if (in == ChatMode) || (in == CmdMode) || (in == TaskMode) || (in == AgentMode) {
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
		nuwa, err := NewNuwaChat(ctx, sysPrompt)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to create NuwaChat,", logger.Args("err", err.Error()))
			return
		}
		rsp, err := nuwa.Chat(ctx, in)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
			return
		}
		fmt.Println("NUWA: " + rsp)
		return
	}

	if CurrentMode == CmdMode {

		rsp, err := GenerateContent(ctx, prompt)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
			return
		}
		fmt.Println("NUWA: " + rsp)

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

		// Check if the current directory has changed
		curDir, err := os.Getwd()
		if err != nil {
			logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
			return
		}

		if CheckDirChanged(curDir) {
			LivePrefixState.LivePrefix = CurrentDir + ">>>"
			LivePrefixState.IsEnable = true
		}

	} else if CurrentMode == TaskMode {
		rsp, err := GenerateContent(ctx, prompt)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
			return
		}
		fmt.Println("NUWA: " + rsp)

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
	} else if CurrentMode == AgentMode {
		llm, err := GetLLMBackend(ctx)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to get LLM backend,", logger.Args("err", err.Error()))
			return
		}

		agentTools := []tools.Tool{
			&agents.ScriptExecutor{},
		}

		agent := agents.NewTroubleshootingAgent(llm, agentTools, "output", nil)
		executor := lcagents.NewExecutor(agent)
		answer, err := chains.Run(context.Background(), executor, in)
		if err != nil {
			logger.Error("NUWA TERMINAL: failed to run agent,", logger.Args("err", err.Error()))
			return
		}
		fmt.Println("NUWA: " + answer)
	}

}

func main() {
	// Initialize a big text display with the letters "Nuwa" and "Terminal"
	// "P" is displayed in cyan and "Term" is displayed in light magenta
	err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Nuwa", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle(" Terminal", pterm.FgLightMagenta.ToStyle())).
		Render() // Render the big text to the terminal
	if err != nil {
		pterm.Error.Printf("Can not render the big text to the terminal: %v\n", err)
		os.Exit(1)
	}

	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

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
	// 获取系统信息
	sysInfo := system.GetSystemInfo()

	// 将系统信息转换为格式化的 JSON 字符串
	jsonStr, err := sysInfo.ToJSON()
	if err != nil {
		log.Fatalf("转换系统信息为 JSON 时出错: %v", err)
	}

	// 打印 JSON 字符串
	fmt.Println("系统信息:")
	fmt.Println(jsonStr)

	//Get current directory path
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatal("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
	}

	defer fmt.Println("Bye!")
	p := goterm.New(
		executor,
		completer,
		goterm.OptionPrefix(currentDir+">>> "),
		goterm.OptionLivePrefix(changeLivePrefix),
		goterm.OptionTitle("NUWA TERMINAL"),
	)
	p.Run()

}

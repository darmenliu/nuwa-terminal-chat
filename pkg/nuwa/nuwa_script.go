package nuwa

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"
	"github.com/pterm/pterm"
	lcllms "github.com/tmc/langchaingo/llms"
)

type NuwaScript struct {
	ctx          context.Context
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	systemPrompt string
	currentDir   string
	catchdir     string
	scriptsdir   string
}

func NewNuwaScript(ctx context.Context, systemPrompt string) (*NuwaScript, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	return &NuwaScript{
		ctx:          ctx,
		model:        model,
		chatHistory:  []lcllms.MessageContent{},
		systemPrompt: systemPrompt,
		currentDir:   "",
		catchdir:     NuwaCatchDir,
		scriptsdir:   NuwaScriptsDir,
	}, nil
}

func (n *NuwaScript) Run(prompt string) error {
	return n.handleNuwaScript(n.ctx, prompt)
}

// handleNuwaScript handles the nuwa script mode
func (n *NuwaScript) handleNuwaScript(ctx context.Context, filepath string) error {
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

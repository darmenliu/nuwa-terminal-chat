package agents

import (
	"context"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/pterm/pterm"
	"github.com/tmc/langchaingo/tools"
)

type ScriptExecutor struct {

}

var _ tools.Tool = &ScriptExecutor{}

// Description returns a string describing the ScriptExecutor tool.
func (e *ScriptExecutor) Description() string {
	return `Useful for execute the shell scripts. 
	The input to this tool should be a shell script file path`
}

// Name returns the name of the tool.
func (e *ScriptExecutor) Name() string {
	return "ScriptExecutor"
}

func (e *ScriptExecutor) Call(ctx context.Context, input string) (string, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	codeParser := &ScriptCodeParser{}
	scriptfile, err := codeParser.ParseScriptAndSave(input)
	if err != nil {
		logger.Error("Failed to parse script from input, error:", logger.Args("err", err.Error()))
		return "", err
	}
	return cmdexe.ExecScriptWithOutput(scriptfile)
}

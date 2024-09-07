package agents

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"time"

	nuwaprmp "github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"
	"github.com/pterm/pterm"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// ToubleshootingAgent will implement langchaingo agent interface
// and will be used to troubleshoot the linux system problems form system logs
// or runtime logs.
type TroubleshootingAgent struct {
	// Chain is the chain used to call with the values. The chain should have an
	// input called "agent_scratchpad" for the agent to put its thoughts in.
	Chain chains.Chain
	// Tools is a list of the tools the agent can use.
	Tools []tools.Tool
	// Output key is the key where the final output is placed.
	OutputKey string
	// CallbacksHandler is the handler for callbacks.
	CallbacksHandler callbacks.Handler
}

const (
	_troubleshootingFinalAnswerAction = "Final Answer:"
)

func NewTroubleshootingAgent(llm llms.Model, tools []tools.Tool, outputkey string, callback callbacks.Handler) *TroubleshootingAgent {
	return &TroubleshootingAgent{
		Chain: chains.NewLLMChain(
			llm,
			CreateTroubleshootingAgentPrompt(tools),
			chains.WithCallback(callback),
		),
		Tools:            tools,
		OutputKey:        outputkey,
		CallbacksHandler: callback,
	}
}

func CreateTroubleshootingAgentPrompt(tools []tools.Tool) prompts.PromptTemplate {
	return prompts.PromptTemplate{
		Template:       nuwaprmp.SysPromptForAgentMode,
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad"},
		PartialVariables: map[string]any{
			"tools":        toolDescriptions(tools),
			"tool_names":   toolNames(tools),
			"ShellScriptFormat": nuwaprmp.ShellScriptFormat,
			"ShellExample": 	nuwaprmp.ShellExample,
			"current_time": time.Now().Format(time.RFC3339),
			"history":           "",
		},
	}
}

func toolNames(tools []tools.Tool) string {
	var tn strings.Builder
	for i, tool := range tools {
		if i > 0 {
			tn.WriteString(", ")
		}
		tn.WriteString(tool.Name())
	}

	return tn.String()
}


func toolDescriptions(tools []tools.Tool) string {
	var ts strings.Builder
	for _, tool := range tools {
		ts.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name(), tool.Description()))
	}

	return ts.String()
}

// Plan decides what action to take or returns the final result of the input.
func (tbs *TroubleshootingAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))
	for key, value := range inputs {
		fullInputs[key] = value
	}

	fullInputs["agent_scratchpad"] = constructScratchPad(intermediateSteps)

	var stream func(ctx context.Context, chunk []byte) error

	if tbs.CallbacksHandler != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			tbs.CallbacksHandler.HandleStreamingFunc(ctx, chunk)
			return nil
		}
	}

	output, err := chains.Predict(
		ctx,
		tbs.Chain,
		fullInputs,
		chains.WithStopWords([]string{"\nObservation:", "\n\tObservation:"}),
		chains.WithStreamingFunc(stream),
	)
	if err != nil {
		return nil, nil, err
	}

	return tbs.parseOutput(output)
}

func (tbs *TroubleshootingAgent) GetInputKeys() []string {
	chainInputs := tbs.Chain.GetInputKeys()

	// Remove inputs given in plan.
	agentInput := make([]string, 0, len(chainInputs))
	for _, v := range chainInputs {
		if v == "agent_scratchpad" {
			continue
		}
		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (tbs *TroubleshootingAgent) GetOutputKeys() []string {
	return []string{tbs.OutputKey}
}

func (tbs *TroubleshootingAgent) GetTools() []tools.Tool {
	return tbs.Tools
}

func constructScratchPad(steps []schema.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += step.Action.Log
			scratchPad += "\nObservation: " + step.Observation
		}
		scratchPad += "\n" + "Thought:"
	}

	return scratchPad
}

func (tbs *TroubleshootingAgent) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
    logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
    if strings.Contains(output, _troubleshootingFinalAnswerAction) {
        splits := strings.Split(output, _troubleshootingFinalAnswerAction)

        finishAction := &schema.AgentFinish{
            ReturnValues: map[string]any{
                tbs.OutputKey: splits[len(splits)-1],
            },
            Log: output,
        }

        return nil, finishAction, nil
    }

    // Normalize line endings to handle different platforms
    normalizedOutput := strings.ReplaceAll(output, "\r\n", "\n")
    
    // Print the normalized output for debugging
    logger.Info("Parsing output:", logger.Args("output", normalizedOutput))

    // Improved regex to handle dynamic script names and multiline content
    r := regexp.MustCompile(`(?s)Action: (.*?)\nAction_input:`)
    matches := r.FindStringSubmatch(normalizedOutput)
    if len(matches) == 0 {
        logger.Error("NUWA TERMINAL: Unable to parse the output,", logger.Args("output", normalizedOutput))
        return nil, nil, fmt.Errorf("%w: %s", agents.ErrUnableToParseOutput, normalizedOutput)
    }
	logger.Info("Matched:", logger.Args("match content for tool name:", matches[0]))
    return []schema.AgentAction{
        {Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(normalizedOutput), Log: normalizedOutput},
    }, nil, nil
}

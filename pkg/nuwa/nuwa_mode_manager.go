package nuwa

import (
	"os"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/prompts"
	"github.com/pterm/pterm"
)

const (
	ChatMode  = "chatmode"
	CmdMode   = "cmdmode"
	TaskMode  = "taskmode"
	AgentMode = "agentmode"
	BashMode  = "bashmode"
)

type NuwaModeManager interface {
	SwitchMode(mode string)
	GetCurrentMode() string
	SetCurrentMode(in string)
	SetCurrentDir(in string)
	GetModePrefix(mode string) string
	CheckDirChanged() bool
	GetSysPrompt() string
	GetLivePrefix() (string, bool)
}

type LivePrefix struct {
	Prefix   string
	IsEnable bool
}

type NuwaModeManagerImpl struct {
	currentMode string
	Prefix      LivePrefix
	currentDir  string
}

func NewNuwaModeManager() *NuwaModeManagerImpl {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		return nil
	}
	return &NuwaModeManagerImpl{
		currentMode: ChatMode,
		Prefix:      LivePrefix{Prefix: curDir + ChatModePrefix + " ", IsEnable: false},
		currentDir:  curDir,
	}
}

func (n *NuwaModeManagerImpl) SwitchMode(mode string) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		curDir = n.currentDir
	}

	n.SetCurrentMode(mode)
	n.SetCurrentDir(curDir)

	logger.Info("NUWA TERMINAL: Mode is " + n.currentMode)
}

func (n *NuwaModeManagerImpl) GetCurrentMode() string {
	return n.currentMode
}

// SetCurrentMode sets the current mode to the specified value.
func (n *NuwaModeManagerImpl) SetCurrentMode(in string) {
	n.currentMode = in
	n.Prefix.Prefix = n.currentDir + GetModePrefix(n.currentMode) + " "
	n.Prefix.IsEnable = true
}

func (n *NuwaModeManagerImpl) SetCurrentDir(in string) {
	n.currentDir = in
	n.Prefix.Prefix = n.currentDir + GetModePrefix(n.currentMode) + " "
	n.Prefix.IsEnable = true
}

func (n *NuwaModeManagerImpl) GetModePrefix(mode string) string {
	return GetModePrefix(mode)
}

func (n *NuwaModeManagerImpl) CheckDirChanged() bool {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		return false
	}
	if curDir == n.currentDir {
		return false
	}
	n.SetCurrentDir(curDir)
	return true
}

func (n *NuwaModeManagerImpl) GetSysPrompt() string {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	switch n.currentMode {
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

func (n *NuwaModeManagerImpl) GetLivePrefix() (string, bool) {
	return n.Prefix.Prefix, n.Prefix.IsEnable
}

func GetModePrefix(mode string) string {
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

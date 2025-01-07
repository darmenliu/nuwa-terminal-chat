package nuwa

type Nuwa interface {
	Run(prompt string) error
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

const (
	ChatModePrefix  = "@"
	CmdModePrefix   = "#"
	TaskModePrefix  = ">"
	AgentModePrefix = "&"

	NuwaCatchDir   = ".nuwa-terminal"
	NuwaScriptsDir = "scripts"
)

package main

import (
	"flag"
	"fmt"
	"os"
)

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

func ParseCmdParams() *CommandFlags {

	// 定义命令行参数
	flags := &CommandFlags{}
	flag.BoolVar(&flags.interactive, "i", false, "Interactive mode")
	flag.BoolVar(&flags.chatMode, "c", false, "Chat mode")
	flag.BoolVar(&flags.cmdMode, "m", false, "Command mode")
	flag.BoolVar(&flags.taskMode, "t", false, "Task mode")
	flag.BoolVar(&flags.agentMode, "a", false, "Agent mode")
	flag.StringVar(&flags.query, "q", "", "Query to process")
	flag.BoolVar(&flags.help, "h", false, "Show help message")
	flag.Parse()

	return flags
}

func PrintHelp(flags *CommandFlags) {
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
		fmt.Println("  Ctrl+B    Switch to Bash mode")
		fmt.Println("\nExamples:")
		fmt.Println("  nuwa-terminal -c -q \"who are you?\"")
		fmt.Println("  nuwa-terminal -i")
		fmt.Println("  nuwa-terminal -m -q \"list all files\"")
		os.Exit(0)
	}
}

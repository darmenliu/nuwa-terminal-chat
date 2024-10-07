package system

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

type SystemInfo struct {
	OS             string
	Arch           string
	CPU            int
	Memory         string
	AvailableTools []string
}

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		CPU:            runtime.NumCPU(),
		Memory:         fmt.Sprintf("%d MB", runtime.MemStats{}.Sys),
		AvailableTools: getAvailableTools(),
	}
}

func getAvailableTools() []string {
	directories := []string{"/usr/bin", "/usr/sbin", "/bin", "/sbin"}
	toolSet := make(map[string]bool)

	for _, dir := range directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if file.Mode()&0111 != 0 { // 检查是否有执行权限
				toolSet[file.Name()] = true
			}
		}
	}

	var availableTools []string
	for tool := range toolSet {
		availableTools = append(availableTools, tool)
	}

	return availableTools
}

func GetAvailableTools() string {
	return strings.Join(getAvailableTools(), ", ")
}

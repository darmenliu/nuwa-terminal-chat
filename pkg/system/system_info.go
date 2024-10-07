package system

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type SystemInfo struct {
	OS             OSInfo   `json:"os"`
	Arch           string   `json:"arch"`
	CPU            int      `json:"cpu"`
	Memory         string   `json:"memory"`
	AvailableTools []string `json:"available_tools"`
}

type OSInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	ID      string `json:"id"`
}

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:             getOSInfo(),
		Arch:           runtime.GOARCH,
		CPU:            runtime.NumCPU(),
		Memory:         fmt.Sprintf("%d MB", runtime.MemStats{}.Sys),
		AvailableTools: getAvailableTools(),
	}
}

func getOSInfo() OSInfo {
	osInfo := OSInfo{
		Name:    runtime.GOOS,
		Version: "unknown",
		ID:      "unknown",
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return osInfo
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

		switch key {
		case "NAME":
			osInfo.Name = value
		case "VERSION_ID":
			osInfo.Version = value
		case "ID":
			osInfo.ID = value
		}
	}

	return osInfo
}

func getAvailableTools() []string {
	directories := []string{"/usr/bin", "/usr/sbin", "/bin", "/sbin"}
	toolSet := make(map[string]bool)

	for _, dir := range directories {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.Mode()&0111 != 0 { // 检查是否有执行权限
				toolSet[entry.Name()] = true
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

// 新增函数：将 SystemInfo 转换为 JSON 字符串
func (si SystemInfo) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(si)
	if err != nil {
		return "", fmt.Errorf("转换 SystemInfo 为 JSON 时出错: %v", err)
	}
	return string(jsonBytes), nil
}

// 新增函数：将 SystemInfo 转换为格式化的 JSON 字符串
func (si SystemInfo) ToPrettyJSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(si, "", "  ")
	if err != nil {
		return "", fmt.Errorf("转换 SystemInfo 为格式化 JSON 时出错: %v", err)
	}
	return string(jsonBytes), nil
}

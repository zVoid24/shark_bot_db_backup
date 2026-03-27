package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func resolvePostgresTool(overridePath string, toolName string) (string, error) {
	if strings.TrimSpace(overridePath) != "" {
		return overridePath, nil
	}

	if path, err := exec.LookPath(toolName); err == nil {
		return path, nil
	}

	if runtime.GOOS == "windows" {
		patterns := []string{
			`C:\Program Files\PostgreSQL\*\bin\` + toolName + `.exe`,
			`C:\Program Files (x86)\PostgreSQL\*\bin\` + toolName + `.exe`,
		}

		for _, pattern := range patterns {
			matches, _ := filepath.Glob(pattern)
			if len(matches) == 0 {
				continue
			}
			sort.Strings(matches)
			return matches[len(matches)-1], nil
		}
	}

	return "", fmt.Errorf("%s executable not found. add it to PATH or set explicit path in config", toolName)
}

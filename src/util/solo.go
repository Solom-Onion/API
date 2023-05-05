package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseConfig() (map[string]string, error) {
	// Custom config file POGGERS

	configFile, err := os.Open("config.solo")

	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments or handle new lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line: %s", line)
		}
		config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return config, nil
}

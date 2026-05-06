package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const configFileName = "config.txt"

var config = func() map[string]string {
	dict := make(map[string]string)

	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatalf("Could not open config file: %s\n", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) <= 0 || line[0] == ';' {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			log.Fatalf("Error interpreting line: %s\n", line)
		}
		value, _, _ = strings.Cut(value, ";")

		dict[key] = strings.TrimSpace(value)
	}

	return dict
}()

func LoadString(key string) string {
	value, ok := config[key]
	if !ok {
		log.Fatalf("Config key not found: %s\n", key)
	}
	return value
}

func LoadInt(key string) int {
	value, err := strconv.Atoi(LoadString(key))
	if err != nil {
		log.Fatal(err)
	}
	return value
}

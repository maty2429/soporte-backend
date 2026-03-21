package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnvFiles(paths ...string) error {
	for _, path := range paths {
		if err := loadEnvFile(path); err != nil {
			return err
		}
	}

	return nil
}

func loadEnvFile(path string) (err error) {
	file, err := os.Open(path) //nolint:gosec // Path comes from trusted source (hardcoded in main)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open env file %s: %w", path, err)
	}
	defer func() {
		closeErr := file.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("close env file %s: %w", path, closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid env line %s:%d", path, lineNumber)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			return fmt.Errorf("invalid env line %s:%d", path, lineNumber)
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %s from %s:%d: %w", key, path, lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan env file %s: %w", path, err)
	}

	return nil
}

// Package themestate reads the shared desktop theme-state file.
package themestate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Path returns the path to the shared theme-state file.
func Path() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(stateHome, "theme-state.json")
}

// Detect returns the theme from the shared state file, defaulting to "dark".
func Detect() string {
	data, err := os.ReadFile(Path())
	if err != nil {
		return "dark"
	}

	var config struct {
		Theme string `json:"theme"`
	}
	if json.Unmarshal(data, &config) != nil || config.Theme == "" {
		return "dark"
	}
	return config.Theme
}

// ModTime returns the modification time of path, or an error if it cannot be statted.
func ModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

var (
	configPath string
	dataPath   string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get home directory")
	}

	switch runtime.GOOS {
	case "darwin":
		configPath = filepath.Join(homeDir, "Library", "Application Support", "lazysmtp")
		dataPath = filepath.Join(homeDir, "Library", "Application Support", "lazysmtp")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		configPath = filepath.Join(appData, "lazysmtp")
		dataPath = filepath.Join(appData, "lazysmtp")
	default: // Linux and other Unix-like systems
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			xdgConfig = filepath.Join(homeDir, ".config")
		}
		configPath = filepath.Join(xdgConfig, "lazysmtp")

		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			xdgData = filepath.Join(homeDir, ".local", "share")
		}
		dataPath = filepath.Join(xdgData, "lazysmtp")
	}

	if err := ensureDir(configPath); err != nil {
		panic(errors.New("failed to create config directory"))
	}
	if err := ensureDir(dataPath); err != nil {
		panic(errors.New("failed to create data directory"))
	}
}

func GetConfigPath() string {
	return configPath
}

func GetDataPath() string {
	return dataPath
}

func GetDefaultDBPath() string {
	return filepath.Join(dataPath, "lazysmtp.db")
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

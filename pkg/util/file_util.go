package util

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yanosea/cleancobra-pkg/errors"
)

func GetXDGDataHome() (string, error) {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to get home directory")
		}
		xdgDataHome = filepath.Join(homeDir, ".local", "share")
	}
	return xdgDataHome, nil
}

func MkdirIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return errors.Wrap(err, "failed to create directory")
		}
	}
	return nil
}

func InitializeJSONFile(filePath string, emptyData interface{}) error {
	file, err := json.MarshalIndent(emptyData, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal empty data")
	}

	if err := os.WriteFile(filePath, file, 0644); err != nil {
		return errors.Wrap(err, "failed to create initial file")
	}
	return nil
}

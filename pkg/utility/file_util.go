package utility

import (
	"path/filepath"

	"github.com/yanosea/cleancobra/pkg/errors"
	"github.com/yanosea/cleancobra/pkg/proxy"
)

type FileUtil interface {
	GetXDGDataHome() (string, error)
	MkdirIfNotExist(dirPath string) error
	InitializeJSONFile(filePath string, emptyData interface{}) error
}

type fileUtil struct {
	os   proxy.Os
	json proxy.Json
}

func NewFileUtil(
	os proxy.Os,
	json proxy.Json,
) FileUtil {
	return &fileUtil{
		os:   os,
		json: json,
	}
}

func (f *fileUtil) GetXDGDataHome() (string, error) {
	xdgDataHome := f.os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		homeDir, err := f.os.UserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to get home directory")
		}
		xdgDataHome = filepath.Join(homeDir, ".local", "share")
	}
	return xdgDataHome, nil
}

func (f *fileUtil) MkdirIfNotExist(dirPath string) error {
	if _, err := f.os.Stat(dirPath); f.os.IsNotExist(err) {
		if err := f.os.MkdirAll(dirPath, 0755); err != nil {
			return errors.Wrap(err, "failed to create directory")
		}
	}
	return nil
}

func (f *fileUtil) InitializeJSONFile(filePath string, emptyData interface{}) error {
	file, err := f.json.MarshalIndent(emptyData, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal empty data")
	}

	if err := f.os.WriteFile(filePath, file, 0644); err != nil {
		return errors.Wrap(err, "failed to create initial file")
	}
	return nil
}

package utils

import (
	"os"
	"path/filepath"
)

func GetBaseDir() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Dir(abs) + "/../../", nil
}

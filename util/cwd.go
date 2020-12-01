package util

import (
	"os"
	"path"
)

func GetCwd() string {
	execPath, _ := os.Executable()
	execDir, _ := path.Split(execPath)
	return execDir
}

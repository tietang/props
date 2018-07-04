package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func GetRunningDir() string {
	file, _ := exec.LookPath(os.Args[0])
	runningDir, _ := filepath.Abs(filepath.Dir(file))
	return runningDir
}
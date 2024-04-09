package runtimeutils

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// GetProjectRoot returns the root directory of the project.
// The project root is the directory that contains the go.mod file.
// If the project root is not found, it returns the current working directory.
//
// No parameters.
// Returns a string.
func GetProjectRoot() string {
	dir := getProjectPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return filepath.Dir(getProjectPathByCaller())
	} else if strings.Contains(dir, "/tmp/GoLand") {
		return filepath.Dir(getProjectPathByCaller())
	}
	return dir
}

// getTmpDir gets the system temporary directory path, compatible with 'go run'.
// Returns the path as a string.
func getTmpDir() string {
	dir := os.TempDir()

	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// getProjectPathByExecutable returns the absolute path of the current executing file.
func getProjectPathByExecutable() string {
	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// Resolve any symlinks and get the directory path
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// getProjectPathByCaller returns the absolute path of the current executing file (using go run).
func getProjectPathByCaller() string {
	var abPath string

	// Get the caller's file path and check if it was retrieved successfully
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}

	return abPath
}

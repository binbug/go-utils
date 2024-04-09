package runtimeutils

import "testing"

func TestGetProjectRoot(t *testing.T) {
	projectRoot := GetProjectRoot()
	t.Log(projectRoot)
}

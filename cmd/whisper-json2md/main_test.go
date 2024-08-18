package main

import (
	"os"
	"testing"
)

func TestTimestampOffset(t *testing.T) {
	tmpDir := os.TempDir()
	f, err := os.CreateTemp(tmpDir, "example.json")
	if err != nil {
		t.Fatalf("Failed to create temporary json file")
	}

	os.Args = []string{
		"whisper-json2md",
		f.Name(),
	}

	main()

}

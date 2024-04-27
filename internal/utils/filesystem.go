package utils

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/marcuswhybrow/ray-peat-rodeo/internal/global"
)

func Files[Result any](pwd, scope string, f func(filePath string) (*Result, error)) []Result {
	results := []Result{}

	entries, err := fs.ReadDir(os.DirFS(pwd), scope)
	if err != nil {
		log.Fatalf("Failed to read directory '%v': %v\n", path.Join(pwd, scope), err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		result, err := f(path.Join(scope, entry.Name()))
		if err != nil {
			log.Fatal("Failed to convert file entry to Result:", err)
		}

		if result != nil {
			results = append(results, *result)
		}
	}

	return results
}

// Convenience function to output HTML page
func MakePage(outPath string) (*os.File, string) {
	return MakeFile(path.Join(outPath, "index.html"))
}

var BuiltFiles []string
var builtFilesMutex sync.RWMutex

// Convenience function to output file
func MakeFile(outPath string) (*os.File, string) {
	buildPath := path.Join(global.BUILD_OUTPUT, outPath)
	parent := filepath.Dir(buildPath)
	illegalChars := "\":<>|*?\r\n"

	if strings.ContainsAny(outPath, illegalChars) {
		log.Fatalf("Failed to create build file '%v': "+
			"name contains an illegal character. "+
			"Illegal characters include: %v",
			outPath,
			illegalChars,
		)

	}

	err := os.MkdirAll(parent, 0755)
	if err != nil {
		log.Panicf("Failed to create directory '%v': %v", parent, err)
	}

	builtFilesMutex.Lock()
	if slices.Contains(BuiltFiles, buildPath) {
		log.Panicf(
			"Multiple writes attempted to the same build path: %v\n"+
				"  Common reasons for this include:\n"+
				"    - Two files in ./assets that have different dates in the filename, but the same wording after the date.\n"+
				"    - Two mentions that have the same name, but different capitalization.\n"+
				"    - A file in ./assets that has the same wording after the date as the wording of a mention.\n",
			buildPath,
		)
	}
	BuiltFiles = append(BuiltFiles, buildPath)
	builtFilesMutex.Unlock()

	f, err := os.Create(buildPath)
	if err != nil {
		log.Panicf("Failed to create file '%v': %v", buildPath, err)
	}

	return f, buildPath
}

package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

const BINARY_DIR = "./lib/bin"

var BIN_TARGZ_URLS = map[string]string{
	"modd":     "https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz",
	"devd":     "https://github.com/cortesi/devd/releases/download/v0.9/devd-0.9-linux64.tgz",
	"pagefind": "https://github.com/CloudCannon/pagefind/releases/download/v0.12.0/pagefind-v0.12.0-x86_64-unknown-linux-musl.tar.gz",
}

func DownloadBinaryIfAbsent(binaryName string) bool {
	fileName := path.Join(BINARY_DIR, binaryName)
	tarGzUrl := BIN_TARGZ_URLS[binaryName]
	if tarGzUrl == "" {
		panic(fmt.Sprintf("'%s' has no associated tar.gz download URL defined in BIN_TARGZ_URLS", binaryName))
	}
	if _, err := os.Stat(fileName); err == nil {
		return false
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Binary %s not installed. Downloading from %s\n", binaryName, tarGzUrl)
		resp := ReturnOrPanic(http.Get(tarGzUrl))
		defer resp.Body.Close()
		uncompressedStream := ReturnOrPanic(gzip.NewReader(resp.Body))
		tarReader := tar.NewReader(uncompressedStream)
		for true {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			PanicOnErr(err)
			if header.Typeflag == tar.TypeReg && path.Base(header.Name) == path.Base(fileName) {
				PanicOnErr(os.MkdirAll(path.Dir(fileName), 0755))
				outFile := ReturnOrPanic(os.Create(fileName))
				outFile.Chmod(755)
				ReturnOrPanic(io.Copy(outFile, tarReader))
				outFile.Close()
				fmt.Printf("Binary %s installed to %s\n", binaryName, fileName)
				return true
			}
		}
		panic(fmt.Sprintf("%s could not be found in %s", fileName, tarGzUrl))
	} else {
		panic(err)
	}
}

func DownloadBinariesIfAbsentAndExecuteLast(commandsWithArgs ...string) {
	{
		var wg sync.WaitGroup
		for _, commandWithArgs := range commandsWithArgs {
			wg.Add(1)
			command := strings.Split(commandWithArgs, " ")[0]
			go func() {
				defer wg.Done()
				DownloadBinaryIfAbsent(command)
			}()
		}
		wg.Wait()
	}
	lastCommandWithArgs := commandsWithArgs[len(commandsWithArgs)-1]
	lastCommand, lastCommandArgs := func() (string, []string) {
		lastCommandParts := strings.Split(lastCommandWithArgs, " ")
		return lastCommandParts[0], func() []string {
			if len(lastCommandParts) > 1 {
				return lastCommandParts[1:]
			}
			return []string{}
		}()
	}()
	modd := exec.Command(path.Join(BINARY_DIR, lastCommand), lastCommandArgs...)
	stdout := ReturnOrPanic(modd.StdoutPipe())
	var wg sync.WaitGroup
	wg.Add(1)
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Print(scanner.Text())
		}
		wg.Done()
	}()
	fmt.Printf("./%s\n", lastCommandWithArgs)
	PanicOnErr(modd.Start())
	wg.Wait()
	PanicOnErr(modd.Wait())
}

func RemoveBinaryDir() {
	PanicOnErr(os.RemoveAll(BINARY_DIR))
}

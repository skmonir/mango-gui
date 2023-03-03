package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func IsFileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsDirExist(folderPath string) bool {
	_, err := os.Stat(folderPath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(folderPath string) error {
	if !IsDirExist(folderPath) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GetFileNamesInDirectory(folderPath string) []string {
	var filenames []string
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println(err)
		return filenames
	}
	for _, e := range entries {
		filenames = append(filenames, e.Name())
	}
	return filenames
}

func GetFilenamesInDir(folderPath string) []string {
	var filenames []string

	if !IsDirExist(folderPath) {
		return filenames
	}

	files, err := ioutil.ReadDir(folderPath)

	if err == nil {
		for _, f := range files {
			fname := f.Name()
			if strings.HasSuffix(fname, ".cpp") {
				fname = strings.TrimSuffix(fname, filepath.Ext(fname))
				filenames = append(filenames, fname)
			}
		}
	}
	return filenames
}

func CreateFile(folderPath string, filename string) error {
	filePath := filepath.Join(folderPath, filename)
	if !IsFileExist(filePath) {
		fmt.Println("Creating file " + filePath)
		if err := CreateDir(folderPath); err != nil {
			return err
		}
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}

func OpenFile(filePath string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", filePath).Run()
	case "windows":
		exec.Command("cmd", filePath).Run()
	case "darwin":
		err = exec.Command("open", filePath).Run()
	default:
		err = errors.New("unsupported os")
	}
	return err
}

func ReadFileContent(filePath string, maxRow int, maxCol int) string {
	fmt.Println("Reading content from fileService " + filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	content := ""

	buf := make([]byte, maxCol+1)
	scanner.Buffer(buf, maxCol+1)

	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if rowCount == maxRow {
			content += "\n....."
			break
		}
		line = strings.TrimRight(line, " ")
		if len(line) > maxCol {
			line = line[:len(line)-1] + "..."
		}
		if len(content) > 0 {
			content += "\n"
		}
		content += line
		rowCount++
	}

	if err = scanner.Err(); err != nil {
		fmt.Println(err)
		return ""
	}
	return content
}

func WriteFileContent(folderPath string, filename string, data []byte) {
	if err := CreateFile(folderPath, filename); err != nil {
		fmt.Println(err)
		return
	}

	filePath := filepath.Join(folderPath, filename)
	fmt.Println("Writing data into " + filePath)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		fmt.Println(err)
		return
	}
}

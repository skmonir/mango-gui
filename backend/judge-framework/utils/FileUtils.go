package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
	if IsFileExist(folderPath) {
		return false
	}
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
		if !strings.Contains(strings.ToLower(e.Name()), ".ds_store") {
			filenames = append(filenames, e.Name())
		}
	}
	sort.Strings(filenames)
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

func ReadFileContent(filePath string, maxRow int, maxCol int) string {
	if !IsFileExist(filePath) {
		return ""
	}
	fmt.Println("Reading content from file " + filePath)
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return ResizeIOContentForUI(file, maxRow, maxCol)
}

func ResizeIOContentForUI(r io.Reader, maxRow int, maxCol int) string {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	content := ""
	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if rowCount == maxRow {
			content += "\n....."
			break
		}
		line = strings.TrimRight(line, " ")
		if len(line) > maxCol {
			line = line[:maxCol] + "..."
		}
		if len(content) > 0 {
			content += "\n"
		}
		content += line
		rowCount++
	}

	if err := scanner.Err(); err != nil {
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

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

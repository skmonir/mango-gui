package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func GetHtmlBody(URL string) (string, error) {
	log.Println("Fetching html from " + URL)
	resp, err := http.Get(URL)
	defer resp.Body.Close()
	log.Println("Fetched html with status ", resp.StatusCode)
	if err != nil || resp.StatusCode >= 400 {
		return "", errors.New("error while fetching web page")
	}
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}

func DecodeBase64(s string) string {
	decodedStr, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal("Error while decoding:", err)
		return ""
	}
	return strings.ToLower(string(decodedStr))
}

func PanicRecovery() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}

func splitUrlPath(path string) []string {
	var values []string
	for _, p := range strings.Split(strings.Trim(path, "/"), "/") {
		if len(strings.TrimSpace(p)) > 0 {
			values = append(values, p)
		}
	}
	return values
}

func ExtractInfoFromUrl(url string) (string, string, string) {
	platform, cid, pid := "", "", ""
	if strings.HasPrefix(url, "custom/") {
		values := splitUrlPath(url)
		platform = "custom"
		if len(values) > 1 {
			cid = values[1]
		}
		if len(values) > 2 {
			pid = values[2]
		}
	} else if strings.Contains(url, "atcoder.jp/contests") {
		platform = "atcoder"
		index := strings.Index(url, "atcoder.jp/contests")
		path := strings.Trim(url[index+len("atcoder.jp/contests"):], "/")
		values := splitUrlPath(path)
		if len(values) > 0 {
			cid = values[0]
		}
		if len(values) >= 3 {
			pid = values[2]
		}
	} else if strings.Contains(url, "codeforces.com") {
		platform = "codeforces"
		ctype := ""
		if strings.Contains(url, "codeforces.com/contest") {
			ctype = "contest"
		} else if strings.Contains(url, "codeforces.com/gym") {
			ctype = "gym"
		} else {
			return platform, cid, pid
		}
		index := strings.Index(url, "codeforces.com/"+ctype)
		path := strings.Trim(url[index+len("codeforces.com/"+ctype):], "/")
		values := splitUrlPath(path)
		if len(values) > 0 {
			cid = values[0]
		}
		if len(values) >= 3 {
			pid = values[2]
		}
	}
	fmt.Println("Extracted", url, ", got", platform, cid, pid)
	return platform, cid, pid
}

func ConvertMemoryInMb(memory uint64) uint64 {
	return memory / 1024
}

func ParseMemoryInMb(memory uint64) string {
	return fmt.Sprintf("%v MB", memory/1024/1024)
}

func ParseMemoryInKb(memory uint64) string {
	return fmt.Sprintf("%v KB", memory/1024)
}

func GetAppDataDirectoryPath() string {
	configPath := ""
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			configPath = os.Getenv("XDG_CONFIG_HOME")
		} else {
			configPath = filepath.Join(os.Getenv("HOME"), ".mango")
		}
	case "windows":
		configPath = filepath.Join(os.Getenv("APPDATA"), "mango")
	case "darwin":
		configPath = filepath.Join(os.Getenv("HOME"), ".mango")
	default:
		configPath = ""
	}

	return configPath
}

func OpenResourceInDefaultApplication(path string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/C", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, path)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		return err
	}
	return nil
}

func GetDefaultTemplateFilePathByLang(lang string) string {
	templateDirectory := filepath.Join(GetAppDataDirectoryPath(), "source_templates")
	templateFilePath := ""
	if lang == "cpp" {
		templateFilePath = filepath.Join(templateDirectory, "template_CPP.txt")
	} else if lang == "java" {
		templateFilePath = filepath.Join(templateDirectory, "template_Java.txt")
	} else if lang == "python" {
		templateFilePath = filepath.Join(templateDirectory, "template_Python.txt")
	}
	return templateFilePath
}

func GetBinaryFileExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

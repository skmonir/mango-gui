package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetHtmlBody(URL string) (string, error) {
	log.Println("Fetching html from " + URL)
	resp, err := http.Get(URL)
	defer resp.Body.Close()
	log.Println("Fetched html with status ", resp.StatusCode)
	if err != nil || resp.StatusCode >= 400 {
		return "", errors.New("Error while fetching web page")
	}
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}

func GetBody(client *http.Client, URL string) ([]byte, error) {
	log.Println("Fetching html from " + URL)
	resp, err := client.Get(URL)
	log.Println("Fetched html with status ", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func PostBody(client *http.Client, URL string, data url.Values) ([]byte, error) {
	log.Println("Posting data to ", URL)
	resp, err := client.PostForm(URL, data)
	log.Println("Posted data with status ", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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

func ExtractInfoFromUrl(url string) (string, string, string, string) {
	platform, cid, pid, ctype := "", "", "", ""
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
		if strings.Contains(url, "codeforces.com/contests") {
			ctype = "contests"
		} else if strings.Contains(url, "codeforces.com/contest") {
			ctype = "contest"
		} else if strings.Contains(url, "codeforces.com/gym") {
			ctype = "gym"
		} else {
			return platform, cid, pid, ctype
		}
		index := strings.Index(url, "codeforces.com/"+ctype)
		path := strings.Trim(url[index+len("codeforces.com/"+ctype):], "/ \n\t")
		values := splitUrlPath(path)
		if len(values) > 0 {
			cid = values[0]
		}
		if len(values) >= 3 {
			pid = values[2]
		}
	}
	fmt.Println("Extracted", url, ", got", platform, cid, pid)
	return platform, cid, pid, ctype
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

func GetAppHomeDirectoryPath() string {
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

func IsOsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsOsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsOsMac() bool {
	return runtime.GOOS == "darwin"
}

func GetDefaultTemplateFilePathByLang(lang string) string {
	templateDirectory := filepath.Join(GetAppHomeDirectoryPath(), "source_templates")
	templateFilePath := ""
	if lang == "cpp" {
		templateFilePath = filepath.Join(templateDirectory, "template_CPP.cpp")
	} else if lang == "java" {
		templateFilePath = filepath.Join(templateDirectory, "template_Java.java")
	} else if lang == "python" {
		templateFilePath = filepath.Join(templateDirectory, "template_Python.py")
	}
	return templateFilePath
}

func GetLangNameByFileExt(fileExt string) string {
	if fileExt == ".cpp" || fileExt == ".cc" {
		return "cpp"
	} else if fileExt == ".java" {
		return "java"
	} else if fileExt == ".py" {
		return "python"
	}
	return ""
}

func GetBinaryFileExt() string {
	if IsOsWindows() {
		return ".exe"
	}
	return ""
}

func ParseQueryMapFromUrl(URL string) url.Values {
	u, err := url.Parse(URL)
	if err != nil {
		return url.Values{}
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return url.Values{}
	}
	return q
}

func IsTimeInFuture(ctime time.Time) bool {
	today := time.Now()
	return today.Before(ctime)
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

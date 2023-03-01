package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func ExtractInfoFromUrl(url string) (string, string, string) {
	platform, cid, pid := "", "", ""
	if strings.Contains(url, "atcoder.jp") {
		platform = "atcoder"
		index := strings.Index(url, "atcoder.jp/contests")
		path := strings.Trim(url[index+len("atcoder.jp/contests"):], "/")
		var values []string
		for _, p := range strings.Split(path, "/") {
			if len(strings.TrimSpace(p)) > 0 {
				values = append(values, p)
			}
		}
		if len(values) > 0 {
			cid = values[0]
		}
		if len(values) >= 3 {
			pid = values[2]
		}
	} else if strings.Contains(url, "codeforces.com") {
		platform = "codeforces"
		index := strings.Index(url, "codeforces.com/contest")
		path := strings.Trim(url[index+len("codeforces.com/contest"):], "/")
		var values []string
		for _, p := range strings.Split(path, "/") {
			if len(strings.TrimSpace(p)) > 0 {
				values = append(values, p)
			}
		}
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

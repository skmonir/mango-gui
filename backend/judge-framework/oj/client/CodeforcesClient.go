package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/oj/cookiejar"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"github.com/skmonir/mango-gui/backend/socket"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type CodeforcesClient struct {
	Jar           *cookiejar.Jar `json:"cookies"`
	Ftaa          string         `json:"ftaa"`
	Bfaa          string         `json:"bfaa"`
	host          string
	proxy         string
	path          string
	httpClient    *http.Client
	Handle        string `json:"handle"`
	HandleOrEmail string `json:"handleOrEmail"`
	Password      string `json:"password"`
}

func createCodeforcesClient() *CodeforcesClient {
	jar, _ := cookiejar.New(nil)
	client := &CodeforcesClient{
		Jar:        jar,
		host:       "https://codeforces.com",
		path:       filepath.Join(utils.GetAppHomeDirectoryPath(), "appdata", "codeforces_session.json"),
		httpClient: nil,
	}
	if err := client.load(); err != nil {
		logger.Error(err.Error())
	}
	Proxy := http.ProxyFromEnvironment
	//if len(proxy) > 0 {
	//	proxyURL, err := url.Parse(proxy)
	//	if err != nil {
	//	} else {
	//		Proxy = http.ProxyURL(proxyURL)
	//	}
	//}
	client.httpClient = &http.Client{Jar: client.Jar, Transport: &http.Transport{Proxy: Proxy}}
	if err := client.save(); err != nil {
		logger.Error(err.Error())
	}
	return client
}

func (c *CodeforcesClient) login() (err error) {
	defer utils.PanicRecovery()

	logger.Info(fmt.Sprintf("login %v...\n", c.HandleOrEmail))

	password, err := utils.Decrypt(c.HandleOrEmail, c.Password)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	jar, _ := cookiejar.New(nil)
	c.httpClient.Jar = jar
	body, err := utils.GetBody(c.httpClient, c.host+"/enter")
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	csrf, err := findCsrf(body)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	ftaa := genFtaa()
	bfaa := genBfaa()

	body, err = utils.PostBody(c.httpClient, c.host+"/enter", url.Values{
		"csrf_token":    {csrf},
		"action":        {"enter"},
		"ftaa":          {ftaa},
		"bfaa":          {bfaa},
		"handleOrEmail": {c.HandleOrEmail},
		"password":      {password},
		"_tta":          {"253"},
		"remember":      {"on"},
	})
	if err != nil {
		return
	}

	handle, err := findHandle(body)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorLoginFailed)
	}

	c.Ftaa = ftaa
	c.Bfaa = bfaa
	c.Handle = handle
	c.Jar = jar
	fmt.Println("Succeed!!")
	fmt.Printf("Welcome %v~\n", handle)
	return c.save()
}

func (c *CodeforcesClient) DoLogin(handleOrEmail, password string) (err error, handle string) {
	defer utils.PanicRecovery()

	c.HandleOrEmail = handleOrEmail
	c.Password, err = utils.Encrypt(handleOrEmail, password)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError), ""
	}

	err = c.login()
	if err != nil {
		return
	}

	handle = c.Handle
	return nil, handle
}

func (c *CodeforcesClient) Submit(problem models.Problem, langId, source string) (err error) {
	defer utils.PanicRecovery()

	submitUrl := getSubmitUrl(c.host, problem.Url)

	body, err := utils.GetBody(c.httpClient, submitUrl)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	handle, err := findHandle(body)
	if err != nil {
		return
	}

	fmt.Printf("Current user: %v\n", handle)

	csrf, err := findCsrf(body)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	body, err = utils.PostBody(c.httpClient, fmt.Sprintf("%v?csrf_token=%v", submitUrl, csrf), url.Values{
		"csrf_token":            {csrf},
		"ftaa":                  {c.Ftaa},
		"bfaa":                  {c.Bfaa},
		"action":                {"submitSolutionFormSubmitted"},
		"submittedProblemIndex": {problem.Label},
		"programTypeId":         {langId},
		"contestId":             {problem.ContestId},
		"source":                {source},
		"tabSize":               {"4"},
		"_tta":                  {"594"},
		"sourceCodeConfirmed":   {"true"},
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	errMsg, err := findErrorMessage(body)
	if err == nil {
		logger.Error(errMsg)
		return errors.New(errMsg)
	}

	msg, err := findMessage(body)
	if err != nil {
		return errors.New(constants.ErrorSubmitFailed)
	}
	if !strings.Contains(msg, "submitted successfully") {
		return errors.New(msg)
	}
	socket.PublishStatusMessage("test_status", "Submitted successfully", "success")

	c.monitorSubmission(problem)

	c.Handle = handle
	return c.save()
}

func (c *CodeforcesClient) monitorSubmission(problem models.Problem) {
	defer utils.PanicRecovery()

	submissionId := ""
	submissionUrl := getSubmissionUrl(c.host, problem.Url)
	time.Sleep(2 * time.Second)
	for {
		st := time.Now()

		body, err := utils.GetBody(c.httpClient, submissionUrl)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		handle, err := findHandle(body)
		if err != nil {
			return
		}
		fmt.Printf("Current user: %v\n", handle)

		soupBody := soup.HTMLParse(string(body))
		var verdictElem soup.Root
		if submissionId == "" {
			submissionElem := soupBody.Find("span", "class", "submissionVerdictWrapper")
			submissionId = submissionElem.Attrs()["submissionid"]
			verdictElem = submissionElem.Find("span")
		} else {
			submissionElem := soupBody.Find("span", "submissionid", submissionId)
			verdictElem = submissionElem.Find("span")
		}

		verdictText := strings.TrimSpace(verdictElem.FullText())
		fmt.Println(submissionId, verdictText)

		socket.PublishStatusMessage("test_status", verdictText, "info")
		if !strings.HasPrefix(verdictText, "Running") && !strings.HasPrefix(verdictText, "In queue") {
			return
		}

		sub := time.Now().Sub(st)
		if sub < time.Second {
			time.Sleep(time.Duration(time.Second - sub))
		}
	}
}

func (c *CodeforcesClient) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, c)
}

func (c *CodeforcesClient) save() (err error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err == nil {
		err = os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(c.path, data, 0644)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return
}

func findHandle(body []byte) (string, error) {
	reg := regexp.MustCompile(`handle = "([\s\S]+?)"`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New(constants.ErrorNotLoggedIn)
	}
	return string(tmp[1]), nil
}

func findCsrf(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrf='(.+?)'`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("Cannot find csrf")
	}
	return string(tmp[1]), nil
}

func findMessage(body []byte) (string, error) {
	reg := regexp.MustCompile(`Codeforces.showMessage\("([^"]*)"\);\s*?Codeforces\.reformatTimes\(\);`)
	tmp := reg.FindSubmatch(body)
	if tmp != nil {
		return string(tmp[1]), nil
	}
	return "", errors.New("Cannot find any message")
}

func findErrorMessage(body []byte) (string, error) {
	reg := regexp.MustCompile(`error[a-zA-Z_\-\ ]*">(.*?)</span>`)
	tmp := reg.FindSubmatch(body)
	if tmp == nil {
		return "", errors.New("Cannot find error")
	}
	return string(tmp[1]), nil
}

func getSubmitUrl(host, url string) string {
	_, cid, _, ctype := utils.ExtractInfoFromUrl(url)
	return fmt.Sprintf(host+"/%v/%v/submit", ctype, cid)
}

func getSubmissionUrl(host, url string) string {
	oj, cid, _, ctype := utils.ExtractInfoFromUrl(url)
	if oj == "codeforces" {
		return fmt.Sprintf(host+"/%v/%v/my", ctype, cid)
	} else if oj == "atcoder" {
		return fmt.Sprintf(host+"/%v/%v/submissions/me", ctype, cid)
	}
	return ""
}

func genFtaa() string {
	return utils.RandString(18)
}

func genBfaa() string {
	return "f1b3f18c715565b589b7823cda7448ce"
}

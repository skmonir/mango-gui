package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/oj/cookiejar"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type CodeforcesClient struct {
	Jar           *cookiejar.Jar `json:"cookies"`
	Handle        string         `json:"handle"`
	HandleOrEmail string         `json:"handleOrEmail"`
	Password      string         `json:"password"`
	Ftaa          string         `json:"ftaa"`
	Bfaa          string         `json:"bfaa"`
	host          string
	proxy         string
	path          string
	httpClient    *http.Client
}

var once sync.Once
var cfClient *CodeforcesClient

func getCodeforcesClient() *CodeforcesClient {
	if cfClient == nil {
		once.Do(func() {
			cfClient = createClient()
		})
	}
	return cfClient
}

func createClient() *CodeforcesClient {
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
	logger.Info(fmt.Sprintf("login %v...\n", c.HandleOrEmail))

	password, err := decrypt(c.HandleOrEmail, c.Password)
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
		"_tta":          {"176"},
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
	c.HandleOrEmail = handleOrEmail
	c.Password, err = encrypt(handleOrEmail, password)
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
		return errors.New(constants.ErrorServerError)
	}

	errMsg, err := findErrorMessage(body)
	if err == nil {
		logger.Error(errMsg)
		return errors.New(constants.ErrorServerError)
	}

	msg, err := findMessage(body)
	if err != nil {
		return errors.New(constants.ErrorSubmitFailed)
	}
	if !strings.Contains(msg, "submitted successfully") {
		return errors.New(msg)
	}

	c.Handle = handle
	return c.save()
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

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
	"strings"
	"time"
)

type AtCoderClient struct {
	Jar           *cookiejar.Jar `json:"cookies"`
	host          string
	proxy         string
	path          string
	httpClient    *http.Client
	Handle        string `json:"handle"`
	HandleOrEmail string `json:"handleOrEmail"`
	Password      string `json:"password"`
}

var acClient *AtCoderClient

func getAtCoderClient() *AtCoderClient {
	if cfClient == nil {
		once.Do(func() {
			acClient = createAtCoderClient()
		})
	}
	return acClient
}

func createAtCoderClient() *AtCoderClient {
	jar, _ := cookiejar.New(nil)
	client := &AtCoderClient{
		Jar:        jar,
		host:       "https://atcoder.jp",
		path:       filepath.Join(utils.GetAppHomeDirectoryPath(), "appdata", "atcoder_session.json"),
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

func (c *AtCoderClient) login() (err error) {
	defer utils.PanicRecovery()

	logger.Info(fmt.Sprintf("login %v...\n", c.HandleOrEmail))

	password, err := utils.Decrypt(c.HandleOrEmail, c.Password)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	jar, _ := cookiejar.New(nil)
	c.httpClient.Jar = jar
	body, err := utils.GetBody(c.httpClient, c.host+"/login")
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	csrf, err := findAtCoderCsrf(soup.HTMLParse(string(body)))
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	body, err = utils.PostBody(c.httpClient, c.host+"/login", url.Values{
		"csrf_token": {csrf},
		"username":   {c.HandleOrEmail},
		"password":   {password},
	})
	if err != nil {
		return
	}

	err = findAtCoderHandle(c.HandleOrEmail, soup.HTMLParse(string(body)))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	c.Handle = c.HandleOrEmail
	c.Jar = jar
	fmt.Println("Succeed!!")
	fmt.Printf("Welcome %v~\n", c.HandleOrEmail)
	return c.save()
}

func (c *AtCoderClient) DoLogin(handleOrEmail, password string) (err error, handle string) {
	fmt.Println(handleOrEmail, password)
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

func (c *AtCoderClient) Submit(problem models.Problem, langId, source string) (err error) {
	defer utils.PanicRecovery()

	submitUrl := getSubmitUrl(c.host, problem.Url)

	body, err := utils.GetBody(c.httpClient, submitUrl)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	err = findAtCoderHandle(c.HandleOrEmail, soup.HTMLParse(string(body)))
	if err != nil {
		return
	}

	csrf, err := findAtCoderCsrf(soup.HTMLParse(string(body)))
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constants.ErrorServerError)
	}

	_, _, pid, _ := utils.ExtractInfoFromUrl(problem.Url)

	body, err = utils.PostBody(c.httpClient, submitUrl, url.Values{
		"csrf_token":          {csrf},
		"data.TaskScreenName": {pid},
		"data.LanguageId":     {langId},
		"sourceCode":          {source},
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = findAtCoderLoginErrorMessage(body)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	socket.PublishStatusMessage("test_status", "Submitted successfully", "success")

	c.monitorSubmission(problem)

	c.Handle = c.HandleOrEmail
	return c.save()
}

func (c *AtCoderClient) monitorSubmission(problem models.Problem) {
	defer utils.PanicRecovery()

	submissionId, verdict := "", ""
	submissionUrl := getSubmissionUrl(c.host, problem.Url)
	time.Sleep(2 * time.Second)
	for {
		st := time.Now()

		body, err := utils.GetBody(c.httpClient, submissionUrl)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		soupBody := soup.HTMLParse(string(body))

		err = findAtCoderHandle(c.HandleOrEmail, soupBody)
		if err != nil {
			return
		}

		trs := soupBody.Find("tbody").FindAll("tr")
		for _, tr := range trs {
			tds := tr.FindAll("td")
			probUrl := tds[1].Find("a").Attrs()["href"]
			subId := tds[4].Attrs()["data-id"]
			if strings.HasSuffix(problem.Url, probUrl) && (submissionId == "" || submissionId == subId) {
				verdict = strings.TrimSpace(tds[6].FullText())
				submissionId = subId
				break
			}
		}

		fmt.Println(submissionId, verdict)

		socket.PublishStatusMessage("test_status", verdict, "info")
		if !utils.SliceContains([]string{"WJ", "WR", "Judging"}, verdict) && !strings.Contains(verdict, "/") {
			return
		}

		sub := time.Now().Sub(st)
		if sub < time.Second {
			time.Sleep(time.Duration(time.Second - sub))
		}
	}
}

func (c *AtCoderClient) load() (err error) {
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

func (c *AtCoderClient) save() (err error) {
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

func findAtCoderHandle(handle string, soupBody soup.Root) error {
	strHtml := strings.ToLower(soupBody.HTML())
	if !strings.Contains(strHtml, fmt.Sprintf("<a href=\"/users/%v\">", handle)) {
		return errors.New(constants.ErrorNotLoggedIn)
	}
	return nil
}

func findAtCoderCsrf(soupBody soup.Root) (string, error) {
	csrfElem := soupBody.Find("input", "name", "csrf_token")
	if csrfElem.Error != nil {
		return "", errors.New("Cannot find csrf")
	}
	return csrfElem.Attrs()["value"], nil
}

func findAtCoderLoginErrorMessage(data []byte) error {
	defer utils.PanicRecovery()
	soupBody := soup.HTMLParse(string(data))
	errorElem := soupBody.Find("div", "class", "alert-danger")
	if errorElem.Error == nil {
		errMessage := errorElem.FullText()
		errMessage = strings.Replace(errMessage, "×", "", -1)
		errMessage = strings.TrimSpace(errMessage)
		return errors.New(errMessage)
	}
	return nil
}

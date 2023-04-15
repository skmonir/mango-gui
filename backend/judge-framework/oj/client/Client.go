package client

import (
	"errors"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"log"
	"net/http"
	"sync"
)

var once sync.Once

type IClient interface {
	DoLogin(handleOrEmail, password string) (err error, handle string)
	Submit(problem models.Problem, langId, source string) (err error)
}

func GetOjClientByPlatform(platform string) (error, IClient) {
	var err error
	var ojClient IClient
	if platform == "codeforces" {
		ojClient = createCodeforcesClient()
	} else if platform == "atcoder" {
		ojClient = createAtCoderClient()
	} else {
		log.Println("Unknown platform")
		err = errors.New("Unknown platform")
	}
	return err, ojClient
}

func GetHttpClientByPlatform(platform string) (error, *http.Client) {
	var err error
	var httpClient *http.Client
	if platform == "codeforces" {
		httpClient = createCodeforcesClient().httpClient
	} else if platform == "atcoder" {
		httpClient = createAtCoderClient().httpClient
	} else {
		log.Println("Unknown platform")
		err = errors.New("Unknown platform")
	}
	return err, httpClient
}

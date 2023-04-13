package client

import (
	"errors"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"log"
	"sync"
)

var once sync.Once

type IClient interface {
	DoLogin(handleOrEmail, password string) (err error, handle string)
	Submit(problem models.Problem, langId, source string) (err error)
}

func GetClientByPlatform(platform string) (error, IClient) {
	var err error
	var httpClient IClient
	if platform == "codeforces" {
		httpClient = createCodeforcesClient()
	} else if platform == "atcoder" {
		httpClient = createAtCoderClient()
	} else {
		log.Println("Unknown platform")
		err = errors.New("Unknown platform")
	}
	return err, httpClient
}

package client

import (
	"errors"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"log"
)

type IClient interface {
	DoLogin(handleOrEmail, password string) (err error, handle string)
	Submit(problem models.Problem, langId, source string) (err error)
}

func GetClientByPlatform(platform string) (error, IClient) {
	var err error
	var httpClient IClient
	if platform == "codeforces" {
		httpClient = getCodeforcesClient()
	} else if platform == "atcoder" {
		err = errors.New("AtCoder isn't supported at this moment")
	} else {
		log.Println("Unknown platform")
		err = errors.New("Unknown platform")
	}
	return err, httpClient
}

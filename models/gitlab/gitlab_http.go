package gitlab

import (
	"fmt"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
)

var client = &http.Client{}

func GetNewGitLabRequest(url string) *http.Request {
	fullUrl := fmt.Sprintf("%s/api/v4%s", settings.Config.GitLabInstanceURL, url)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("PRIVATE-TOKEN", settings.Config.BotPrivateToken)
	return req
}

func DoRequestBytes(req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	// https://docs.gitlab.com/ee/api/#status-codes
	if resp.StatusCode != 200 {
		return nil, errors.New("Non-200 response")
	}

	return body, nil
}

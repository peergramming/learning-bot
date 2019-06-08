package gitlab

import (
	"fmt"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"io/ioutil"
	"log"
	"net/http"
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

func DoRequestBytes(req *http.Request) []byte {
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	return body
}

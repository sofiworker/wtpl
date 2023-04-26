package req

import (
	reqV3 "github.com/imroc/req/v3"
	"time"
)

var (
	client *reqV3.Client
)

func init() {
	client = reqV3.NewClient()
	client.SetTimeout(3 * time.Second)
}

func Post(url string, body interface{}, respBody interface{}) error {
	response := client.Post(url).SetBody(body).Do()
	err := response.Unmarshal(respBody)
	if err != nil {
		return err
	}
	return nil
}

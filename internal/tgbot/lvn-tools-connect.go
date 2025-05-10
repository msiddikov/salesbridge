package tgbot

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	configData struct {
		Data string `json:"data"`
	}
)

const (
	host   = "https://tools.lavina.uz"
	key    = "Keir*o42Ysdf"
	secret = "Lkjsdf&&^&-sdlfksji54654"
)

// sends post to https://tools.lavina.tech/config/:id the config
func saveConfig(id, data string) error {
	// send config to lvn-tools
	cfg := configData{Data: data}

	// create request
	path := "/config/" + id
	body, _ := json.Marshal(cfg)
	req, err := http.NewRequest("POST", host+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// create signature
	signString := string(body) + secret
	sign := md5.Sum([]byte(signString))

	// add headers
	req.Header.Add("k", key)
	req.Header.Add("s", hex.EncodeToString(sign[:]))

	// send request
	_, err = http.DefaultClient.Do(req)
	return err
}

func getConfig(id, data string) (string, error) {
	// create request
	path := "/config/" + id
	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		return "", err
	}

	// create signature
	signString := path + secret
	sign := md5.Sum([]byte(signString))

	// add headers
	req.Header.Add("k", key)
	req.Header.Add("s", hex.EncodeToString(sign[:]))

	// send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// read response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// parse response
	var cfg configData
	err = json.Unmarshal(body, &cfg)
	if err != nil {
		return "", err
	}

	return cfg.Data, nil

}

func (s *Service) saveTopics() {

	// marshal topics
	data, err := json.Marshal(topics)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = saveConfig("topics", string(data))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Service) loadTopics() {
	data, err := getConfig("topics", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal([]byte(data), &topics)
	if err != nil {
		fmt.Println(err)
	}
}

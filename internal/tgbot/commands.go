package tgbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type body struct {
	Commands []command `json:"commands"`
}

func setCommands() {
	commands := getCommandList()
	body, _ := json.Marshal(commands)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.telegram.org/bot"+os.Getenv("TG_TOKEN")+"/setMyCommands", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	bodybytes, _ := ioutil.ReadAll(res.Body)

	if err != nil || res.StatusCode > 299 {
		fmt.Printf("%s: %s", err, bodybytes)
	}
}

func getCommandList() body {
	return body{
		Commands: []command{
			{
				Command:     "id",
				Description: "Returns chat id",
			},
			{
				Command:     "servers",
				Description: "Returns online servers",
			},
		},
	}
}

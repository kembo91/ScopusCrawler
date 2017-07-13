package query

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"fmt"

	"../config"
	"../storage"
)

func MakeQuery(address string, id string, params map[string]string, timeoutSec int,
	storage storage.MySqlStorage, config config.Configuration) (string, error) {
	requestPath := address
	if id != "" {
		requestPath = strings.Replace(requestPath, "{_id_}", id, -1)
		//requestPath = requestPath + "/" + id
	}
	for key, value := range params {
		requestPath += key + "=" + value + "&"
	}
	var data string
	var body []byte
	finishedRequest, _ := storage.GetFinishedRequest(requestPath)
	if finishedRequest == "" {
		request := requestPath
		authKey := config.GetKey()
		requestPath = requestPath + "apiKey=" + authKey
		fmt.Println(requestPath)
		req, err := http.NewRequest("GET", requestPath, nil)
		if err != nil {
			return "", err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			config.RemoveKey(authKey)
			return "", err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		data = string(body)
		err = storage.CreateFinishedRequest(request, data)
		if err != nil{
			return "", err
		}
	} else {
		data = finishedRequest
	}

	duration := time.Duration(timeoutSec) * time.Second
	time.Sleep(duration)
	return data, nil
}

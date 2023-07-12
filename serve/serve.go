package serve

import (
	"awesomeProject/module"
	"bufio"
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

func sendRequest(url, method string, payload *module.SendUuidMode, headers map[string]string) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return client.Do(req)
}

func ToGetUuid(uuid string) {
	url := "https://claude.sbai.chat/api/organizations/13b52fbc-790b-4f61-8da4-4abdd21d17a2/chat_conversations"
	method := "POST"
	payload := &module.SendUuidMode{
		UUID: uuid,
		Name: "",
	}
	headers := map[string]string{
		"Cookie": "sessionKey=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; intercom-device-id-lupk8zyo=7c3574c3-1390-49a8-a55c-79d218a83c53; intercom-session-lupk8zyo=MVo5a0JQMG1EemtRZC9tMXM4NlpQNWJ4Y0pKOHFwRzR2bGhmeTE4SkRUV2R4QnZsZlZaK3Q4dGV3cjZZT1BLVC0tYnNkMlRsbE52SC9VL2diVUtFSnpHUT09--bb47d23ce5d820b23ba28f17e387c6d700e5e1be",
		//"User-Agent":   "Apifox/1.0.0 (https://apifox.com)",
		"Content-Type": "application/json",
	}

	res, err := sendRequest(url, method, payload, headers)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)
	var response module.UuidResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return
	}
	if response.UUID == "" {
		myUuid = uuid
	} else {
		myUuid = response.UUID
	}
	fmt.Println(myUuid)

}

var myUuid string
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToSendMsg(jsonData chan<- []byte, wg *sync.WaitGroup, uuid string) {
	ToGetUuid(uuid)
	url := "https://claude.sbai.chat/api/append_message"
	method := "POST"

	p := module.SendMsgMode{
		Completion: struct {
			Prompt   string `json:"prompt"`
			Timezone string `json:"timezone"`
			Model    string `json:"model"`
		}{
			Prompt:   "",
			Timezone: "Asia/Shanghai",
			Model:    "claude-2",
		},
		OrganizationUUID: "13b52fbc-790b-4f61-8da4-4abdd21d17a2",
		ConversationUUID: myUuid,
		Text:             "写一个雪花算法",
		Attachments:      []string{},
	}

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	payload := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Cookie", "sessionKey=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; intercom-device-id-lupk8zyo=7c3574c3-1390-49a8-a55c-79d218a83c53; intercom-session-lupk8zyo=MVo5a0JQMG1EemtRZC9tMXM4NlpQNWJ4Y0pKOHFwRzR2bGhmeTE4SkRUV2R4QnZsZlZaK3Q4dGV3cjZZT1BLVC0tYnNkMlRsbE52SC9VL2diVUtFSnpHUT09--bb47d23ce5d820b23ba28f17e387c6d700e5e1be")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()

		splitBody := strings.SplitN(line, "data:", 2) // Use SplitN to split the line into at most 2 parts
		if len(splitBody) < 2 {
			continue
		}
		jsonBody := strings.TrimSpace(splitBody[1]) // Get the second part, which is the JSON string

		var result module.MsgResponse
		err = json.Unmarshal([]byte(jsonBody), &result)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gptMsg := module.JsonData{
			ID:      result.LogID,
			Object:  "chat.completion.chunk",
			Created: myUuid,
			Model:   result.Model,
			Choices: []module.Choice{
				{
					Index: 0,
					Delta: map[string]string{
						"content": result.Completion,
					},
					FinishReason: nil,
				},
			},
		}
		data, err := json.Marshal(gptMsg)
		if err != nil {
			log.Println(err)
		}

		jsonData <- data

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	defer wg.Done()
}

// Fetch data from the API
func FetchData(c chan<- string, wg *sync.WaitGroup) {
	url := "https://claude.sbai.chat/api/append_message"
	method := "POST"

	p := module.SendMsgMode{
		Completion: struct {
			Prompt   string `json:"prompt"`
			Timezone string `json:"timezone"`
			Model    string `json:"model"`
		}{
			Prompt:   "",
			Timezone: "Asia/Shanghai",
			Model:    "claude-2",
		},
		OrganizationUUID: "13b52fbc-790b-4f61-8da4-4abdd21d17a2",
		ConversationUUID: myUuid,
		Text:             "你好",
		Attachments:      []string{},
	}

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	payload := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Cookie", "sessionKey=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; intercom-device-id-lupk8zyo=7c3574c3-1390-49a8-a55c-79d218a83c53; intercom-session-lupk8zyo=MVo5a0JQMG1EemtRZC9tMXM4NlpQNWJ4Y0pKOHFwRzR2bGhmeTE4SkRUV2R4QnZsZlZaK3Q4dGV3cjZZT1BLVC0tYnNkMlRsbE52SC9VL2diVUtFSnpHUT09--bb47d23ce5d820b23ba28f17e387c6d700e5e1be")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		c <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	close(c)
	wg.Done()
}

// Process the fetched data
func ProcessData(c <-chan string, out chan<- []byte, wg *sync.WaitGroup) {
	for line := range c {
		splitBody := strings.SplitN(line, "data:", 2) // Use SplitN to split the line into at most 2 parts
		if len(splitBody) < 2 {
			continue
		}
		jsonBody := strings.TrimSpace(splitBody[1]) // Get the second part, which is the JSON string

		var result module.MsgResponse
		err := json.Unmarshal([]byte(jsonBody), &result)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gptMsg := module.JsonData{
			ID:      result.LogID,
			Object:  "chat.completion.chunk",
			Created: myUuid,
			Model:   result.Model,
			Choices: []module.Choice{
				{
					Index: 0,
					Delta: map[string]string{
						"content": result.Completion,
					},
					FinishReason: nil,
				},
			},
		}
		data, err := json.Marshal(gptMsg)
		if err != nil {
			log.Println(err)
		}

		out <- data
	}
	close(out)
	wg.Done()
}

// Send the processed data
func SendData(c <-chan []byte, jsonData chan<- []byte, wg *sync.WaitGroup) {
	for data := range c {
		jsonData <- data
	}

	wg.Done()
}

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

func ToGetUuid(uuid string, apikey string) {
	url := "https://claude.ai/api/organizations/13b52fbc-790b-4f61-8da4-4abdd21d17a2/chat_conversations"
	method := "POST"
	payload := &module.SendUuidMode{
		UUID: uuid,
		Name: "",
	}
	headers := map[string]string{
		"Cookie":       apikey,
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
func getOrganizationUuid(apikey string) string {
	method := "GET"
	url := "https://claude.ai/api/organizations"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	req.Header.Add("Cookie", apikey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println(string(body))

	var orgs []module.Organization
	err = json.Unmarshal(body, &orgs)
	if err != nil {
		log.Println(err)
		return ""
	}
	if len(orgs) > 0 {
		return orgs[0].Uuid
	}
	log.Println("No organizations found")
	return ""
}

var myUuid string
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToSendMsg(jsonData chan<- []byte, wg *sync.WaitGroup, uuid string, myreq module.OpenAIRequest, apiKey string) {
	ToGetUuid(uuid, apiKey)
	OrganizationUUID := getOrganizationUuid(apiKey)
	var messages string
	var usermsg string
	for _, msg := range myreq.Messages {
		// 根据msg.Role添加消息内容到messages字符串
		if msg.Role == "user" {
			messages += fmt.Sprintf("Human:%s\n\n", msg.Content)
		} else {
			usermsg += fmt.Sprintf("Assistant:%s\n\n", msg.Content)
		}
	}
	if usermsg == "" {
		usermsg = messages
	}
	if messages == "" {
		messages = usermsg
	}
	url := "https://claude.ai/api/append_message"
	method := "POST"

	p := module.SendMsgMode{
		Completion: struct {
			Prompt   string `json:"prompt"`
			Timezone string `json:"timezone"`
			Model    string `json:"model"`
		}{
			Prompt:   usermsg,
			Timezone: "Asia/Shanghai",
			Model:    "claude-2",
		},
		OrganizationUUID: OrganizationUUID,
		ConversationUUID: myUuid,
		Text:             messages,
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
	req.Header.Add("Cookie", apiKey)
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
		//gptMsg := module.ClaudeRes{
		//	Completion: result.Completion,
		//	StopReason: nil,
		//	Model:      result.Model,
		//}
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

func ToSendClaudeMsg(jsonData chan<- []byte, wg *sync.WaitGroup, uuid string, myreq module.AssistantRequest, apiKey string) {
	ToGetUuid(uuid, apiKey)
	OrganizationUUID := getOrganizationUuid(apiKey)
	//var msg rune
	//var messages, usermsg string
	//parts := strings.Split(myreq.Prompt, "\n\n")
	//for _, part := range parts {
	//	if strings.HasPrefix(part, "Human:") {
	//		messages += strings.TrimPrefix(part, "Human:")
	//	} else if strings.HasPrefix(part, "Assistant:") {
	//		usermsg += strings.TrimPrefix(part, "Assistant:")
	//	}
	//}
	//fmt.Println(myreq.Prompt)
	//fmt.Println(messages, usermsg)
	url := "https://claude.ai/api/append_message"
	method := "POST"

	p := module.SendMsgMode{
		Completion: struct {
			Prompt   string `json:"prompt"`
			Timezone string `json:"timezone"`
			Model    string `json:"model"`
		}{
			Prompt:   myreq.Prompt,
			Timezone: "Asia/Shanghai",
			Model:    "claude-2",
		},
		OrganizationUUID: OrganizationUUID,
		ConversationUUID: myUuid,
		Text:             myreq.Prompt,
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
	req.Header.Add("Cookie", apiKey)
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
		//gptMsg := module.JsonData{
		//	ID:      result.LogID,
		//	Object:  "chat.completion.chunk",
		//	Created: myUuid,
		//	Model:   result.Model,
		//	Choices: []module.Choice{
		//		{
		//			Index: 0,
		//			Delta: map[string]string{
		//				"content": result.Completion,
		//			},
		//			FinishReason: nil,
		//		},
		//	},
		//}
		gptMsg := module.ClaudeRes{
			Completion: result.Completion,
			StopReason: nil,
			Model:      result.Model,
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

package main

import (
	"awesomeProject/module"
	"awesomeProject/serve"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	//serve.ToGetUuid("bf52dc95-0ffa-4dc4-9b07-fdd90bea4035")
	//jsonData := make(chan []byte)
	//go func() {
	//	if err := serve.ToSendMsg(jsonData); err != nil {
	//		log.Println(err)
	//	}
	//}()
	//
	//for data := range jsonData {
	//	fmt.Printf("data:%s\n", string(data))
	//}
	r := gin.Default()
	//r.POST("/v2/chat/completions", func(c *gin.Context) {
	//	//c.Writer.Header().Set("Content-Type", "text/event-stream")
	//	//c.Writer.Header().Set("Cache-Control", "no-cache")
	//	//c.Writer.Header().Set("Connection", "keep-alive")
	//
	//	jsonData := make(chan []byte)
	//	var wg sync.WaitGroup
	//
	//	wg.Add(1)
	//	//go serve.ToSendMsg(jsonData, &wg, "bf52dc95-0ffa-4dc4-9b07-fdd90bea4031",)
	//
	//	go func() {
	//		wg.Wait()
	//		close(jsonData)
	//	}()
	//
	//	for d := range jsonData {
	//		c.Writer.Write(d)
	//		c.Writer.Write([]byte("\n"))
	//		c.Writer.Flush()
	//	}
	//})
	r.POST("/v1/chat/completions", func(c *gin.Context) {
		jsonData := make(chan []byte, 100)
		var req module.OpenAIRequest
		auth := c.GetHeader("Authorization")
		// 从Authorization字段中获取API密钥
		apiKey := strings.TrimPrefix(auth, "Bearer ")

		// 获取请求体
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var wg sync.WaitGroup

		// Generate a random UUID
		muid := uuid.New().String()

		// Start a goroutine to send requests
		wg.Add(1)
		go func() {
			start := time.Now()
			serve.ToSendMsg(jsonData, &wg, muid, req, apiKey)
			elapsed := time.Since(start)

			// Set the timeout to be twice the time it took to handle the request
			timeout := elapsed * 2
			time.AfterFunc(timeout, func() {
				close(jsonData)
			})
		}()

		// Use select to read from jsonData channel
		for d := range jsonData {
			// If data is received, write it to HTTP response
			c.Writer.Write(d)
			c.Writer.Write([]byte("\n"))
			c.Writer.Flush()
		}
	})

	r.POST("/v1/complete", func(c *gin.Context) {
		jsonData := make(chan []byte, 100)
		var req module.AssistantRequest
		auth := c.GetHeader("x-api-key")
		// 从Authorization字段中获取API密钥

		// 获取请求体
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var wg sync.WaitGroup

		// Generate a random UUID
		muid := uuid.New().String()

		// Start a goroutine to send requests
		wg.Add(1)
		go func() {
			start := time.Now()
			serve.ToSendClaudeMsg(jsonData, &wg, muid, req, auth)
			elapsed := time.Since(start)

			// Set the timeout to be twice the time it took to handle the request
			timeout := elapsed * 2
			time.AfterFunc(timeout, func() {
				close(jsonData)
			})
		}()

		// Use select to read from jsonData channel
		for d := range jsonData {
			// If data is received, write it to HTTP response
			c.Writer.Write([]byte("data: "))
			c.Writer.Write(d)
			c.Writer.Write([]byte("\n"))
			c.Writer.Flush()
		}
	})
	//r.GET("/getdata", func(c *gin.Context) {
	//	//c.Writer.Header().Set("Content-Type", "text/event-stream")
	//	//c.Writer.Header().Set("Cache-Control", "no-cache")
	//	//c.Writer.Header().Set("Connection", "keep-alive")
	//
	//	fetchedData := make(chan string)
	//	processedData := make(chan []byte)
	//	//jsonData := make(chan []byte)
	//	var wg sync.WaitGroup
	//
	//	wg.Add(2)
	//	go serve.FetchData(fetchedData, &wg)
	//	go serve.ProcessData(fetchedData, processedData, &wg)
	//
	//	wg.Wait()
	//
	//	for d := range processedData {
	//		c.Writer.Write(d)
	//		c.Writer.Write([]byte("\n\n"))
	//		c.Writer.Flush()
	//	}
	//})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//r.POST("/v3/chat/completions", func(c *gin.Context) {
	//	var req module.OpenAIRequest
	//
	//	// 获取请求头中的Authorization字段
	//	auth := c.GetHeader("Authorization")
	//	// 从Authorization字段中获取API密钥
	//	apiKey := strings.TrimPrefix(auth, "Bearer ")
	//
	//	// 输出API密钥
	//	fmt.Printf("API Key:%s\n", apiKey)
	//
	//	// 获取请求体
	//	if err := c.ShouldBindJSON(&req); err != nil {
	//		c.JSON(400, gin.H{"error": err.Error()})
	//		return
	//	}
	//
	//	// 输出请求体
	//	fmt.Printf("Model: %s, Messages: %v, Temperature: %f\n", req.Model, req.Messages, req.Temperature)
	//
	//	// 这里你可以处理你的逻辑
	//
	//	res := module.OpenAIResponse{
	//		Role:    "system",
	//		Content: "This is a test!",
	//	}
	//
	//	c.JSON(200, res)
	//})
	r.Run(":8084")
}

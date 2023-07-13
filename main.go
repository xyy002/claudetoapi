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
	r := gin.Default()
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
			timeout := elapsed
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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}

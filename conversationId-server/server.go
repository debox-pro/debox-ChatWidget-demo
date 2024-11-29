package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func calculateSignature(appSecret, nonce, timestamp string) string {
	h := sha1.New()
	h.Write([]byte(appSecret + nonce + timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
    // Load environment variables
    if err := godotenv.Load(".env.local"); err != nil {
        fmt.Println("Error loading .env file")
    }
	appSecret := os.Getenv("APP_SECRET")
	apiKey    := os.Getenv("API_KEY")
	appId     := os.Getenv("APP_ID")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Define a route on the severside to get conversation id
	r.GET("/getconversationid", func(c *gin.Context) {

		// Get parameters from query string
		groupName := c.Query("group_name")
		chainID := c.Query("chain_id")
		contractAddress := c.Query("contract_address")

		if groupName == "" && chainID == "" && contractAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
			return
		}

		// setting request url
		url := "http://open.debox.pro/openapi/chatwidget/conversation/id"
		method := "POST"

		nonce := fmt.Sprintf("%08d", rand.Intn(100000000))
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		signature := calculateSignature(appSecret, nonce, timestamp)

		// Construct request payload
		payloadMap := make(map[string]interface{})
		if groupName != "" {
			payloadMap["group_name"] = groupName
		}
		if chainID != "" {
			chainIDInt, err := strconv.Atoi(chainID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chain_id format"})
				return
			}
			payloadMap["chain_id"] = chainIDInt
		}
		if contractAddress != "" {
			payloadMap["contract_address"] = contractAddress
		}

		// Convert payload to JSON
		payloadBytes, err := json.Marshal(payloadMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to construct request payload"})
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest(method, url, strings.NewReader(string(payloadBytes)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add headers
		req.Header.Add("app_id", appId)
		req.Header.Add("X-API-KEY", apiKey)
		req.Header.Add("nonce", nonce)
		req.Header.Add("timestamp", timestamp)
		req.Header.Add("signature", signature)
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "*/*")

		// Send request
		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response to client
		c.Data(res.StatusCode, "application/json", body)
	})

	// Start server
	r.Run(":1444")
}

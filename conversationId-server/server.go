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
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var (
	rdb *redis.Client
)

func initializeRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("Invalid REDIS_DB value in .env file:", err)
		return nil
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Check the redis connection
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		fmt.Println("Redis connection failed:", err)
		return nil
	}
	return rdb
}

func calculateSignature(appSecret, nonce, timestamp string) string {
	h := sha1.New()
	h.Write([]byte(appSecret + nonce + timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	// Load environment variables from .env.local file
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("Error loading .env file")
	}
	appSecret := os.Getenv("APP_SECRET")
	apiKey := os.Getenv("API_KEY")
	appId := os.Getenv("APP_ID")
	port := os.Getenv("PORT")

	redisExpiryTime := os.Getenv("REDIS_EXPIRY_TIME")
	redisExpiryDuration, err := time.ParseDuration(redisExpiryTime)
	if err != nil {
		fmt.Println("Invalid REDIS_EXPIRY_TIME value in .env file:", err)
		return
	}

	// Initialize Redis
	rdb := initializeRedis()
	if rdb == nil {
		fmt.Println("Redis initialization failed, exiting")
		return
	}
	fmt.Println("Redis initialized")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Define a route to get conversation ID
	r.GET("/getconversationid", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get parameters from the query string
		groupName := c.Query("group_name")
		chainID := c.Query("chain_id")
		contractAddress := c.Query("contract_address")

		if groupName == "" && chainID == "" && contractAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters", "success": false})
			return
		}

		// Set request URL
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chain_id format", "success": false})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to construct request payload", "success": false})
			return
		}

		// Redis cache key
		cacheKey := fmt.Sprintf("conversation:%s:%s:%s", groupName, chainID, contractAddress)

		// Check if the response is already cached
		cachedResponse, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache hit, return cached data
			c.Data(http.StatusOK, "application/json", []byte(cachedResponse))
			return
		} else if err != redis.Nil {
			// Redis error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server Redis error", "success": false})
			return
		}

		// Make the request to Debox Open API
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
		req.Header.Add("Content-Type", "application/json")

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

		// Cache the response for future use
		err = rdb.Set(ctx, cacheKey, body, redisExpiryDuration).Err()
		if err != nil {
			fmt.Println("Failed to cache the response:", err)
		}

		c.Data(res.StatusCode, "application/json", body)
	})

	// Start the server
	fmt.Println("Server is running on port " + port)
	r.Run(":" + port)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CatResponse struct {
	Message string `json: "message"`
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://foo.com", "http://localhost:8888"},
		AllowMethods:  []string{"GET", "PUT", "PATCH"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	r.GET("/bark", func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "spice.com")
		// c.Header("Access-Control-Allow-Methods", "GET")

		mode := append(c.Request.Header["X-Mode"], "ACTIVE")[0]
		url := c.Request.Host + c.Request.URL.Path
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("[Dog Server v7.7 - %s from %s] %s", mode, url, callCat(mode))})
	})

	r.Run()
}

func getBaseURL(mode string) string {
	url := os.Getenv("BACKEND_BASE_URL")
	port := os.Getenv(fmt.Sprintf("CAT_SERVICE_%s_SERVICE_PORT", mode))

	if url == "" {
		url = "localhost:3000"
	} else {
		url = fmt.Sprintf("%s:%s", url, port)
		// url = fmt.Sprintf("%s-%s", url, mode)
	}

	return url
}

func callCat(mode string) string {
	// url := fmt.Sprintf("http://%s/meow", getBaseURL(mode))
	url := fmt.Sprintf("http://%s/animal/cat/meow", getBaseURL(mode))

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	var thisRes CatResponse
	parseErr := json.Unmarshal(body, &thisRes)

	if parseErr != nil {
		panic(parseErr)
	}
	return thisRes.Message
}

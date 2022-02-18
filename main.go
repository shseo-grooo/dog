package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CatResponse struct {
	Message string `json: "message"`
}

func main() {
	r := gin.Default()

	r.GET("/bark", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")

		mode := append(c.Request.Header["X-Mode"], "")[0]

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("[Dog Server v6.13 - %s] %s", mode, callCat(mode))})
	})

	r.Run()
}

func getBaseURL(mode string) string {
	url := os.Getenv("BACKEND_BASE_URL")

	if url == "" {
		url = "localhost:3000"
	} else {
		url = fmt.Sprintf("%s-%s", url, mode)
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

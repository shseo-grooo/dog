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

		mode := getMode(c.Request.Header["X-Mode"])

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("[Dog Server v6.3 - %s] %s", mode, callCat(mode))})
	})

	r.Run()
}

func getMode(header []string) (mode string) {
	x_mode := header
	x_mode = append(x_mode, "")

	if x_mode[0] == "active" || x_mode[0] == "preview" {
		mode = x_mode[0]
	} else {
		mode = "active"
	}

	return
}

func getBaseURL(mode string) string {
	url := os.Getenv("BACKEND_BASE_URL")

	if url == "" {
		url = "localhost:3000"
	} else {
		var port string
		if mode == "preview" {
			port = "8888"
		} else {
			port = "80"
		}
		url = fmt.Sprintf("%s:%s", url, port)
	}

	return url
}

func callCat(mode string) string {
	url := fmt.Sprintf("http://%s/animal/cat/meow", getBaseURL(mode))

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
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

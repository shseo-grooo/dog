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

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("[Dog Server v6.10 - %s] %s", mode, callCat(mode))})
	})

	r.Run()
}

func getMode(header []string) (mode string) {
	x_mode := header
	x_mode = append(x_mode, "")

	if x_mode[0] == "ACTIVE" || x_mode[0] == "PREVIEW" {
		mode = x_mode[0]
	} else {
		mode = "ACTIVE"
	}

	return
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

	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Add("X-Mode", mode)
	// client := &http.Client{}
	// resp, err := client.Do(req)
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

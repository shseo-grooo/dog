package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/bark", func(c *gin.Context) {
		mode := getMode(c.Request.Header["X-Mode"])
		c.String(http.StatusOK, "[Dog Server v3 - %s] %s", mode, callCat(mode))
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

func callCat(mode string) string {
	url := fmt.Sprintf("https://cat-service-%s/meow", mode)
	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("X-Mode", mode)

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	return string(bytes)
}

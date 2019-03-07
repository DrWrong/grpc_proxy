package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(ReverseProxy())
	r.POST("/*path", handleProxy)
	r.Run(":7001")
}

func handleProxy(c *gin.Context) {
	host := c.Request.URL.Host
	path := c.Param("path")

	response, err := getGrpcResponse(host, path, c.Request.Body)
	if err != nil {
		fmt.Printf("Error is %+v\n", err)
		c.AbortWithStatus(400)
		return
	}

	c.Data(200, "application/json", response)

}

func getGrpcResponse(host string, path string, body io.Reader) ([]byte, error) {
	path = strings.TrimLeft(path, "/")
	fmt.Println("Host: ", host)
	fmt.Println("Path: ", path)
	cmd := exec.Command(
		"grpcurl",
		"-d",
		"@",
		"-plaintext",
		host,
		path,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = body
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		print(out.String())
		print(stderr.String())
		return nil, err
	}

	return out.Bytes(), nil
}

func ReverseProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("host: " + c.Request.Host)
		if c.GetHeader("proxy-grpc") != "" {
			log.Println("continue for next")
			c.Next()
		}
		httpClient := http.Client{}
		c.Request.RequestURI = ""
		resp, err := httpClient.Do(c.Request)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for h, val := range resp.Header {
			c.Writer.Header()[h] = val
		}

		bodyContent, _ := ioutil.ReadAll(resp.Body)
		c.Writer.Write(bodyContent)
		c.Abort()
	}
}

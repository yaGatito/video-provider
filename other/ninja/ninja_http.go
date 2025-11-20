package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type loggingMiddleWare struct {
	logger io.Writer
	next   http.RoundTripper
}

func (l *loggingMiddleWare) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := fmt.Fprintf(l.logger, "[%s] %s %s %s\n", time.Now().Format(time.RFC3339), req.Method, req.RequestURI, req.Host)
	if err != nil {
		return nil, err
	}
	res, err := l.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main2() {
	client := http.Client{
		Transport: &loggingMiddleWare{
			logger: os.Stdout,
			next:   http.DefaultTransport,
		},
		Timeout: time.Second * 10,
	}
	resp, err := client.Get("https://rest.coincap.io/v3/assets")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("Status code:", resp.StatusCode, err)
	fmt.Println(string(bytes), err)
}

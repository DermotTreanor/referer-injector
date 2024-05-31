package main

import (
	"io"
	"net/http"
	"os"
	"fmt"
)

func main() {
	resp := proxy()
	io.Copy(os.Stdout, resp.Body)
}


func proxy() *http.Response{
	request, err := http.NewRequest("GET", "http://www.something.com/", nil)
	if err != nil{
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil{
		fmt.Println(err)
	}
	return resp
}
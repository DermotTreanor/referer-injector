package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

var addr string = "localhost:8080"
func main() {
	// response, err := http.Get("https://www.youtube.com")
	// if err != nil{
	// 	fmt.Println(err)
	// }
	// fmt.Println(response)

	http.HandleFunc("/", our_handler)
	//Remember that this will run the handler as a goroutine. 
	//So, expect concurrency if using a browser as another automatic favicon request launches
	err := http.ListenAndServe(addr, nil)
	if err != nil{
		log.Println(err)
	}
}

func our_handler(rw http.ResponseWriter, req *http.Request) {	
	alt := proxy_request(req)

	if alt == nil{
		fmt.Fprint(rw, "There was a problem with the proxy response")
		return
	}
	for k, v := range alt.Header{
		fmt.Println(k, v)
			rw.Header()[k] = v
	}
	
	my_file, err:= os.OpenFile("/home/dermottreanor/test", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		fmt.Printf("The file could not be opened: %v\n", err)
	}
	defer my_file.Close()

	byte_store := []byte{}
	buf := bytes.NewBuffer(byte_store)
	io.Copy(buf, alt.Body)
	buf2 := *buf

	io.Copy(rw, buf)
	io.Copy(my_file, &buf2)
	//rw.Header().Add("Request URL", "localhost:8080")
}

func proxy_request(req *http.Request) *http.Response{
	//PROBLEM: The response shouldn't be what we are getting back
	//We will need to try and fix this part before we worry about
	//whether we are sending it back correctly to the client.
	var err error
	fmt.Println(req.Host, req.URL.Path)
	req.RequestURI = ""

	newHost := ""
	newPath := ""
	endHostInd := strings.Index(req.URL.Path[1:], "/")
	if endHostInd >= 0{
		newHost = req.URL.Path[1:endHostInd + 1]
		newPath = req.URL.Path[endHostInd + 1:]
	} else{
		newHost = req.URL.Path[1:]
	}
	req.URL, err = url.Parse("http://www." + newHost + ".com" + newPath)
	if err != nil{
		fmt.Printf("Parse Error when creating new response's url, %v\n", err)
		return nil
	}
	req = req.Clone(context.Background())
	req.Host = "www." + newHost
	fmt.Println(req.Host, req.URL.Path)

	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		fmt.Printf("Error when using incoming request to get proxy response: %v\n", err)
		return nil
	}

	fmt.Printf("%v\n", resp)
	return resp
}


func reflecting(req *http.Request){
	reflection := reflect.ValueOf(*req.URL)
	for i := 0; i < reflection.Type().NumField(); i++{
		field := reflection.Type().Field(i)
		fmt.Printf("Type: %s. Key: %s. Value: %v\n\n", field.Name, field.Type, reflection.Field(i).Interface())
	}
}
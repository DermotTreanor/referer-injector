package main

import (
	"fmt"
	"net/http"
	"log"
	"reflect"
	"net/url"
	"context"
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
	go proxy_request(req)
	rw.Header()["Date"] = nil
	for k, v := range req.Header{
		rw.Header()[k] = v
	}
	fmt.Fprintf(rw, "The Headers we received are as follows:\n%s\n\n\nHere are the response headers to send to you:\n%s", req.Header, rw.Header())

}

func proxy_request(req *http.Request){
	var err error
	req.RequestURI = ""
	req.URL, err = url.Parse("https://www.youtube.com/hello")
	if err != nil{
		fmt.Printf("PARSE ERROR, %v. TYPE %T\n\n.", err, err)
	}
	req = req.Clone(context.Background())
	req.Host = "youtube.com"


	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		fmt.Printf("GOT ERROR: %v. TYPE, %T\n\n", err, err)
	}
	fmt.Printf("%s\n\n\n", resp)
}


func reflecting(req *http.Request){
	reflection := reflect.ValueOf(*req.URL)
	for i := 0; i < reflection.Type().NumField(); i++{
		field := reflection.Type().Field(i)
		fmt.Printf("Type: %s. Key: %s. Value: %v\n\n", field.Name, field.Type, reflection.Field(i).Interface())
	}
}
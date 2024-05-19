package main

import(
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/", our_handler)
	http.ListenAndServe("localhost:8080", nil)
}

func our_handler(rw http.ResponseWriter, req *http.Request){
	fmt.Fprintf(rw, "The method you used was: %s", req.Method)
}
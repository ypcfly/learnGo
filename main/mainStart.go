package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"program/com.ypc/learnGo/webService"
)

func helloHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("----> handler start <----")
	//writer.Header().Set("content-type","application/json")
	fmt.Fprintf(writer, "hello go, rquest url is %s",req.URL.Path)
}

// 自定义全局变量
var mu webService.CustomMux

func main() {
	fmt.Println("----> go start <----")
	//http.HandleFunc("/hello", helloHandler)
	//http.HandleFunc("/index",indexHandler)

	mux := &webService.CustomMux{}
	err := http.ListenAndServe("localhost:9090", mux)
	//err := http.ListenAndServe("localhost:9090", &mu)
	checkErr(err)
}
func indexHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("----> index handler start <----")
	t,err := template.ParseFiles("template/index.html")
	checkErr(err)
	t.Execute(writer,nil)
}


// 处理错误函数
func checkErr(e error) {
	if e != nil{
		log.Fatal(e)
	}
}


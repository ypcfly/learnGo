package webService

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type CustomMux struct {
	
}

func (mux *CustomMux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		IndexHandler(writer,request)
		return
	case "/home":
		HomePageHandler(writer,request)
		return
	case "/login":
		LoginHandler(writer,request)
		return
	default:
		http.NotFound(writer,request)
	}
}
func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> login handler start <-----")
	if request.Method == "Get" {
		t,err := template.ParseFiles("template/login.html")
		checkErr(err)
		t.Execute(writer,request)
	} else {
		username := request.FormValue("username")
		password := request.Form["password"]
		fmt.Println("username: ", username)
		fmt.Println("password: ", password)
		http.Redirect(writer,request,"/home",302)
	}
}
func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> home page handler start <-----")
	t,err := template.ParseFiles("template/homePage.html")
	checkErr(err)
	t.Execute(writer,"这是网站首页")
}
func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> index handler start <-----")
	t,err := template.ParseFiles("template/login.html")
	checkErr(err)
	Name := "go template"
	t.Execute(writer,Name)
}
func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

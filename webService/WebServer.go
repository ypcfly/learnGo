package webService

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"program/com.ypc/learnGo/database"
	"program/com.ypc/learnGo/model"
	"strconv"
)

type CustomMux struct {
}

func (mux *CustomMux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		IndexHandler(writer, request)
		return
	case "/home":
		HomePageHandler(writer, request)
		return
	case "/login":
		LoginHandler(writer, request)
		return
	case "/add":
		RenderAddHandler(writer, request)
		return
	case "/user/insert":
		InsertUserHandler(writer, request)
		return
	case "/user/select":
		QueryByIdHandler(writer, request)
		return
	case "/upload/picture":
		UploadPictureHandler(writer, request)
		return
	case "/do/upload":
		DoUploadActionHandler(writer, request)
		return
	case "/json/param":
		JsonHandler(writer, request)
		return
	default:
		http.NotFound(writer, request)
	}
}

// 请求参数是json类型
func JsonHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> request json param handler start <-----")
	bytes, err := ioutil.ReadAll(request.Body)
	checkErr(err)
	var user model.User
	//讲json参数和对应的model进行映射
	error := json.Unmarshal(bytes, &user)
	checkErr(error)

	// 保存到数据库
	db := database.GetConn()
	insert := "insert into t_user (username,password,address,age,mobile,sex,status) values($1,$2,$3,$4,$5,$6,$7)"
	rs, _ := db.Exec(insert, user.Username, user.Password, user.Address, user.Age, user.Mobile, user.Sex, 1)
	count, _ := rs.RowsAffected()

	var res model.ComRes
	if count != 1 {
		res.Code = "0002"
		res.Success = false
		res.Message = "insert to database failed"
	} else {
		res.Code = "0001"
		res.Success = true
		res.Message = "insert to database success"
	}
	// 返回自定义的响应结果
	writer.Header().Set("Content-Typ", "application/json")
	json.NewEncoder(writer).Encode(res)
}

// 处理上传逻辑
func DoUploadActionHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> upload picture handler start <-----")
	if request.Method == "GET" {
		t, _ := template.ParseFiles("template/error.html")
		t.Execute(writer, request)
	} else {
		file1, head, err := request.FormFile("picture1")
		checkErr(err)
		fmt.Println("----> filename: " + head.Filename + " <----")
		defer file1.Close()

		fileBytes, err := ioutil.ReadAll(file1)
		checkErr(err)
		// 图片类型
		fileType := http.DetectContentType(fileBytes)
		fmt.Println(fileType)

		// 创建存放图片文件夹
		dest := "/home/ypcfly/ypcfly/upload/"
		exist := dirExist(dest)
		if exist {
			fmt.Println("----> directory has exist <----")
		} else {
			error := os.Mkdir(dest, os.ModePerm)
			checkErr(error)
		}

		newFile, err := os.Create(dest + head.Filename)
		checkErr(err)
		defer newFile.Close()
		len, err := newFile.Write(fileBytes)
		if err != nil {
			fmt.Println("----> error occurred while write file to disk <----")
		}

		fmt.Println(len)
		t, _ := template.ParseFiles("template/picture.html")
		t.Execute(writer, dest+head.Filename)
	}

}

// 文件是否存在
func dirExist(s string) bool {
	var exist = true
	if _, err := os.Stat(s); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 跳转上传图片
func UploadPictureHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> forward to picture file upload handler start <-----")
	t, err := template.ParseFiles("template/upload.html")
	checkErr(err)
	t.Execute(writer, nil)
}

// 根据id查找用户信息
func QueryByIdHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> query by id handler start <-----")
	//id := request.FormValue("id")
	id := 1
	//query := request.URL.Query()
	//id := query["id"][0]
	fmt.Println(id)
	db := database.GetConn()
	query := "select username,sex,age,address,mobile,role from t_user where id = $1"
	//rows := db.QueryRow(query,id)
	rows, err := db.Query(query, id)
	checkErr(err)
	// 获取用户信息

	var user model.User
	for rows.Next() {
		//var username string
		//var sex string
		//var age int
		//var address string
		//var mobile string
		//var role string
		//err = rows.Scan(&username,&sex,&age,&address,&mobile,&role)

		err = rows.Scan(&user.Username, &user.Sex, &user.Age, &user.Address, &user.Mobile, &user.Role)
		checkErr(err)
	}
	// 将用户数据写回页面
	t, _ := template.ParseFiles("template/user/userDetail.html")
	t.Execute(writer, user)

	// 返回json数据
	//writer.Header().Set("content-type","application/json")
	//json.NewEncoder(writer).Encode(user)
}

// 添加新用户
func InsertUserHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> insert new user handler start <-----")
	// 获取表单参数
	username := request.PostFormValue("username")
	password := request.PostFormValue("password")
	age := request.PostFormValue("age")
	mobile := request.PostFormValue("mobile")
	address := request.PostFormValue("address")
	//省略校验逻辑....

	// 插入数据库
	db := database.GetConn()
	insert := "insert into t_user (username,password,address,age,mobile,status) values($1,$2,$3,$4,$5,$6)"
	//stmt,err := db.Prepare(insert)
	//checkErr(err)
	//rs,err := stmt.Exec(username,password,address,age,mobile,1)
	//checkErr(err)
	rs, err := db.Exec(insert, username, password, address, age, mobile, 1)
	checkErr(err)
	row, err := rs.RowsAffected()
	checkErr(err)
	fmt.Println("----> insert account = " + strconv.FormatInt(row, 10) + " <----")
	if row != 1 {
		log.Fatal("----> error occurred <----")
	} else {
		fmt.Println("----> insert user to database succeed <----")
	}
	//关闭数据库连接
	checkErr(db.Close())
}

// 跳转新加用户页面
func RenderAddHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> insert handler start <-----")
	t, err := template.ParseFiles("template/user/addUser.html")
	checkErr(err)
	t.Execute(writer, nil)
}

// 处理用户登录
func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> login handler start <-----")
	if request.Method == "Get" {
		t, err := template.ParseFiles("template/login.html")
		checkErr(err)
		t.Execute(writer, request)
	} else {
		username := request.FormValue("username")
		password := request.Form["password"]
		fmt.Println("username: ", username)
		fmt.Println("password: ", password)
		http.Redirect(writer, request, "/home", 302)
	}
}

// 跳转home page页面
func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> home page handler start <-----")
	t, err := template.ParseFiles("template/homePage.html")
	checkErr(err)
	t.Execute(writer, "这是网站首页")
}
func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("-----> index handler start <-----")
	t, err := template.ParseFiles("template/login.html")
	checkErr(err)
	t.Execute(writer, nil)
}

// 处理错误
func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

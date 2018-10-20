package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DBType struct{

}

func GetConn() {
	fmt.Println("----> get postgresql connection <----")
	db,err := sql.Open("postgres", "root:123456@tcp(127.0.0.1:5432)/pgsql?charset=utf8")
	checkErr(err,"-----> open datasources failed <-----")
	err = db.Ping()
	checkErr(err,"-----> get connection failed <-----")


	// 查询
	queryFunc(db)
	// 更新

	// 插入

	// 删除

}
func queryFunc(db *sql.DB) {
	selectById := "select * from user"
	rows,err := db.Query(selectById)
	checkErr(err,err.Error())

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}
}
func checkErr(e error,msg string) {
	if e != nil {
		log.Fatal(e,msg)
	}
}
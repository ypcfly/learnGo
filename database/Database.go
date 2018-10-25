package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//var DBType struct{
//	db *sql.DB
//	err error
//}

func GetConn() *sql.DB {
	sqlURL := "postgres://postgres:123456@localhost/pgsql?sslmode=disable"
	fmt.Println("----> get postgresql connection <----")
	db, err := sql.Open("postgres", sqlURL)
	checkErr(err, "-----> open datasources failed <-----")

	// 返回数据库连接
	return db
}
func queryFunc(db *sql.DB) {
	selectById := "select * from user"
	rows, err := db.Query(selectById)
	checkErr(err, err.Error())

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
func checkErr(e error, msg string) {
	if e != nil {
		log.Fatal(e, msg)
	}
}

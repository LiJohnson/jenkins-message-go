package test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE string = "/tmp/foo.db"

func Test_db_main(t *testing.T) {
	db, err := sql.Open("sqlite3", DB_FILE)
	checkErr(err)

	table := `
    CREATE TABLE IF NOT EXISTS userinfo (
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(128) NULL,
		departname VARCHAR(128) NULL,
        created DATE NULL
    );
    `
	create_table, err := db.Exec(table)
	checkErr(err)
	fmt.Println("create_table", create_table)
	//插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "12🌧研发部门", time.Now())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println("LastInsertId", id)
	//更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid, username, department, created)
	}

	// //删除数据
	// stmt, err = db.Prepare("delete from userinfo where uid=?")
	// checkErr(err)

	// res, err = stmt.Exec(id)
	// checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	db.Close()

}

func Test_insert(t *testing.T) {
	db, err := sql.Open("sqlite3", DB_FILE)
	checkErr(err)
	for i := 0; i < 100; i++ {
		stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
		checkErr(err)
		_, err = stmt.Exec("astaxie"+fmt.Sprintf("_%v", i), "研发部门", time.Now())
		checkErr(err)
	}
	//插入数据

}

func checkErr(err error) {
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}
}

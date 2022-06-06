package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CacheFile struct {
	Id         int64
	FileName   string
	Content    []byte
	CreateTime time.Time
}

const CACHE_SIZE int32 = 300

type Cache struct {
	Index  int32
	DBFile string
	db     *sql.DB
}

//添加消息
func (c *Cache) addFile(name string, content []byte) int64 {
	stmt, err := c.db.Prepare("INSERT INTO build_log_message(file_name, content, create_time) values(?,?,?)")
	if err != nil {
		log.Println(err)
		return 0
	}
	res, err := stmt.Exec(name, content, time.Now())
	if err != nil {
		log.Println(err)
		return 0
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0
	}
	return id
}

//删除消息
func (c *Cache) delFile(ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("ids not empty")
	}
	var deleteSql strings.Builder
	var idData = make([]interface{}, len(ids))
	deleteSql.WriteString("DELETE FROM build_log_message WHERE id in (?")
	idData[0] = ids[0]
	for i := 1; i < len(ids); i++ {
		deleteSql.WriteString(",?")
		idData[i] = ids[i]
	}
	deleteSql.WriteString(")")
	res, err := c.db.Exec(deleteSql.String(), idData...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return affect, nil
}

//获取消息
func (c *Cache) getFiles(size int32) ([]CacheFile, error) {
	if size <= 0 {
		return nil, errors.New("size must > 0")
	}
	rows, err := c.db.Query("SELECT id , file_name , content , create_time FROM build_log_message order by id limit ?", size)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var list []CacheFile
	for rows.Next() {
		var cacheFile CacheFile
		err = rows.Scan(&cacheFile.Id, &cacheFile.FileName, &cacheFile.Content, &cacheFile.CreateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, cacheFile)
	}
	return list, nil
}

//初始化 DB/Table
func (c *Cache) init() error {
	db, err := sql.Open("sqlite3", c.DBFile)
	if err != nil {
		return err
	}
	c.db = db
	table := `
    CREATE TABLE IF NOT EXISTS build_log_message (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        file_name VARCHAR(128) NULL,
		content VARCHAR(10240) NULL,
        create_time DATE NULL
    );
    `
	_, err = db.Exec(table)
	if err != nil {
		return err
	}
	return nil
}

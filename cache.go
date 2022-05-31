package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

type CacheFile struct {
	FileName string
	Content  []byte
}

const CACHE_SIZE int32 = 300

type Cache struct {
	List     [CACHE_SIZE]CacheFile
	Index    int32
	JsonFile string
}

//添加消息
func (c *Cache) addFile(name string, content []byte) {
	cur := c.Index % CACHE_SIZE
	c.List[cur] = CacheFile{FileName: name, Content: content}
	c.Index = c.Index + 1
	c.save()
}

//获取消息
func (c *Cache) getFiles(size int32) ([]CacheFile, error) {
	if size <= 0 {
		return nil, errors.New("size must > 0")
	}
	if size > CACHE_SIZE {
		return nil, fmt.Errorf("size must < %v", CACHE_SIZE)
	}

	cur := c.Index % CACHE_SIZE
	list := c.List[:cur]

	if c.Index > CACHE_SIZE {
		app := c.List[cur:]
		list = append(app, list...)
	}
	first := len(list) - int(size)
	if first < 0 {
		first = 0
	}
	return list[first:], nil
}

//初始化，从jons文件中加载消息
func (c *Cache) init() error {
	jsonData, err := ioutil.ReadFile(c.JsonFile)
	if err != nil {
		return nil
	}
	var tmp Cache
	err = json.Unmarshal(jsonData, &tmp)
	if err != nil {
		return nil
	}
	c.Index = tmp.Index
	c.List = tmp.List
	return nil
}

var sysLock sync.Mutex

// 将信息写到文件
func (c *Cache) save() error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	sysLock.Lock()
	err = ioutil.WriteFile(c.JsonFile, jsonData, 0600)
	log.Println("write file error", c.JsonFile, err)
	sysLock.Unlock()
	return nil
}

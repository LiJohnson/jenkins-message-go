package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gomarkdown/markdown"
)

const TMP string = "/tmp/logs"

type Message struct {
	Markdown map[string]string `json:"markdown"`
	File     map[string]string `json:"file"`
}

var messageCache Cache = Cache{jsonFile: fmt.Sprintf("%v/cache.db.json", TMP)}

//处理构建进度消息
func sendMessage(hub *Hub) func(http.ResponseWriter, *http.Request) {
	messageCache.init()
	return func(w http.ResponseWriter, r *http.Request) {
		// r.GetBody()
		var message Message
		json.NewDecoder(r.Body).Decode(&message)
		var content string
		if message.Markdown != nil {
			content = message.Markdown["content"]
		} else if message.File != nil {
			fileName := message.File["media_id"]
			content = fmt.Sprintf("[%v](%v)", fileName, fileName)
		} else {
			fmt.Println(r)
			return
		}
		urlReg, _ := regexp.Compile(`http[\w:\?\.\=\#%&/]+`)
		content = urlReg.ReplaceAllString(content, "")

		byteContent := []byte(content)
		if len(byteContent) > 1024*10 {
			byteContent = byteContent[:1024*10]
		}

		byteContent = markdown.ToHTML(byteContent, nil, nil)
		hub.broadcast <- byteContent
		messageCache.addFile("fileName", byteContent)
	}
}

//日志文件上传
func uploadMedia(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	b, _ := ioutil.ReadAll(file)
	fileName := fmt.Sprintf("%v/%v", TMP, header.Filename)
	fileErr := ioutil.WriteFile(fileName, b, 0600)
	log.Println(fileName, len(b), err, fileErr)
	res := map[string]string{
		"media_id": header.Filename,
	}
	json.NewEncoder(w).Encode(res)
}

// 返回日志文件
func logFile(w http.ResponseWriter, r *http.Request) error {
	fileName := fmt.Sprintf("%v%v", TMP, r.URL)
	reg, _ := regexp.Compile(".log$")
	if !reg.Match([]byte(fileName)) {
		return fmt.Errorf("not log file")
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("is dir")
	}

	return nil
}

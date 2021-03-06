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

type Message struct {
	Markdown map[string]string `json:"markdown"`
	File     map[string]string `json:"file"`
	MediaId  string            `json:"media_id"`
}

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
			log.Printf("content=> %v\n", content)
		} else {
			fmt.Println(r)
			return
		}
		urlReg, _ := regexp.Compile(`\n>\s*http[\w:\?\.\-\=\#%&/]+\n?`)
		content = urlReg.ReplaceAllString(content, "")

		urlReg, _ = regexp.Compile(`http[\w:\?\.\-\=\#%&/]+`)
		content = urlReg.ReplaceAllString(content, message.MediaId)

		log.Printf("no http content=> %v\n", content)

		byteContent := []byte(content)
		if len(byteContent) > 1024*10 {
			byteContent = byteContent[:1024*10]
		}
		byteContent = markdown.ToHTML(byteContent, nil, nil)
		log.Printf("byteContent => %v\n", string(byteContent))
		id := messageCache.addFile("fileName", byteContent)
		idhtml := fmt.Sprintf("<span data-message-id=%v style='display:none'></span>", id)
		hub.broadcast <- append([]byte(idhtml), byteContent...)

	}
}

//日志文件上传
func uploadMedia(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	b, _ := ioutil.ReadAll(file)
	fileName := fmt.Sprintf("%v/%v", *logPath, header.Filename)
	fileErr := ioutil.WriteFile(fileName, b, 0600)
	log.Println(fileName, len(b), err, fileErr)
	res := map[string]string{
		"media_id": header.Filename,
	}
	json.NewEncoder(w).Encode(res)
}

// 返回日志文件
func logFile(w http.ResponseWriter, r *http.Request) error {
	fileName := fmt.Sprintf("%v%v", *logPath, r.URL)
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
	http.ServeFile(w, r, fileName)
	return nil
}

// 返回日志文件
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var ids []int64
	json.NewDecoder(r.Body).Decode(&ids)
	res, _ := messageCache.delFile(ids)
	w.Write([]byte(fmt.Sprint(res)))
}

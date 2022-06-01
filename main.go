// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8082", "http service address")
var logPath = flag.String("logPath", "/tmp/logs", "path save buil log file on")
var (
	gitHash   string
	buildTime string
	runMac    string
)
var messageCache Cache = Cache{}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	err := logFile(w, r)
	if err == nil {
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()

	if err := checkMac(); err != nil {
		log.Println(err)
		return
	}

	log.Println("gitHash : ", gitHash)
	log.Println("buildTime : ", buildTime)
	log.Println("logPath : ", *logPath)
	log.Println("addr : ", *addr)

	_, err := os.Stat(fmt.Sprintf("%v", *logPath))
	if err != nil {
		log.Printf("%v path error", *logPath)
		return
	}

	messageCache.JsonFile = fmt.Sprintf("%v/cache.db.json", *logPath)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/sendMessage", sendMessage(hub))
	http.HandleFunc("/uploadMedia", uploadMedia)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//校验mac地址
func checkMac() error {
	if len(runMac) == 0 {
		return nil
	}
	for _, mac := range getMacAddrs() {
		if mac == runMac {
			return nil
		}
	}
	return fmt.Errorf("run time error")
}

func getMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}
	// fmt.Println(netInterfaces)
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()

		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return
}

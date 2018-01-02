package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8888", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if r.URL.Path != "/" {
		log.Println(r.URL.Path + " Not found 404")
		http.Error(w, "Page not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}

	http.ServeFile(w, r, dir+"/home.html")
}

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("Listen and serve : ", err)
	}
}

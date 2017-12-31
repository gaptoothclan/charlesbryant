package main

import (
	"flag"
	"log"
  "net/http"
)

var addr = flag.String("addr", ":8888", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}

	http.ServeFile(w, r, "./home.html")
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
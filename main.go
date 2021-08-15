package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"text/template"
	"time"
)

var (
	reqs = 0
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute)
	return tc, nil
}

func ListenAndServe(addr string, num int) error {
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			log.Println(http.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}, nil))
			wg.Done()
		}(i)
	}
	log.Println("ALL THREAD STARTED")
	wg.Wait()
	return nil
}

func api(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"reqs":"%v"}`, reqs)
}

func index(w http.ResponseWriter, req *http.Request) {
	reqs++
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func dstat(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("dstat.html"))
	tmpl.Execute(w, nil)
}

func main() {
	//========SET==========
	threads := 1000
	port := "80"
	//=======END SET=======
	http.HandleFunc("/api", api)
	http.HandleFunc("/", index)
	http.HandleFunc("/dstat", dstat)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Println(ListenAndServe(":"+port, threads))
}

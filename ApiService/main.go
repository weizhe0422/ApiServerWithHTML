package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"time"
)

func handleDisplay(resp http.ResponseWriter, req *http.Request) {
	var (
		inputParams map[string][]string
		respString  *bytes.Buffer
	)

	inputParams = make(map[string][]string, 0)
	inputParams = req.URL.Query()

	respString = bytes.NewBuffer([]byte("Request String:\n"))

	for key, param := range inputParams {
		for _, value := range param {
			respString.Write([]byte(key + ":" + value + "\n"))
		}
	}
	resp.Write(respString.Bytes())
}


func main() {
	var (
		mux            *http.ServeMux
		httpSvr        *http.Server
		err            error
		listener       net.Listener
		httpDir        http.Dir
		webPageHandler http.Handler
		waitCh			chan string
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/test/show", handleDisplay)

	httpDir = http.Dir("./webPage")
	webPageHandler = http.FileServer(httpDir)
	mux.Handle("/", http.StripPrefix("/", webPageHandler))

	httpSvr = &http.Server{
		ReadTimeout:  5 * time.Millisecond,
		WriteTimeout: 5 * time.Millisecond,
		Handler:      mux,
	}

	if listener, err = net.Listen("tcp", ":8080"); err != nil {
		log.Fatalf("failed to new listener: %s", err)
		return
	}

	go httpSvr.Serve(listener)

	<- waitCh

}

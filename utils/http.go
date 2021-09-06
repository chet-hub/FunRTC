package utils

import (
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Msg struct {
	Request  chan string
	Response chan string
}

func HTTPServer() chan Msg {
	port := flag.Int("port", 8080, "http server port")
	flag.Parse()

	msgChan := make(chan Msg)
	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		body, _ := ioutil.ReadAll(r.Body)
		var msg = Msg{Request: make(chan string), Response: make(chan string)}
		msgChan <- msg
		msg.Request <- string(body)
		//fmt.Printf("request: %s\n", string(body))
		response := <-msg.Response
		w.Write([]byte(response))
		//fmt.Printf("response: %s\n", response)
		close(msg.Request)
		close(msg.Response)
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return msgChan
}


//func main() {
//	var m = HTTPServer()
//	for {
//		msg := <-m
//		re := <-msg.Request
//		fmt.Printf("request: %s\n", re)
//		//time.Sleep(1 * time.Second)
//		msg.Response <- re
//		fmt.Printf("response: %s\n", re)
//	}
//}
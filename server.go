package main

import (
	"fmt"
    "net/http"
    "log"
	"time"
	// "bytes"
	// "io"
)

func main(){

	http.HandleFunc("/watch", watchAPI)
	log.Printf("http server started at :%d...", 8897)
	http.ListenAndServe(fmt.Sprintf(":%d", 8897), nil)
}

func watchAPI(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	flusher, ok := resp.(http.Flusher)
	if !ok {
		msg := fmt.Sprintf("unable to start watch - can't get http.Flusher: %#v", resp)
		log.Println(msg)
		resp.Write([]byte(msg))
		return
	}
	cn, ok := resp.(http.CloseNotifier)
	if !ok {
		// 感知到客户端断开连接的行为的方式, 类似于一个 stop channel.
		// 否则一旦开始, 下面的 for{} 循环就无法停止(多个客户端同时连接时, 还会同时进行多个循环).
		msg := fmt.Sprintf("unable to start watch - can't get http.CloseNotifier: %#v", resp)
		log.Println(msg)
		resp.Write([]byte(msg))
		return
	}

	resp.Header().Set("Content-Type", "text/plain")
	// Transfer-Encoding 与 Content-Length 不能同时出现
	resp.Header().Set("Transfer-Encoding", "chunked")
	resp.WriteHeader(http.StatusOK)
	flusher.Flush()

	var ch chan string
	ch = make(chan string, 1)
	go func() {
		for {
			ch <- time.Now().Format("2006-01-02 15:04:05")
			time.Sleep(time.Second * 1)
		}
	}()
	for {
		select {
		case <-cn.CloseNotify():
			log.Printf("client %v disconnected from the server", req.RemoteAddr)
			return
		case info := <- ch:
			log.Printf("info: %s\n", info)

			resp.Write([]byte(info+"\n"))
			// 下面两种方法也可以
			// fmt.Fprintf(resp, "%s\n", info)
			// io.WriteString(resp, info+"\n")

			flusher.Flush()
		}
	}
}

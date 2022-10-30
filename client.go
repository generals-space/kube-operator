package main

import (
	"bufio"
	"log"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("http://127.0.0.1:8897/watch")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Printf("response: %T\n", resp)
	log.Printf("response.Body: %T\n", resp.Body)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			log.Printf("%s", line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	log.Println("done")
}

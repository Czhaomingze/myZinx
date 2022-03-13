package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start")

	time.Sleep(1)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error:", err)
		return
	}
	for {
		_, err := conn.Write([]byte("Hello zinx"))
		if err != nil {
			fmt.Println("write conn error:", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error:", err)
			return
		}

		fmt.Printf("server call back:%s,cnt =%d\n", buf, cnt)
		time.Sleep(10 * time.Second)
	}
}

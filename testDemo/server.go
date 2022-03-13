package main

import "zinx/znet"

func main() {
	//创建s服务器
	s := znet.NewServer("[Zinx V0.1]")

	//启动服务
	s.Serve()
}

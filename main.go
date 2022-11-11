package main

import (
	"nepal/proxy"
)

func main()  {
	go proxy.DefaultProxy.Run() //运行代理
	proxy.RunApp()              //启动UI
}

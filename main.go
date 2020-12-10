//go:generate go-bindata -prefix "./www/dist" -fs  ./www/dist/...
package main

import (
	"fmt"

	"myapp/libs"

	"myapp/web_server"
)

func main() {
	libs.InitConfig("/Users/mac/myapp/config.yml")

	irisServer := web_server.NewServer(nil)
	if irisServer == nil {
		panic("Http 初始化失败")
	}
	irisServer.NewApp()

	if libs.IsPortInUse(libs.Config.Port) {
		panic(fmt.Sprintf("端口 %d 已被使用", libs.Config.Port))
	}

	err := irisServer.Serve()
	if err != nil {
		panic(err)
	}

}

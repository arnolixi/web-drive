package main

import (
	"flag"
	"log"

	arc "gitee.com/Arno-lixi/web-drive"
	"gitee.com/Arno-lixi/web-drive/demo/conf"
	"gitee.com/Arno-lixi/web-drive/demo/src/controller"

	"github.com/gin-gonic/gin"
)

func main() {

	app := arc.Start()
	confFile := flag.String("config", "./conf/config.yaml", "指定配置文件")
	flag.Parse()
	app.LoadConf(*confFile, conf.NewConfig())

	app.GroupUse("/api/v2", H1()).
		GroupMount(
			controller.NewMyTestClass(),
		)
	app.GroupUse("/api/v1", H1(), H2())
	app.Mount("/api/v1", controller.NewMyTestClass())
	app.GoRun()

}

func H1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("H1 开始执行")
		ctx.Next()
	}

}

func H2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("H2 开始执行")
		ctx.Next()

	}

}

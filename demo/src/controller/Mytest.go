package controller

import (
	arc "web-drive"

	"github.com/gin-gonic/gin"
)

type MyTestClass struct {
	*arc.GormAdapter
}

func NewMyTestClass() *MyTestClass {
	return &MyTestClass{}
}

func (this *MyTestClass) Name() string {
	return "MyTestClass"
}

func (this *MyTestClass) test(ctx *gin.Context) arc.Json {
	return gin.H{"mmsg": "ok"}
}

func (this *MyTestClass) Build(arc *arc.Reactor) {
	arc.Handle("GET", "test", this.test)
}

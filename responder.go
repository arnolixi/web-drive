package arc

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var responderList []Responder

var once_resp_list sync.Once

func get_resp_list() []Responder {
	once_resp_list.Do(func() {
		responderList = []Responder{}
	})
	return responderList
}

type Responder interface {
	RespondTo() gin.HandlerFunc
}

type Json interface{}
type ResponderForJson func(ctx *gin.Context) Json

func (r ResponderForJson) RespondTo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ctx.String(200)
	}

}

package arc

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

var responderList []Responder

var once_resp_list sync.Once

func get_resp_list() []Responder {
	once_resp_list.Do(func() {
		responderList = []Responder{
			ResponderForJson(nil),
		}
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
		// ctx.String(200,r(ctx))
		ctx.JSON(200, r(ctx))
	}

}

func Convert(handler interface{}) gin.HandlerFunc {
	h_ref := reflect.ValueOf(handler)
	for _, resp := range get_resp_list() {
		r_ref := reflect.TypeOf(resp)
		if h_ref.Type().ConvertibleTo(r_ref) {
			return h_ref.Convert(r_ref).Interface().(Responder).RespondTo()
		}
	}
	return nil
}

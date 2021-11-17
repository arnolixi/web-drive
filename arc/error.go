package arc

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
)

var _ error = &Error{}

type ErrorType int
type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) SetType(flags ErrorType) *Error {
	e.Type = flags
	return e
}

func (e *Error) SetMeta(data interface{}) *Error {
	e.Meta = data
	return e
}

func (e *Error) JSON() interface{} {
	jsonData := gin.H{}
	if e.Meta != nil {
		value := reflect.ValueOf(e.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return e.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				jsonData[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			jsonData["meta"] = e.Meta
		}
	}
	if _, ok := jsonData["error"]; !ok {
		jsonData["error"] = e.Error()
	}
	return jsonData
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.JSON())
}

func (e *Error) IsType(flags ErrorType) bool {
	return e.Type == flags
}

func (e *Error) Unwrap() error {
	return e.Err
}

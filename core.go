package arc

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type Reactor struct {
	Engine       *gin.Engine
	g            *gin.RouterGroup
	currentGroup string
	*Config
}

func Start(ginMiddlewares ...gin.HandlerFunc) *Reactor {
	r := &Reactor{
		Engine: gin.New(),
	}
	// 使用错误处理
	r.Engine.Use(ErrorHandler())
	r.Engine.Use(ginMiddlewares...)

	return r
}

func (r *Reactor) Load() {
	r.Engine.Run(r.Config.Addr)

}

func (r *Reactor) Mount(group string, classes ...IClass) *Reactor {
	r.g = r.Engine.Group(group)
	for _, class := range classes {
		r.currentGroup = group
		class.Build(r)
		r.Beans(class)
	}

	return r
}

func (r *Reactor) Beans(beans ...Bean) *Reactor {

	return r

}
func (r *Reactor) InjectConfig(cfgs ...interface{}) *Reactor {
	BeanFactory.InjectConfig(cfgs...)
	return r
}

func (r *Reactor) applyAll() {
	for t, v := range BeanFactory.GetBM() {
		if t.Elem().Kind() == reflect.Struct {
			BeanFactory.Apply(v.Interface())
		}
	}
}

package arc

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Reactor struct {
	Engine *gin.Engine
	g      *gin.RouterGroup
	gMap   map[string]*gin.RouterGroup
	Config IConfig
}

func Start(ginMiddlewares ...gin.HandlerFunc) *Reactor {
	r := &Reactor{
		Engine: gin.New(),
		gMap:   make(map[string]*gin.RouterGroup),
	}
	// 使用错误处理
	r.Engine.Use(ErrorHandler())
	r.Engine.Use(ginMiddlewares...)

	return r
}

func (r *Reactor) LoadConf(filePath string, configPtr IConfig) *Reactor {
	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Failed to Viper ReadInConfig,Error:", err)
	}
	if err = viper.Unmarshal(configPtr); err != nil {
		fmt.Println("Viper Unmarshal Failed,Error:", err)
	}
	r.Config = configPtr
	return r
}

func (r *Reactor) GoRun() {
	r.applyAll()
	r.Engine.Run(r.Config.GetServer().GetAddr())
}

func (r *Reactor) Handle(httpMethod, path string, handler interface{}) *Reactor {
	if h := Convert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, method := range methods {
			r.g.Handle(method, path, h)
		}
	}
	return r
}

func (r *Reactor) GroupUse(group string, handlers ...gin.HandlerFunc) *Reactor {
	var ok bool
	r.g, ok = r.gMap[group]
	if !ok {
		r.g = r.Engine.Group(group)
		r.gMap[group] = r.g
	}
	r.g.Use(handlers...)
	return r
}

func (r *Reactor) Mount(group string, classes ...IClass) *Reactor {
	var ok bool
	r.g, ok = r.gMap[group]
	if !ok {
		r.g = r.Engine.Group(group)
		r.gMap[group] = r.g
	}

	for _, class := range classes {
		class.Build(r)
		r.Beans(class)
	}

	return r
}

func (r *Reactor) GroupMount(classes ...IClass) *Reactor {
	log.Println(r.g.BasePath())
	for _, class := range classes {
		class.Build(r)
		r.Beans(class)
	}
	return r
}

func (r *Reactor) Beans(beans ...Bean) *Reactor {
	for _, b := range beans {
		BeanFactory.Set(b)
	}
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

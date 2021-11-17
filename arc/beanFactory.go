package arc

import (
	"reflect"
)

type BeanFactoryImpl struct {
	beanMapper BeanMapper
}

func NewBeanFactoryImpl() *BeanFactoryImpl {
	return &BeanFactoryImpl{
		beanMapper: make(BeanMapper),
	}
}

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactoryImpl()
}

func (this *BeanFactoryImpl) Set(vList ...interface{}) {
	if vList == nil || len(vList) == 0 {
		return
	}
	for _, v := range vList {
		this.beanMapper.add(v)
	}
}

func (this *BeanFactoryImpl) Get(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	getVal := this.beanMapper.get(v)
	if getVal.IsValid() {
		return getVal.Interface()
	}
	return nil
}

func (this *BeanFactoryImpl) GetBM() BeanMapper {
	return this.beanMapper
}

func (this *BeanFactoryImpl) Apply(bean interface{}) {
	if bean == nil {
		return
	}
	v := reflect.ValueOf(bean)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if v.Field(i).CanSet() && field.Tag.Get("inject") != "" {
			if field.Tag.Get("inject") == "-" {
				// 注入指定 struct
				if getV := this.Get(field.Type); getV != nil {
					v.Field(i).Set(reflect.ValueOf(getV))
					this.Apply(getV)
				}
			} else {
				// TODO 支持表达式注入
				continue
			}
		}

	}

}

func (this *BeanFactoryImpl) InjectConfig(cfgs ...interface{}) {
	for _, cfg := range cfgs {
		t := reflect.TypeOf(cfgs)
		if t.Kind() != reflect.Ptr {
			panic("Required Ptr Object!")
		}
		if t.Elem().Kind() != reflect.Struct {
			continue
		}
		this.Set(cfg)
		this.Apply(cfg)
		v := reflect.ValueOf(cfg)
		for i := 0; i < t.NumMethod(); i++ {
			m := v.Method(i)
			callRet := m.Call(nil)
			if callRet != nil && len(callRet) == 1 {
				this.Set(callRet[0].Interface())
			}
		}
	}

}

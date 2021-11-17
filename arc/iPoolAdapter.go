package arc

import (
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type IPool struct {
	*ants.Pool
}

func (this *IPool) Name() string {
	return "IPool"
}

func NewIPool() *IPool {
	pool, err := ants.NewPool(100)
	if err != nil {
		zap.L().Error("Init Goroutine pool Failed.", zap.Error(err))
		return nil
	}
	return &IPool{
		pool,
	}
}

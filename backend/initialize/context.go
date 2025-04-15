package initialize

import (
	"context"
	"sync"
)

type AppContext struct {
	ctx  context.Context
	lock sync.RWMutex
}

func InitDefaultContext() *AppContext {
	return &AppContext{
		ctx: context.Background(),
	}
}

func (c *AppContext) SetContext(ctx context.Context) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ctx = ctx
}

func (c *AppContext) Context() context.Context {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.ctx
}

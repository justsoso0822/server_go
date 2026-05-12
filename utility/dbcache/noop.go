package dbcache

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
)

// NoopAdapter 是一个空操作缓存适配器，所有读写都不生效。
// 设置此适配器后，业务代码的 .Cache() 链式调用变成 no-op。
type NoopAdapter struct{}

func (n *NoopAdapter) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	return nil
}

func (n *NoopAdapter) SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	return nil
}

func (n *NoopAdapter) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error) {
	return true, nil
}

func (n *NoopAdapter) SetIfNotExistFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (bool, error) {
	return true, nil
}

func (n *NoopAdapter) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (bool, error) {
	return true, nil
}

func (n *NoopAdapter) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	return nil, nil
}

func (n *NoopAdapter) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (*gvar.Var, error) {
	return gvar.New(value), nil
}

func (n *NoopAdapter) GetOrSetFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (*gvar.Var, error) {
	v, err := f(ctx)
	if err != nil {
		return nil, err
	}
	return gvar.New(v), nil
}

func (n *NoopAdapter) GetOrSetFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (*gvar.Var, error) {
	v, err := f(ctx)
	if err != nil {
		return nil, err
	}
	return gvar.New(v), nil
}

func (n *NoopAdapter) Contains(ctx context.Context, key interface{}) (bool, error) {
	return false, nil
}

func (n *NoopAdapter) Size(ctx context.Context) (int, error) {
	return 0, nil
}

func (n *NoopAdapter) Data(ctx context.Context) (map[interface{}]interface{}, error) {
	return nil, nil
}

func (n *NoopAdapter) Keys(ctx context.Context) ([]interface{}, error) {
	return nil, nil
}

func (n *NoopAdapter) Values(ctx context.Context) ([]interface{}, error) {
	return nil, nil
}

func (n *NoopAdapter) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	return nil, false, nil
}

func (n *NoopAdapter) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return 0, nil
}

func (n *NoopAdapter) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	return 0, nil
}

func (n *NoopAdapter) Remove(ctx context.Context, keys ...interface{}) (*gvar.Var, error) {
	return nil, nil
}

func (n *NoopAdapter) Clear(ctx context.Context) error {
	return nil
}

func (n *NoopAdapter) Close(ctx context.Context) error {
	return nil
}
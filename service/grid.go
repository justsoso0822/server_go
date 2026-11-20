package service

import (
	"context"
	"fmt"
	"sync"
)

func GetGrid(ctx context.Context, uid int64, chapter int) (map[string]any, error) {
	out := map[string]any{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	collect := func(key string, fn func() (any, error)) {
		wg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = fmt.Errorf("%s panic: %v", key, r)
					}
					mu.Unlock()
				}
				wg.Done()
			}()
			result, err := fn()
			mu.Lock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
			} else {
				out[key] = result
			}
			mu.Unlock()
		}()
	}

	collect("bag", func() (any, error) { return GetUserBag(ctx, uid, chapter) })
	collect("bag_tp", func() (any, error) { return GetUserBagTp(ctx, uid, chapter) })
	collect("tasks", func() (any, error) { return InitTasks(ctx, uid) })

	wg.Wait()
	return out, firstErr
}

package grid

import (
	"context"
	"fmt"
	"sync"

	"server_go/internal/service/bag"
	"server_go/internal/service/task"

	"github.com/gogf/gf/v2/frame/g"
)

func GetGrid(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	out := g.Map{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	wg.Add(3)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("GetUserBag panic: %v", r)
				}
				mu.Unlock()
			}
			wg.Done()
		}()
		result, err := bag.GetUserBag(ctx, uid, chapter)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out["bag"] = result
		mu.Unlock()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("GetUserBagTp panic: %v", r)
				}
				mu.Unlock()
			}
			wg.Done()
		}()
		result, err := bag.GetUserBagTp(ctx, uid, chapter)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out["bag_tp"] = result
		mu.Unlock()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("InitTasks panic: %v", r)
				}
				mu.Unlock()
			}
			wg.Done()
		}()
		tasks, err := task.InitTasks(ctx, uid)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out["tasks"] = tasks
		mu.Unlock()
	}()

	wg.Wait()
	return out, firstErr
}

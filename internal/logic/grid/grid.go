package grid

import (
	"context"
	"sync"

	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sGrid struct{}

func init() {
	service.RegisterGrid(&sGrid{})
}

func (s *sGrid) GetGrid(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	out := g.Map{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	wg.Add(3)

	go func() {
		defer wg.Done()
		result, err := service.Bag().GetUserBag(ctx, uid, chapter)
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
		defer wg.Done()
		result, err := service.Bag().GetUserBagTp(ctx, uid, chapter)
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
		defer wg.Done()
		tasks, err := service.Task().InitTasks(ctx, uid)
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

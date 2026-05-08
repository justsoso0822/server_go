package grid

import (
	"context"
	"sync"

	"server_go/internal/logic/bag"
	"server_go/internal/logic/task"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sGrid struct{}

func init() {
	service.RegisterGrid(&sGrid{})
}

func (s *sGrid) GetGrid(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	return GetGrid(ctx, uid, chapter)
}

// GetGrid fetches bag + bag_tp + tasks concurrently.
func GetGrid(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	ret := g.Map{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	wg.Add(3)

	go func() {
		defer wg.Done()
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
		ret["bag"] = result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
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
		ret["bag_tp"] = result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
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
		ret["tasks"] = tasks
		mu.Unlock()
	}()

	wg.Wait()
	return ret, firstErr
}
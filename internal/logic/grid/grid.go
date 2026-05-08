package grid

import (
	"context"
	"sync"

	"server_go/internal/model"
	"server_go/internal/service"
)

type sGrid struct{}

func init() {
	service.RegisterGrid(&sGrid{})
}

func (s *sGrid) GetGrid(ctx context.Context, in *model.BagInput) (*model.GridOutput, error) {
	out := &model.GridOutput{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	wg.Add(3)

	go func() {
		defer wg.Done()
		result, err := service.Bag().GetUserBag(ctx, in)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out.Bag = result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		result, err := service.Bag().GetUserBagTp(ctx, in)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out.BagTp = result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		tasks, err := service.Task().InitTasks(ctx, in.Uid)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out.Tasks = tasks
		mu.Unlock()
	}()

	wg.Wait()
	return out, firstErr
}
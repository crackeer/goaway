package container

import (
	"context"
	"sync"
	"time"
)

// StartSchedule
//
//	@param ctx
//	@param cfg
func StartSchedule(ctx context.Context, wg *sync.WaitGroup, cfg *AppConfig) {
	if cfg.SyncInterval > 0 {
		syncRouter(ctx, wg, cfg.SyncInterval)
	}
}

func syncRouter(ctx context.Context, wg *sync.WaitGroup, interval int64) {
	wg.Add(1)
	defer wg.Done()

loop:
	for {
		ticker := time.NewTicker(time.Duration(interval) * time.Duration(time.Second))
		select {
		case <-ctx.Done():
			Log(map[string]interface{}{
				"message": "Recieved Stop Signal",
			}, "StopSyncRouter")
			break loop
		case <-ticker.C:
			ticker.Stop()
			Log(map[string]interface{}{}, "SyncAPIOnce")
		}
	}
	Log(map[string]interface{}{
		"interval": interval,
	}, "SyncRouterStop")
}

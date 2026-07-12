package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/robfig/cron/v3"
)

const (
	upstreamRateSyncSchedule  = "0 2,10,18 * * *"
	upstreamRateSyncTimeout   = 20 * time.Minute
	upstreamRateSyncComponent = "upstream.cost_rate_sync"
)

var upstreamRateSyncCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

type UpstreamRateSyncRunner struct {
	syncService    *UpstreamRateSyncService
	settingService *SettingService

	cron      *cron.Cron
	startOnce sync.Once
	stopOnce  sync.Once
	runMu     sync.Mutex
}

func NewUpstreamRateSyncRunner(syncService *UpstreamRateSyncService, settingService *SettingService) *UpstreamRateSyncRunner {
	return &UpstreamRateSyncRunner{syncService: syncService, settingService: settingService}
}

func (r *UpstreamRateSyncRunner) Start() {
	if r == nil || r.syncService == nil || r.settingService == nil {
		return
	}
	r.startOnce.Do(func() {
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			location = time.FixedZone("Asia/Shanghai", 8*60*60)
		}
		scheduler := cron.New(cron.WithParser(upstreamRateSyncCronParser), cron.WithLocation(location))
		if _, err := scheduler.AddFunc(upstreamRateSyncSchedule, r.runScheduled); err != nil {
			logger.WriteSinkEvent("error", upstreamRateSyncComponent, "自动同步成本倍率调度器启动失败", map[string]any{"error": err.Error()})
			return
		}
		r.cron = scheduler
		r.cron.Start()
	})
}

func (r *UpstreamRateSyncRunner) Stop() {
	if r == nil {
		return
	}
	r.stopOnce.Do(func() {
		if r.cron == nil {
			return
		}
		ctx := r.cron.Stop()
		select {
		case <-ctx.Done():
		case <-time.After(3 * time.Second):
		}
	})
}

func (r *UpstreamRateSyncRunner) runScheduled() {
	if r == nil || r.syncService == nil || r.settingService == nil {
		return
	}
	if !r.settingService.IsUpstreamRateSyncEnabled(context.Background()) {
		return
	}
	if !r.runMu.TryLock() {
		logger.WriteSinkEvent("warn", upstreamRateSyncComponent, "上一次自动同步成本倍率仍在运行，本次已跳过", map[string]any{"schedule": upstreamRateSyncSchedule})
		return
	}
	defer r.runMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), upstreamRateSyncTimeout)
	defer cancel()
	result, err := r.syncService.SyncAllAPIKeyAccounts(ctx)
	if err != nil {
		logger.WriteSinkEvent("error", upstreamRateSyncComponent, "自动同步成本倍率失败", map[string]any{
			"schedule": upstreamRateSyncSchedule,
			"error":    safeUpstreamRateSyncError(err),
		})
		return
	}

	for _, item := range result.Items {
		if item.Error == nil {
			continue
		}
		logger.WriteSinkEvent("error", upstreamRateSyncComponent, "账号自动同步成本倍率失败", map[string]any{
			"account_id":   item.AccountID,
			"account_name": item.AccountName,
			"schedule":     upstreamRateSyncSchedule,
			"error":        safeUpstreamRateSyncError(item.Error),
		})
	}

	level := "info"
	message := "自动同步成本倍率完成"
	if result.Failed > 0 {
		level = "warn"
		message = "自动同步成本倍率部分失败"
	}
	logger.WriteSinkEvent(level, upstreamRateSyncComponent, message, map[string]any{
		"schedule": upstreamRateSyncSchedule,
		"total":    result.Total,
		"success":  result.Success,
		"failed":   result.Failed,
	})
}

func safeUpstreamRateSyncError(err error) string {
	if err == nil {
		return ""
	}
	var syncErr *UpstreamAccountMetaSyncError
	if errors.As(err, &syncErr) {
		return syncErr.SafeMessage()
	}
	return err.Error()
}

func ProvideUpstreamRateSyncRunner(syncService *UpstreamRateSyncService, settingService *SettingService) *UpstreamRateSyncRunner {
	runner := NewUpstreamRateSyncRunner(syncService, settingService)
	runner.Start()
	return runner
}

package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"golang.org/x/sync/errgroup"
)

const upstreamRateSyncConcurrency = 3

type UpstreamRateSyncResult struct {
	Account                *Account
	Meta                   *UpstreamAccountMeta
	PreviousRateMultiplier float64
	RateMultiplier         float64
	Changed                bool
	SignificantChange      bool
}

type UpstreamRateSyncBatchItem struct {
	AccountID   int64
	AccountName string
	Result      *UpstreamRateSyncResult
	Error       error
}

type UpstreamRateSyncBatchResult struct {
	Total   int
	Success int
	Failed  int
	Items   []UpstreamRateSyncBatchItem
}

type upstreamRateSyncAdmin interface {
	ListAccounts(ctx context.Context, page, pageSize int, platform, accountType, status, anomalyReason, search string, groupID int64, privacyMode string, sortBy, sortOrder string) ([]Account, int64, error)
	GetAccountsByIDs(ctx context.Context, ids []int64) ([]*Account, error)
	UpdateAccount(ctx context.Context, id int64, input *UpdateAccountInput) (*Account, error)
}

type UpstreamRateSyncService struct {
	adminService       upstreamRateSyncAdmin
	accountTestService *AccountTestService
}

func NewUpstreamRateSyncService(adminService upstreamRateSyncAdmin, accountTestService *AccountTestService) *UpstreamRateSyncService {
	return &UpstreamRateSyncService{adminService: adminService, accountTestService: accountTestService}
}

func ProvideUpstreamRateSyncService(adminService AdminService, accountTestService *AccountTestService) *UpstreamRateSyncService {
	return NewUpstreamRateSyncService(adminService, accountTestService)
}

func (s *UpstreamRateSyncService) SyncAccount(ctx context.Context, account *Account) (*UpstreamRateSyncResult, error) {
	if s == nil || s.adminService == nil || s.accountTestService == nil {
		return nil, errors.New("upstream rate sync service is not configured")
	}
	if account == nil {
		return nil, errors.New("account is required")
	}

	previousRateMultiplier := account.BillingRateMultiplier()
	meta, err := s.accountTestService.FetchCompatibleUpstreamAccountMeta(ctx, account)
	if err != nil {
		return nil, err
	}
	updated, err := s.adminService.UpdateAccount(ctx, account.ID, &UpdateAccountInput{RateMultiplier: &meta.RateMultiplier})
	if err != nil {
		return nil, err
	}

	changed := math.Abs(meta.RateMultiplier-previousRateMultiplier) > 1e-9
	significantChange := IsSignificantRateMultiplierChange(previousRateMultiplier, meta.RateMultiplier)
	if changed {
		changeRatio := rateMultiplierChangeRatio(previousRateMultiplier, meta.RateMultiplier)
		level := "info"
		if significantChange {
			level = "warn"
		}
		logger.WriteSinkEvent(level, "upstream.cost_rate_sync", "账号上游成本倍率已同步", map[string]any{
			"account_id":               account.ID,
			"account_name":             account.Name,
			"platform":                 account.Platform,
			"previous_rate_multiplier": previousRateMultiplier,
			"rate_multiplier":          meta.RateMultiplier,
			"change_ratio":             changeRatio,
			"significant_change":       significantChange,
			"rate_source":              meta.RateSource,
		})
		if significantChange {
			logger.WriteSinkEvent("warn", "upstream.health_alert", "账号上游成本倍率发生大幅变化", map[string]any{
				"alert_key":                fmt.Sprintf("cost_multiplier_jump:account:%d", account.ID),
				"alert_type":               "cost_multiplier_jump",
				"alert_status":             "firing",
				"severity":                 "warning",
				"account_id":               account.ID,
				"account_name":             account.Name,
				"platform":                 account.Platform,
				"previous_rate_multiplier": previousRateMultiplier,
				"rate_multiplier":          meta.RateMultiplier,
				"change_ratio":             changeRatio,
				"rate_source":              meta.RateSource,
			})
		}
	}

	return &UpstreamRateSyncResult{
		Account:                updated,
		Meta:                   meta,
		PreviousRateMultiplier: previousRateMultiplier,
		RateMultiplier:         meta.RateMultiplier,
		Changed:                changed,
		SignificantChange:      significantChange,
	}, nil
}

func (s *UpstreamRateSyncService) SyncAllAPIKeyAccounts(ctx context.Context) (*UpstreamRateSyncBatchResult, error) {
	if s == nil || s.adminService == nil {
		return nil, errors.New("upstream rate sync service is not configured")
	}
	const pageSize = 30
	accountIDs := make([]int64, 0, pageSize)
	for page := 1; ; page++ {
		accounts, total, err := s.adminService.ListAccounts(ctx, page, pageSize, "", AccountTypeAPIKey, "", "", "", 0, "", "id", "asc")
		if err != nil {
			return nil, err
		}
		for index := range accounts {
			accountIDs = append(accountIDs, accounts[index].ID)
		}
		if len(accountIDs) >= int(total) || len(accounts) == 0 {
			break
		}
	}
	return s.SyncAccounts(ctx, accountIDs)
}

func (s *UpstreamRateSyncService) SyncAccounts(ctx context.Context, accountIDs []int64) (*UpstreamRateSyncBatchResult, error) {
	accountIDs = normalizePositiveInt64IDs(accountIDs)
	if len(accountIDs) == 0 {
		return &UpstreamRateSyncBatchResult{Items: []UpstreamRateSyncBatchItem{}}, nil
	}
	accounts, err := s.adminService.GetAccountsByIDs(ctx, accountIDs)
	if err != nil {
		return nil, err
	}
	byID := make(map[int64]*Account, len(accounts))
	for _, account := range accounts {
		if account != nil {
			byID[account.ID] = account
		}
	}

	items := make([]UpstreamRateSyncBatchItem, len(accountIDs))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(upstreamRateSyncConcurrency)
	for index, accountID := range accountIDs {
		index, accountID := index, accountID
		g.Go(func() error {
			account := byID[accountID]
			if account == nil {
				items[index] = UpstreamRateSyncBatchItem{AccountID: accountID, Error: errors.New("account not found")}
				return nil
			}
			item := UpstreamRateSyncBatchItem{AccountID: accountID, AccountName: account.Name}
			if account.Type != AccountTypeAPIKey {
				item.Error = errors.New("upstream rate sync only supports API-key accounts")
				items[index] = item
				return nil
			}
			item.Result, item.Error = s.SyncAccount(gctx, account)
			items[index] = item
			return nil
		})
	}
	_ = g.Wait()

	result := &UpstreamRateSyncBatchResult{Total: len(items), Items: items}
	for _, item := range items {
		if item.Error != nil {
			result.Failed++
		} else {
			result.Success++
		}
	}
	return result, nil
}

func IsSignificantRateMultiplierChange(previous, next float64) bool {
	if math.Abs(next-previous) <= 1e-9 {
		return false
	}
	if previous <= 0 {
		return next > 0
	}
	return math.Abs(next-previous)/previous >= 0.5-1e-9
}

func rateMultiplierChangeRatio(previous, next float64) float64 {
	if previous <= 0 {
		return 0
	}
	return (next - previous) / previous
}

func normalizePositiveInt64IDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

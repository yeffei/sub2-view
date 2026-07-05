package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strings"
	"time"
)

const poolHealthSlotDuration = 5 * time.Minute

type PoolHealthService struct {
	upstreamPoolRepo   UpstreamPoolRepository
	accountRepo        AccountRepository
	accountMonitorRepo AccountMonitorRepository
}

func NewPoolHealthService(
	upstreamPoolRepo UpstreamPoolRepository,
	accountRepo AccountRepository,
	accountMonitorRepo AccountMonitorRepository,
) *PoolHealthService {
	return &PoolHealthService{
		upstreamPoolRepo:   upstreamPoolRepo,
		accountRepo:        accountRepo,
		accountMonitorRepo: accountMonitorRepo,
	}
}

func (s *PoolHealthService) ListUserPoolHealth(ctx context.Context) ([]*PoolHealthView, error) {
	details, err := s.buildPoolHealth(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*PoolHealthView, 0, len(details))
	for _, detail := range details {
		out = append(out, poolHealthDetailToView(detail))
	}
	return out, nil
}

func (s *PoolHealthService) GetUserPoolHealthDetail(ctx context.Context, poolID int64) (*PoolHealthDetail, error) {
	if poolID <= 0 {
		return nil, ErrPoolHealthNotFound
	}
	details, err := s.buildPoolHealth(ctx)
	if err != nil {
		return nil, err
	}
	for _, detail := range details {
		if detail.ID == poolID {
			return detail, nil
		}
	}
	return nil, ErrPoolHealthNotFound
}

func poolHealthDetailToView(detail *PoolHealthDetail) *PoolHealthView {
	if detail == nil {
		return nil
	}
	return &PoolHealthView{
		ID:                  detail.ID,
		Name:                detail.Name,
		Provider:            detail.Provider,
		GroupID:             detail.GroupID,
		GroupName:           detail.GroupName,
		Status:              detail.Status,
		Availability7d:      detail.Availability7d,
		BestLatencyMs:       detail.BestLatencyMs,
		BestPingLatencyMs:   detail.BestPingLatencyMs,
		HealthyMemberCount:  detail.HealthyMemberCount,
		DegradedMemberCount: detail.DegradedMemberCount,
		FailedMemberCount:   detail.FailedMemberCount,
		TotalMemberCount:    detail.TotalMemberCount,
		Timeline:            detail.Timeline,
	}
}

func (s *PoolHealthService) buildPoolHealth(ctx context.Context) ([]*PoolHealthDetail, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil {
		return []*PoolHealthDetail{}, nil
	}

	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("list upstream pools: %w", err)
	}
	bindings, err := s.upstreamPoolRepo.ListUpstreamPoolBindings(ctx)
	if err != nil {
		return nil, fmt.Errorf("list upstream pool bindings: %w", err)
	}
	enabledPools := filterEnabledPools(pools)
	poolBindings := groupPoolBindings(bindings, enabledPools)
	poolMembers, accountsByID := s.loadPoolMembersAndAccounts(ctx, enabledPools)
	accountFacts, err := s.loadPoolAccountFacts(ctx, enabledPools, poolBindings, poolMembers)
	if err != nil {
		return nil, err
	}

	out := make([]*PoolHealthDetail, 0, len(enabledPools))
	for _, pool := range enabledPools {
		bindings := poolBindings[pool.ID]
		detail := buildPoolHealthDetail(pool, bindings, poolMembers[pool.ID], accountsByID, accountFacts)
		out = append(out, detail)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func filterEnabledPools(pools []UpstreamPool) []UpstreamPool {
	out := make([]UpstreamPool, 0, len(pools))
	for _, pool := range pools {
		if pool.Enabled {
			out = append(out, pool)
		}
	}
	return out
}

func groupPoolBindings(bindings []UpstreamPoolBinding, pools []UpstreamPool) map[int64][]UpstreamPoolBinding {
	poolIDs := make(map[int64]struct{}, len(pools))
	for _, pool := range pools {
		poolIDs[pool.ID] = struct{}{}
	}
	out := make(map[int64][]UpstreamPoolBinding, len(pools))
	for _, binding := range bindings {
		if !binding.Enabled {
			continue
		}
		if _, ok := poolIDs[binding.PoolID]; !ok {
			continue
		}
		out[binding.PoolID] = append(out[binding.PoolID], binding)
	}
	for poolID := range out {
		sort.SliceStable(out[poolID], func(i, j int) bool {
			if out[poolID][i].Priority != out[poolID][j].Priority {
				return out[poolID][i].Priority < out[poolID][j].Priority
			}
			return out[poolID][i].ID < out[poolID][j].ID
		})
	}
	return out
}

func (s *PoolHealthService) loadPoolMembersAndAccounts(ctx context.Context, pools []UpstreamPool) (map[int64][]UpstreamPoolMember, map[int64]*Account) {
	membersByPool := make(map[int64][]UpstreamPoolMember, len(pools))
	accountIDs := make([]int64, 0)
	seenAccounts := make(map[int64]struct{})
	for _, pool := range pools {
		members, err := s.upstreamPoolRepo.ListUpstreamPoolMembers(ctx, pool.ID)
		if err != nil {
			slog.Warn("pool_health: load pool members failed", "pool_id", pool.ID, "error", err)
			continue
		}
		membersByPool[pool.ID] = members
		for _, member := range members {
			if member.AccountID <= 0 {
				continue
			}
			if _, ok := seenAccounts[member.AccountID]; ok {
				continue
			}
			seenAccounts[member.AccountID] = struct{}{}
			accountIDs = append(accountIDs, member.AccountID)
		}
	}
	accountsByID := make(map[int64]*Account, len(accountIDs))
	if len(accountIDs) == 0 {
		return membersByPool, accountsByID
	}
	accounts, err := s.accountRepo.GetByIDs(ctx, accountIDs)
	if err != nil {
		slog.Warn("pool_health: load accounts failed", "error", err)
		return membersByPool, accountsByID
	}
	for _, account := range accounts {
		if account != nil {
			accountsByID[account.ID] = account
		}
	}
	return membersByPool, accountsByID
}

type poolAccountFacts struct {
	latestByAccount  map[int64][]*AccountMonitorLatest
	avail7ByAccount  map[int64][]*AccountMonitorAvailability
	avail15ByAccount map[int64][]*AccountMonitorAvailability
	avail30ByAccount map[int64][]*AccountMonitorAvailability
	historyByAccount map[int64][]*AccountMonitorHistoryEntry
}

func (s *PoolHealthService) loadPoolAccountFacts(
	ctx context.Context,
	pools []UpstreamPool,
	bindings map[int64][]UpstreamPoolBinding,
	members map[int64][]UpstreamPoolMember,
) (*poolAccountFacts, error) {
	facts := emptyPoolAccountFacts()
	if s == nil || s.accountMonitorRepo == nil {
		return facts, nil
	}
	accountIDs := collectPoolAccountIDs(pools, members)
	if len(accountIDs) == 0 {
		return facts, nil
	}
	var err error
	facts.latestByAccount, err = s.accountMonitorRepo.ListLatestForAccountIDs(ctx, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("load account monitor latest: %w", err)
	}
	if facts.avail7ByAccount, err = s.accountMonitorRepo.ComputeAvailabilityForAccounts(ctx, accountIDs, monitorAvailability7Days); err != nil {
		return nil, fmt.Errorf("load account 7d availability: %w", err)
	}
	if facts.avail15ByAccount, err = s.accountMonitorRepo.ComputeAvailabilityForAccounts(ctx, accountIDs, monitorAvailability15Days); err != nil {
		return nil, fmt.Errorf("load account 15d availability: %w", err)
	}
	if facts.avail30ByAccount, err = s.accountMonitorRepo.ComputeAvailabilityForAccounts(ctx, accountIDs, monitorAvailability30Days); err != nil {
		return nil, fmt.Errorf("load account 30d availability: %w", err)
	}
	since := time.Now().UTC().AddDate(0, 0, -monitorAvailability30Days)
	facts.historyByAccount, err = s.accountMonitorRepo.ListHistorySinceForAccounts(ctx, accountIDs, nil, since)
	if err != nil {
		return nil, fmt.Errorf("load account monitor history window: %w", err)
	}
	return facts, nil
}

func emptyPoolAccountFacts() *poolAccountFacts {
	return &poolAccountFacts{
		latestByAccount:  map[int64][]*AccountMonitorLatest{},
		avail7ByAccount:  map[int64][]*AccountMonitorAvailability{},
		avail15ByAccount: map[int64][]*AccountMonitorAvailability{},
		avail30ByAccount: map[int64][]*AccountMonitorAvailability{},
		historyByAccount: map[int64][]*AccountMonitorHistoryEntry{},
	}
}

func collectPoolAccountIDs(
	pools []UpstreamPool,
	members map[int64][]UpstreamPoolMember,
) []int64 {
	ids := make([]int64, 0)
	seen := make(map[int64]struct{})
	for _, pool := range pools {
		for _, member := range members[pool.ID] {
			if member.AccountID <= 0 {
				continue
			}
			if _, ok := seen[member.AccountID]; !ok {
				seen[member.AccountID] = struct{}{}
				ids = append(ids, member.AccountID)
			}
		}
	}
	return ids
}

func buildPoolHealthDetail(
	pool UpstreamPool,
	bindings []UpstreamPoolBinding,
	members []UpstreamPoolMember,
	accountsByID map[int64]*Account,
	facts *poolAccountFacts,
) *PoolHealthDetail {
	groupID, groupName := primaryPoolHealthGroup(bindings)
	probeModel := probeModelForPool(pool.Platform, bindings)
	probeStatus, bestLatency, bestPing := aggregatePoolAccountLatest(members, probeModel, facts.latestByAccount)
	timeline := aggregatePoolHealthTimeline(members, probeModel, facts.historyByAccount)
	memberItems, healthy, degraded, failed := buildPoolMemberHealthItems(members, accountsByID)
	memberAvailability := aggregatePoolMemberAvailability(len(memberItems), healthy, degraded)
	availability7d := aggregatePoolWindowAvailability(members, probeModel, facts.historyByAccount, monitorAvailability7Days, memberAvailability)
	availability15d := aggregatePoolWindowAvailability(members, probeModel, facts.historyByAccount, monitorAvailability15Days, memberAvailability)
	availability30d := aggregatePoolWindowAvailability(members, probeModel, facts.historyByAccount, monitorAvailability30Days, memberAvailability)
	memberStatus := aggregatePoolMemberServingStatus(len(memberItems), healthy, degraded, failed)
	status := combinePoolMemberAndProbeStatus(memberStatus, probeStatus)
	if len(bindings) == 0 {
		status = MonitorStatusError
	}

	return &PoolHealthDetail{
		ID:                  pool.ID,
		Name:                pool.Name,
		Provider:            pool.Platform,
		GroupID:             groupID,
		GroupName:           groupName,
		Status:              status,
		Availability7d:      availability7d,
		Availability15d:     availability15d,
		Availability30d:     availability30d,
		BestLatencyMs:       bestLatency,
		BestPingLatencyMs:   bestPing,
		HealthyMemberCount:  healthy,
		DegradedMemberCount: degraded,
		FailedMemberCount:   failed,
		TotalMemberCount:    len(memberItems),
		Timeline:            timeline,
		Members:             memberItems,
		Lines:               buildPoolHealthLines(bindings, members, accountsByID, probeModel, facts),
	}
}

func primaryPoolHealthGroup(bindings []UpstreamPoolBinding) (*int64, string) {
	if len(bindings) == 0 {
		return nil, ""
	}
	id := bindings[0].GroupID
	return &id, bindings[0].GroupName
}

func aggregatePoolAccountLatest(members []UpstreamPoolMember, model string, latestByAccount map[int64][]*AccountMonitorLatest) (string, *int, *int) {
	status := ""
	var bestLatency, bestPing *int
	for _, member := range members {
		latest := pickAccountLatest(latestByAccount[member.AccountID], model)
		if latest == nil {
			continue
		}
		status = betterMonitorStatus(status, latest.Status)
		if isRoutableMonitorStatus(latest.Status) {
			bestLatency = minIntPtr(bestLatency, latest.LatencyMs)
			bestPing = minIntPtr(bestPing, latest.PingLatencyMs)
		}
	}
	return status, bestLatency, bestPing
}

func buildPoolMemberHealthItems(
	members []UpstreamPoolMember,
	accountsByID map[int64]*Account,
) ([]PoolHealthMemberItem, int, int, int) {
	out := make([]PoolHealthMemberItem, 0, len(members))
	healthy, degraded, failed := 0, 0, 0
	for _, member := range members {
		account := accountsByID[member.AccountID]
		healthStatus := resolvePoolMemberHealthStatus(member, account)
		switch healthStatus {
		case MonitorStatusOperational:
			healthy++
		case MonitorStatusDegraded:
			degraded++
		default:
			failed++
		}
		out = append(out, PoolHealthMemberItem{
			AccountID:     member.AccountID,
			AccountName:   firstNonEmptyTrimmed(member.AccountName, accountName(account)),
			Platform:      firstNonEmptyTrimmed(member.AccountPlatform, accountPlatform(account)),
			AccountStatus: firstNonEmptyTrimmed(accountStatus(account), member.AccountStatus),
			HealthStatus:  healthStatus,
			RuntimeStatus: member.RuntimeStatus,
			RuntimeReason: member.RuntimeReason,
			Schedulable:   poolMemberSchedulable(member, account),
			Enabled:       member.Enabled,
			ManualDrained: member.ManualDrained,
			Weight:        member.Weight,
			SourceType:    member.SourceType,
			SourceSetName: member.SourceSetName,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].HealthStatus != out[j].HealthStatus {
			return monitorStatusRank(out[i].HealthStatus) < monitorStatusRank(out[j].HealthStatus)
		}
		return out[i].AccountID < out[j].AccountID
	})
	return out, healthy, degraded, failed
}

func resolvePoolMemberHealthStatus(member UpstreamPoolMember, account *Account) string {
	if !member.Enabled || member.ManualDrained || (member.SchedulableOverride != nil && !*member.SchedulableOverride) {
		return MonitorStatusFailed
	}
	if account == nil || !account.IsActive() {
		return MonitorStatusFailed
	}
	if !account.IsSchedulable() {
		return MonitorStatusDegraded
	}
	return MonitorStatusOperational
}

func poolMemberSchedulable(member UpstreamPoolMember, account *Account) bool {
	if !member.Enabled || member.ManualDrained || (member.SchedulableOverride != nil && !*member.SchedulableOverride) {
		return false
	}
	return account != nil && account.IsSchedulable()
}

func aggregatePoolHealthStatus(total, healthy, degraded, failed int) string {
	if total == 0 {
		return MonitorStatusError
	}
	if healthy == total {
		return MonitorStatusOperational
	}
	if healthy > 0 || degraded > 0 {
		return MonitorStatusDegraded
	}
	if failed > 0 {
		return MonitorStatusFailed
	}
	return MonitorStatusError
}

func aggregatePoolMemberServingStatus(total, healthy, degraded, failed int) string {
	if total == 0 {
		return MonitorStatusError
	}
	if healthy > 0 {
		return MonitorStatusOperational
	}
	if degraded > 0 {
		return MonitorStatusDegraded
	}
	if failed > 0 {
		return MonitorStatusFailed
	}
	return MonitorStatusError
}

func combinePoolMemberAndProbeStatus(memberStatus, probeStatus string) string {
	if probeStatus == "" {
		return memberStatus
	}
	if memberStatus == MonitorStatusError || memberStatus == MonitorStatusFailed {
		return memberStatus
	}
	if memberStatus == MonitorStatusOperational {
		return MonitorStatusOperational
	}
	if isRoutableMonitorStatus(probeStatus) && isRoutableMonitorStatus(memberStatus) {
		return MonitorStatusOperational
	}
	return MonitorStatusDegraded
}

func buildPoolHealthLines(
	bindings []UpstreamPoolBinding,
	members []UpstreamPoolMember,
	accountsByID map[int64]*Account,
	probeModel string,
	facts *poolAccountFacts,
) []PoolHealthLineItem {
	groupID, groupName := primaryPoolHealthGroup(bindings)
	out := make([]PoolHealthLineItem, 0, len(members))
	for _, member := range members {
		account := accountsByID[member.AccountID]
		latest := pickAccountLatest(facts.latestByAccount[member.AccountID], probeModel)
		line := PoolHealthLineItem{
			AccountID:   member.AccountID,
			AccountName: firstNonEmptyTrimmed(member.AccountName, accountName(account)),
			GroupID:     groupID,
			GroupName:   groupName,
			ProbeModel:  probeModel,
		}
		if latest != nil {
			line.LatestStatus = latest.Status
			line.LatestLatencyMs = latest.LatencyMs
			line.LatestPingLatencyMs = latest.PingLatencyMs
			line.LastCheckedAt = &latest.CheckedAt
		}
		line.Availability7d = accountAvailabilityForModel(facts.avail7ByAccount[member.AccountID], probeModel)
		line.Availability15d = accountAvailabilityForModel(facts.avail15ByAccount[member.AccountID], probeModel)
		line.Availability30d = accountAvailabilityForModel(facts.avail30ByAccount[member.AccountID], probeModel)
		out = append(out, line)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].AccountID < out[j].AccountID
	})
	return out
}

func accountAvailabilityForModel(rows []*AccountMonitorAvailability, model string) float64 {
	for _, row := range rows {
		if row.Model == model {
			return row.AvailabilityPct
		}
	}
	return 0
}

func pickAccountLatest(rows []*AccountMonitorLatest, model string) *AccountMonitorLatest {
	if model == "" {
		return nil
	}
	for _, row := range rows {
		if row.Model == model {
			return row
		}
	}
	return nil
}

type poolTimelineSlot struct {
	status        string
	latencyMs     *int
	pingLatencyMs *int
	checkedAt     time.Time
}

func aggregatePoolWindowAvailability(
	members []UpstreamPoolMember,
	model string,
	historyByAccount map[int64][]*AccountMonitorHistoryEntry,
	windowDays int,
	fallback float64,
) float64 {
	if windowDays <= 0 {
		return 0
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -windowDays)
	entries := collectPoolHistoryEntries(members, model, historyByAccount, cutoff)
	if len(entries) == 0 {
		return fallback
	}
	ok := 0
	for _, entry := range entries {
		if isRoutableMonitorStatus(entry.Status) {
			ok++
		}
	}
	return math.Round(float64(ok)*10000/float64(len(entries))) / 100
}

func aggregatePoolMemberAvailability(total, healthy, degraded int) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(healthy)*10000/float64(total)) / 100
}

func aggregatePoolHealthTimeline(
	members []UpstreamPoolMember,
	model string,
	historyByAccount map[int64][]*AccountMonitorHistoryEntry,
) []PoolHealthTimelinePoint {
	slots := collectPoolHistorySlots(members, model, historyByAccount, time.Time{})
	if len(slots) == 0 {
		return nil
	}
	items := make([]poolTimelineSlot, 0, len(slots))
	for _, slot := range slots {
		items = append(items, slot)
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].checkedAt.After(items[j].checkedAt)
	})
	if len(items) > monitorTimelineMaxPoints {
		items = items[:monitorTimelineMaxPoints]
	}
	out := make([]PoolHealthTimelinePoint, 0, len(items))
	for _, slot := range items {
		out = append(out, PoolHealthTimelinePoint{
			Status:        slot.status,
			LatencyMs:     slot.latencyMs,
			PingLatencyMs: slot.pingLatencyMs,
			CheckedAt:     slot.checkedAt,
		})
	}
	return out
}

func collectPoolHistorySlots(
	members []UpstreamPoolMember,
	model string,
	historyByAccount map[int64][]*AccountMonitorHistoryEntry,
	cutoff time.Time,
) map[int64]poolTimelineSlot {
	out := make(map[int64]poolTimelineSlot)
	for _, entry := range collectPoolHistoryEntries(members, model, historyByAccount, cutoff) {
		checkedAt := entry.CheckedAt.UTC()
		key := checkedAt.Unix() / int64(poolHealthSlotDuration.Seconds())
		slot := out[key]
		slot.status = betterMonitorStatus(slot.status, entry.Status)
		if checkedAt.After(slot.checkedAt) {
			slot.checkedAt = checkedAt
		}
		if isRoutableMonitorStatus(entry.Status) {
			slot.latencyMs = minIntPtr(slot.latencyMs, entry.LatencyMs)
			slot.pingLatencyMs = minIntPtr(slot.pingLatencyMs, entry.PingLatencyMs)
		}
		out[key] = slot
	}
	return out
}

func collectPoolHistoryEntries(
	members []UpstreamPoolMember,
	model string,
	historyByAccount map[int64][]*AccountMonitorHistoryEntry,
	cutoff time.Time,
) []*AccountMonitorHistoryEntry {
	out := make([]*AccountMonitorHistoryEntry, 0)
	for _, member := range members {
		for _, entry := range historyByAccount[member.AccountID] {
			if entry == nil {
				continue
			}
			if model != "" && entry.Model != model {
				continue
			}
			checkedAt := entry.CheckedAt.UTC()
			if !cutoff.IsZero() && checkedAt.Before(cutoff) {
				continue
			}
			out = append(out, entry)
		}
	}
	return out
}

func firstNonEmptyTrimmed(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func accountName(account *Account) string {
	if account == nil {
		return ""
	}
	return account.Name
}

func accountPlatform(account *Account) string {
	if account == nil {
		return ""
	}
	return account.Platform
}

func accountStatus(account *Account) string {
	if account == nil {
		return ""
	}
	return account.Status
}

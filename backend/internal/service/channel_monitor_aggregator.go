package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strings"
)

// 渠道监控聚合层：把 latest + availability 拼成 admin/user 视图所需的 summary / detail。
// 所有方法都遵守"失败仅日志，返回零值"的原则，避免 N+1 查询失败拖垮列表渲染。

// BatchMonitorStatusSummary 批量聚合多个监控的 latest + 7d 可用率（admin/user list 用，消除 N+1）。
// 失败时返回空 map，错误仅日志，不影响列表渲染。
//
// 参数：
//   - ids: 要聚合的 monitor ID 列表
//   - primaryByID: monitor ID -> primary model（用于读 7d 可用率与 latest 状态）
//   - extrasByID: monitor ID -> extra models 列表（用于读 latest 状态填充 ExtraModels）
func (s *ChannelMonitorService) BatchMonitorStatusSummary(
	ctx context.Context,
	ids []int64,
	primaryByID map[int64]string,
	extrasByID map[int64][]string,
) map[int64]MonitorStatusSummary {
	out := make(map[int64]MonitorStatusSummary, len(ids))
	if len(ids) == 0 {
		return out
	}
	latestMap, err := s.repo.ListLatestForMonitorIDs(ctx, ids)
	if err != nil {
		slog.Warn("channel_monitor: batch load latest failed", "error", err)
		latestMap = map[int64][]*ChannelMonitorLatest{}
	}
	availMap, err := s.repo.ComputeAvailabilityForMonitors(ctx, ids, monitorAvailability7Days)
	if err != nil {
		slog.Warn("channel_monitor: batch compute availability failed", "error", err)
		availMap = map[int64][]*ChannelMonitorAvailability{}
	}

	for _, id := range ids {
		out[id] = buildStatusSummary(
			indexLatestByModel(latestMap[id]),
			indexAvailabilityByModel(availMap[id]),
			primaryByID[id],
			extrasByID[id],
		)
	}
	return out
}

// ListUserView 用户只读视图：列出所有 enabled 监控的概览。
// 使用批量聚合接口避免 N+1：
//
//	1 次查 monitors；
//	1 次批量 latest（含 ping_latency_ms）；
//	1 次批量 7d availability；
//	1 次批量 timeline（主模型最近 N 条）。
func (s *ChannelMonitorService) ListUserView(ctx context.Context) ([]*UserMonitorView, error) {
	monitors, err := s.repo.ListEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("list enabled monitors: %w", err)
	}
	if len(monitors) == 0 {
		return []*UserMonitorView{}, nil
	}

	ids, primaryByID, extrasByID := collectMonitorIndexes(monitors)
	summaries := s.BatchMonitorStatusSummary(ctx, ids, primaryByID, extrasByID)
	latestMap := s.batchLatest(ctx, ids)
	timelineMap := s.batchTimeline(ctx, ids, primaryByID)

	views := make([]*UserMonitorView, 0, len(monitors))
	for _, m := range monitors {
		primaryLatest := pickLatest(latestMap[m.ID], m.PrimaryModel)
		views = append(views, buildUserViewFromSummary(m, summaries[m.ID], primaryLatest, timelineMap[m.ID]))
	}
	return s.mergeUserViewsByUpstreamPools(ctx, views)
}

// collectMonitorIndexes 把 monitors 列表按 ID 展开为聚合查询所需的三个索引结构。
func collectMonitorIndexes(monitors []*ChannelMonitor) ([]int64, map[int64]string, map[int64][]string) {
	ids := make([]int64, 0, len(monitors))
	primaryByID := make(map[int64]string, len(monitors))
	extrasByID := make(map[int64][]string, len(monitors))
	for _, m := range monitors {
		ids = append(ids, m.ID)
		primaryByID[m.ID] = m.PrimaryModel
		extrasByID[m.ID] = m.ExtraModels
	}
	return ids, primaryByID, extrasByID
}

// batchLatest 批量取 latest per model，失败仅日志（与现有 BatchMonitorStatusSummary 一致，不阻断列表渲染）。
func (s *ChannelMonitorService) batchLatest(ctx context.Context, ids []int64) map[int64][]*ChannelMonitorLatest {
	latestMap, err := s.repo.ListLatestForMonitorIDs(ctx, ids)
	if err != nil {
		slog.Warn("channel_monitor: user view batch latest failed", "error", err)
		return map[int64][]*ChannelMonitorLatest{}
	}
	return latestMap
}

// batchTimeline 批量取每个 monitor 主模型最近 monitorTimelineMaxPoints 条历史。
func (s *ChannelMonitorService) batchTimeline(
	ctx context.Context,
	ids []int64,
	primaryByID map[int64]string,
) map[int64][]*ChannelMonitorHistoryEntry {
	timelineMap, err := s.repo.ListRecentHistoryForMonitors(ctx, ids, primaryByID, monitorTimelineMaxPoints)
	if err != nil {
		slog.Warn("channel_monitor: user view batch timeline failed", "error", err)
		return map[int64][]*ChannelMonitorHistoryEntry{}
	}
	return timelineMap
}

// pickLatest 从 latest 切片中挑出指定 model 对应项，未命中返回 nil。
func pickLatest(rows []*ChannelMonitorLatest, model string) *ChannelMonitorLatest {
	if model == "" {
		return nil
	}
	for _, r := range rows {
		if r.Model == model {
			return r
		}
	}
	return nil
}

type userMonitorPoolBinding struct {
	pool    UpstreamPool
	binding UpstreamPoolBinding
}

type userMonitorPoolBucket struct {
	binding  userMonitorPoolBinding
	monitors []*UserMonitorView
	runtime  userMonitorPoolRuntime
}

type userMonitorPoolRuntime struct {
	activeMembers      int
	schedulableMembers int
}

func (s *ChannelMonitorService) mergeUserViewsByUpstreamPools(ctx context.Context, views []*UserMonitorView) ([]*UserMonitorView, error) {
	if s.upstreamPoolRepo == nil || len(views) == 0 {
		return views, nil
	}
	poolBindings, err := s.loadUserMonitorPoolBindings(ctx)
	if err != nil {
		slog.Warn("channel_monitor: load upstream pools for user view failed", "error", err)
		return views, nil
	}
	if len(poolBindings) == 0 {
		return views, nil
	}
	poolRuntime, err := s.loadUserPoolRuntime(ctx, poolBindings)
	if err != nil {
		slog.Warn("channel_monitor: load upstream pool runtime for user view failed", "error", err)
	}

	buckets := make(map[int64]*userMonitorPoolBucket)
	plain := make([]*UserMonitorView, 0, len(views))
	for _, view := range views {
		binding, ok := matchUserMonitorPoolBinding(view, poolBindings)
		if !ok {
			plain = append(plain, view)
			continue
		}
		bucket := buckets[binding.pool.ID]
		if bucket == nil {
			bucket = &userMonitorPoolBucket{binding: binding, runtime: poolRuntime[binding.pool.ID]}
			buckets[binding.pool.ID] = bucket
		}
		bucket.monitors = append(bucket.monitors, view)
	}

	out := make([]*UserMonitorView, 0, len(plain)+len(buckets))
	used := make(map[int64]struct{}, len(buckets))
	for _, view := range views {
		binding, ok := matchUserMonitorPoolBinding(view, poolBindings)
		if !ok {
			out = append(out, view)
			continue
		}
		if _, exists := used[binding.pool.ID]; exists {
			continue
		}
		used[binding.pool.ID] = struct{}{}
		if merged := buildUserPoolView(buckets[binding.pool.ID]); merged != nil {
			out = append(out, merged)
		}
	}
	return out, nil
}

func (s *ChannelMonitorService) loadUserMonitorPoolBindings(ctx context.Context) ([]userMonitorPoolBinding, error) {
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return nil, err
	}
	poolByID := make(map[int64]UpstreamPool, len(pools))
	for _, pool := range pools {
		if pool.Enabled {
			poolByID[pool.ID] = pool
		}
	}
	if len(poolByID) == 0 {
		return nil, nil
	}
	bindings, err := s.upstreamPoolRepo.ListUpstreamPoolBindings(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]userMonitorPoolBinding, 0, len(bindings))
	for _, binding := range bindings {
		if !binding.Enabled {
			continue
		}
		pool, ok := poolByID[binding.PoolID]
		if !ok {
			continue
		}
		out = append(out, userMonitorPoolBinding{pool: pool, binding: binding})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].binding.Priority != out[j].binding.Priority {
			return out[i].binding.Priority < out[j].binding.Priority
		}
		return out[i].binding.ID < out[j].binding.ID
	})
	return out, nil
}

func matchUserMonitorPoolBinding(view *UserMonitorView, bindings []userMonitorPoolBinding) (userMonitorPoolBinding, bool) {
	if view == nil {
		return userMonitorPoolBinding{}, false
	}
	for _, binding := range bindings {
		if !strings.EqualFold(strings.TrimSpace(binding.pool.Platform), strings.TrimSpace(view.Provider)) {
			continue
		}
		if view.GroupID != nil && *view.GroupID > 0 {
			if binding.binding.GroupID == *view.GroupID {
				return binding, true
			}
			continue
		}
		if monitorGroupMatchesBinding(view.GroupName, binding.binding.GroupName) {
			return binding, true
		}
	}
	return userMonitorPoolBinding{}, false
}

func monitorGroupMatchesBinding(monitorGroup, bindingGroup string) bool {
	monitor := normalizeMonitorGroupName(monitorGroup)
	group := normalizeMonitorGroupName(bindingGroup)
	if monitor == "" || group == "" {
		return false
	}
	return monitor == group
}

func normalizeMonitorGroupName(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.TrimSuffix(v, "组")
	v = strings.TrimSuffix(v, "group")
	return strings.TrimSpace(v)
}

func monitorModelMatchesBinding(primaryModel string, models []string) bool {
	if len(models) == 0 {
		return true
	}
	primary := strings.TrimSpace(primaryModel)
	for _, model := range models {
		if strings.EqualFold(strings.TrimSpace(model), primary) {
			return true
		}
	}
	return false
}

func buildUserPoolView(bucket *userMonitorPoolBucket) *UserMonitorView {
	if bucket == nil || len(bucket.monitors) == 0 {
		return nil
	}
	monitors := bucket.monitors
	base := monitors[0]
	status, latency, ping, availability := aggregateUserMonitorHealth(monitors, bucket.runtime)
	return &UserMonitorView{
		ID:                   userPoolMonitorID(bucket.binding.pool.ID),
		Name:                 bucket.binding.pool.Name,
		Provider:             bucket.binding.pool.Platform,
		GroupID:              &bucket.binding.binding.GroupID,
		GroupName:            bucket.binding.binding.GroupName,
		PrimaryModel:         base.PrimaryModel,
		PrimaryStatus:        status,
		PrimaryLatencyMs:     latency,
		PrimaryPingLatencyMs: ping,
		Availability7d:       availability,
		ExtraModels:          aggregateUserMonitorExtraModels(monitors),
		Timeline:             aggregateUserMonitorTimeline(monitors),
	}
}

func aggregateUserMonitorHealth(monitors []*UserMonitorView, runtime userMonitorPoolRuntime) (string, *int, *int, float64) {
	status := ""
	availability := 0.0
	var latencyMin, pingMin *int
	for _, monitor := range monitors {
		status = betterMonitorStatus(status, monitor.PrimaryStatus)
		if monitor.Availability7d > availability {
			availability = monitor.Availability7d
		}
		if isRoutableMonitorStatus(monitor.PrimaryStatus) {
			latencyMin = minIntPtr(latencyMin, monitor.PrimaryLatencyMs)
			pingMin = minIntPtr(pingMin, monitor.PrimaryPingLatencyMs)
		}
	}
	if status == "" {
		status = MonitorStatusError
	}
	if runtime.schedulableMembers > 0 && !isRoutableMonitorStatus(status) {
		status = MonitorStatusDegraded
	}
	return status, latencyMin, pingMin, availability
}

func (s *ChannelMonitorService) loadUserPoolRuntime(ctx context.Context, bindings []userMonitorPoolBinding) (map[int64]userMonitorPoolRuntime, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil || len(bindings) == 0 {
		return map[int64]userMonitorPoolRuntime{}, nil
	}

	type poolMemberRef struct {
		poolID   int64
		accountID int64
	}

	refs := make([]poolMemberRef, 0)
	accountIDs := make([]int64, 0)
	seenAccounts := make(map[int64]struct{})
	seenPools := make(map[int64]struct{})
	for _, binding := range bindings {
		poolID := binding.pool.ID
		if poolID <= 0 {
			continue
		}
		if _, ok := seenPools[poolID]; ok {
			continue
		}
		seenPools[poolID] = struct{}{}
		members, err := s.upstreamPoolRepo.ListUpstreamPoolMembers(ctx, poolID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			if !member.Enabled || member.ManualDrained {
				continue
			}
			if member.SchedulableOverride != nil && !*member.SchedulableOverride {
				continue
			}
			if member.AccountID <= 0 {
				continue
			}
			refs = append(refs, poolMemberRef{poolID: poolID, accountID: member.AccountID})
			if _, ok := seenAccounts[member.AccountID]; !ok {
				seenAccounts[member.AccountID] = struct{}{}
				accountIDs = append(accountIDs, member.AccountID)
			}
		}
	}
	if len(accountIDs) == 0 {
		return map[int64]userMonitorPoolRuntime{}, nil
	}

	accounts, err := s.accountRepo.GetByIDs(ctx, accountIDs)
	if err != nil {
		return nil, err
	}
	accountByID := make(map[int64]*Account, len(accounts))
	for _, account := range accounts {
		if account != nil {
			accountByID[account.ID] = account
		}
	}

	out := make(map[int64]userMonitorPoolRuntime, len(seenPools))
	for _, ref := range refs {
		runtime := out[ref.poolID]
		account := accountByID[ref.accountID]
		if account == nil {
			out[ref.poolID] = runtime
			continue
		}
		if account.IsActive() {
			runtime.activeMembers++
		}
		if account.IsSchedulable() {
			runtime.schedulableMembers++
		}
		out[ref.poolID] = runtime
	}
	return out, nil
}

func worseMonitorStatus(a, b string) string {
	if monitorStatusRank(b) > monitorStatusRank(a) {
		return b
	}
	return a
}

func betterMonitorStatus(a, b string) string {
	if a == "" || monitorStatusRank(b) < monitorStatusRank(a) {
		return b
	}
	return a
}

func isRoutableMonitorStatus(status string) bool {
	return status == MonitorStatusOperational || status == MonitorStatusDegraded
}

func monitorStatusRank(status string) int {
	switch status {
	case MonitorStatusError:
		return 4
	case MonitorStatusFailed:
		return 3
	case MonitorStatusDegraded:
		return 2
	case MonitorStatusOperational:
		return 1
	default:
		return 0
	}
}

func minIntPtr(current *int, candidate *int) *int {
	if candidate == nil {
		return current
	}
	if current == nil || *candidate < *current {
		v := *candidate
		return &v
	}
	return current
}

func averageIntPtr(total, count int) *int {
	if count == 0 {
		return nil
	}
	v := int(math.Round(float64(total) / float64(count)))
	return &v
}

func averageFloat(total float64, count int) float64 {
	if count == 0 {
		return 0
	}
	return math.Round(total/float64(count)*100) / 100
}

func aggregateUserMonitorExtraModels(monitors []*UserMonitorView) []ExtraModelStatus {
	seen := make(map[string]ExtraModelStatus)
	order := make([]string, 0)
	for _, monitor := range monitors {
		for _, extra := range monitor.ExtraModels {
			key := strings.TrimSpace(extra.Model)
			if key == "" {
				continue
			}
			if _, ok := seen[key]; !ok {
				order = append(order, key)
				seen[key] = extra
				continue
			}
			current := seen[key]
			current.Status = betterMonitorStatus(current.Status, extra.Status)
			if isRoutableMonitorStatus(extra.Status) {
				current.LatencyMs = minIntPtr(current.LatencyMs, extra.LatencyMs)
			}
			seen[key] = current
		}
	}
	out := make([]ExtraModelStatus, 0, len(order))
	for _, key := range order {
		out = append(out, seen[key])
	}
	return out
}

func aggregateUserMonitorTimeline(monitors []*UserMonitorView) []UserMonitorTimelinePoint {
	if len(monitors) == 0 {
		return nil
	}
	maxLen := 0
	for _, monitor := range monitors {
		if len(monitor.Timeline) > maxLen {
			maxLen = len(monitor.Timeline)
		}
	}
	out := make([]UserMonitorTimelinePoint, 0, maxLen)
	for i := 0; i < maxLen; i++ {
		var points []UserMonitorTimelinePoint
		for _, monitor := range monitors {
			if i < len(monitor.Timeline) {
				points = append(points, monitor.Timeline[i])
			}
		}
		if len(points) == 0 {
			continue
		}
		out = append(out, aggregateTimelineSlot(points))
	}
	return out
}

func aggregateTimelineSlot(points []UserMonitorTimelinePoint) UserMonitorTimelinePoint {
	status := ""
	latest := points[0].CheckedAt
	var latencyMin, pingMin *int
	for _, point := range points {
		status = betterMonitorStatus(status, point.Status)
		if point.CheckedAt.After(latest) {
			latest = point.CheckedAt
		}
		if isRoutableMonitorStatus(point.Status) {
			latencyMin = minIntPtr(latencyMin, point.LatencyMs)
			pingMin = minIntPtr(pingMin, point.PingLatencyMs)
		}
	}
	if status == "" {
		status = MonitorStatusError
	}
	return UserMonitorTimelinePoint{
		Status:        status,
		LatencyMs:     latencyMin,
		PingLatencyMs: pingMin,
		CheckedAt:     latest,
	}
}

func userPoolMonitorID(poolID int64) int64 {
	if poolID <= 0 {
		return 0
	}
	return -poolID
}

func userPoolIDFromMonitorID(id int64) (int64, bool) {
	if id >= 0 {
		return 0, false
	}
	return -id, true
}

// GetUserDetail 用户只读视图：单个监控详情（每个模型 7d/15d/30d 可用率与平均延迟）。
// 不暴露 api_key。
func (s *ChannelMonitorService) GetUserDetail(ctx context.Context, id int64) (*UserMonitorDetail, error) {
	if poolID, ok := userPoolIDFromMonitorID(id); ok {
		return s.GetUserPoolDetail(ctx, poolID)
	}
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !m.Enabled {
		return nil, ErrChannelMonitorNotFound
	}

	latest, err := s.repo.ListLatestPerModel(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list latest per model: %w", err)
	}
	availMap, err := s.collectAvailabilityWindows(ctx, id)
	if err != nil {
		return nil, err
	}

	models := mergeModelDetails(m, latest, availMap)
	return &UserMonitorDetail{
		ID:        m.ID,
		Name:      m.Name,
		Provider:  m.Provider,
		GroupName: m.GroupName,
		Models:    models,
	}, nil
}

func (s *ChannelMonitorService) GetUserPoolDetail(ctx context.Context, poolID int64) (*UserMonitorDetail, error) {
	views, err := s.ListUserView(ctx)
	if err != nil {
		return nil, err
	}
	targetID := userPoolMonitorID(poolID)
	for _, view := range views {
		if view.ID != targetID {
			continue
		}
		return &UserMonitorDetail{
			ID:        view.ID,
			Name:      view.Name,
			Provider:  view.Provider,
			GroupName: view.GroupName,
			Models: []ModelDetail{{
				Model:           view.PrimaryModel,
				LatestStatus:    view.PrimaryStatus,
				LatestLatencyMs: view.PrimaryLatencyMs,
				Availability7d:  view.Availability7d,
				Availability15d: view.Availability7d,
				Availability30d: view.Availability7d,
				AvgLatency7dMs:  view.PrimaryLatencyMs,
			}},
		}, nil
	}
	return nil, ErrChannelMonitorNotFound
}

// collectAvailabilityWindows 一次性查询 7/15/30 天三个窗口，按模型组织。
func (s *ChannelMonitorService) collectAvailabilityWindows(ctx context.Context, monitorID int64) (map[int]map[string]*ChannelMonitorAvailability, error) {
	out := make(map[int]map[string]*ChannelMonitorAvailability, 3)
	windows := []int{monitorAvailability7Days, monitorAvailability15Days, monitorAvailability30Days}
	for _, w := range windows {
		rows, err := s.repo.ComputeAvailability(ctx, monitorID, w)
		if err != nil {
			return nil, fmt.Errorf("compute availability %dd: %w", w, err)
		}
		out[w] = indexAvailabilityByModel(rows)
	}
	return out, nil
}

// ---------- 纯函数 helper（无 IO，可在 batch / 单 monitor / detail 路径复用）----------

// indexLatestByModel 把 latest 切片按 model 索引（小工具，避免在 hot path 重复写）。
func indexLatestByModel(rows []*ChannelMonitorLatest) map[string]*ChannelMonitorLatest {
	m := make(map[string]*ChannelMonitorLatest, len(rows))
	for _, r := range rows {
		m[r.Model] = r
	}
	return m
}

// indexAvailabilityByModel 把 availability 切片按 model 索引。
func indexAvailabilityByModel(rows []*ChannelMonitorAvailability) map[string]*ChannelMonitorAvailability {
	m := make(map[string]*ChannelMonitorAvailability, len(rows))
	for _, r := range rows {
		m[r.Model] = r
	}
	return m
}

// buildStatusSummary 由 latest + availability 字典构造 MonitorStatusSummary。
// 不做任何 IO，纯组装，便于在 batch 与单 monitor 路径复用。
func buildStatusSummary(
	latestByModel map[string]*ChannelMonitorLatest,
	availByModel map[string]*ChannelMonitorAvailability,
	primary string,
	extras []string,
) MonitorStatusSummary {
	summary := MonitorStatusSummary{ExtraModels: make([]ExtraModelStatus, 0, len(extras))}
	if primary != "" {
		if l, ok := latestByModel[primary]; ok {
			summary.PrimaryStatus = l.Status
			summary.PrimaryLatencyMs = l.LatencyMs
		}
		if a, ok := availByModel[primary]; ok {
			summary.Availability7d = a.AvailabilityPct
		}
	}
	for _, model := range extras {
		entry := ExtraModelStatus{Model: model}
		if l, ok := latestByModel[model]; ok {
			entry.Status = l.Status
			entry.LatencyMs = l.LatencyMs
		}
		summary.ExtraModels = append(summary.ExtraModels, entry)
	}
	return summary
}

// buildUserViewFromSummary 用预聚合好的 MonitorStatusSummary + 主模型 latest + timeline 装填 UserMonitorView（无 IO）。
// primaryLatest 可能为 nil（该监控尚无历史）；timelineEntries 可能为空。
func buildUserViewFromSummary(
	m *ChannelMonitor,
	summary MonitorStatusSummary,
	primaryLatest *ChannelMonitorLatest,
	timelineEntries []*ChannelMonitorHistoryEntry,
) *UserMonitorView {
	view := &UserMonitorView{
		ID:               m.ID,
		Name:             m.Name,
		Provider:         m.Provider,
		GroupID:          m.GroupID,
		GroupName:        m.GroupName,
		PrimaryModel:     m.PrimaryModel,
		PrimaryStatus:    summary.PrimaryStatus,
		PrimaryLatencyMs: summary.PrimaryLatencyMs,
		Availability7d:   summary.Availability7d,
		ExtraModels:      summary.ExtraModels,
		Timeline:         buildTimelinePoints(timelineEntries),
	}
	if primaryLatest != nil {
		view.PrimaryPingLatencyMs = primaryLatest.PingLatencyMs
	}
	return view
}

// buildTimelinePoints 把 history entry 裁剪为 timeline 点（去除 message/ID/Model，减小响应体）。
func buildTimelinePoints(entries []*ChannelMonitorHistoryEntry) []UserMonitorTimelinePoint {
	out := make([]UserMonitorTimelinePoint, 0, len(entries))
	for _, e := range entries {
		out = append(out, UserMonitorTimelinePoint{
			Status:        e.Status,
			LatencyMs:     e.LatencyMs,
			PingLatencyMs: e.PingLatencyMs,
			CheckedAt:     e.CheckedAt,
		})
	}
	return out
}

// mergeModelDetails 合并 latest + availability 三个窗口为 ModelDetail 列表。
// 复用 indexLatestByModel，避免在多处重复写 build map 逻辑。
func mergeModelDetails(
	m *ChannelMonitor,
	latest []*ChannelMonitorLatest,
	availMap map[int]map[string]*ChannelMonitorAvailability,
) []ModelDetail {
	all := append([]string{m.PrimaryModel}, m.ExtraModels...)
	latestByModel := indexLatestByModel(latest)
	out := make([]ModelDetail, 0, len(all))
	for _, model := range all {
		d := ModelDetail{Model: model}
		if l, ok := latestByModel[model]; ok {
			d.LatestStatus = l.Status
			d.LatestLatencyMs = l.LatencyMs
		}
		if a, ok := availMap[monitorAvailability7Days][model]; ok {
			d.Availability7d = a.AvailabilityPct
			d.AvgLatency7dMs = a.AvgLatencyMs
		}
		if a, ok := availMap[monitorAvailability15Days][model]; ok {
			d.Availability15d = a.AvailabilityPct
		}
		if a, ok := availMap[monitorAvailability30Days][model]; ok {
			d.Availability30d = a.AvailabilityPct
		}
		out = append(out, d)
	}
	return out
}

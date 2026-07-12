package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

const (
	upstreamHealthAlertRequiredObservations = 3
	upstreamHealthAlertReminderCooldown     = 6 * time.Hour
	upstreamHealthErrorRateThreshold        = 0.30
	upstreamHealthLatencyRatioThreshold     = 1.75
)

type upstreamHealthAlertState struct {
	breaches     int
	active       bool
	lastEmitted  time.Time
	lastObserved upstreamHealthAlertObservation
}

type upstreamHealthAlertObservation struct {
	key              string
	alertType        string
	severity         string
	poolID           int64
	poolName         string
	poolCode         string
	accountID        int64
	accountName      string
	platform         string
	model            string
	reason           string
	message          string
	latencyMs        *int
	medianLatencyMs  *int
	errorRate        *float64
	weightFactor     *float64
	availableMembers *int
	totalMembers     *int
}

type upstreamHealthAlertEvent struct {
	status      string
	observation upstreamHealthAlertObservation
	at          time.Time
}

func (s *AccountMonitorService) evaluateUpstreamHealthAlerts(ctx context.Context, probeRows []*AccountMonitorHistoryRow) error {
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return fmt.Errorf("list upstream pools for health alerts: %w", err)
	}
	enabledPools := filterEnabledPools(pools)
	if len(enabledPools) == 0 {
		s.reconcileUpstreamHealthAlerts(nil)
		return nil
	}

	membersByPool, accountsByID := s.loadPoolMembersAndAccounts(ctx, enabledPools)
	rowsByPool := groupHealthAlertProbeRows(probeRows)
	statesByPool, err := s.loadHealthAlertRuntimeWeights(ctx, enabledPools)
	if err != nil {
		return err
	}

	now := s.currentTime().UTC()
	observations := make([]upstreamHealthAlertObservation, 0)
	for _, pool := range enabledPools {
		members := membersByPool[pool.ID]
		available, configured := poolHealthAlertCapacity(members, accountsByID)
		if configured > 1 && available <= 1 {
			alertType := "pool_capacity_low"
			severity := "warning"
			message := "上游池只剩一个可用账号"
			if available == 0 {
				alertType = "pool_unavailable"
				severity = "critical"
				message = "上游池当前没有可用账号"
			}
			observations = append(observations, upstreamHealthAlertObservation{
				key: fmt.Sprintf("%s:pool:%d", alertType, pool.ID), alertType: alertType, severity: severity,
				poolID: pool.ID, poolName: pool.Name, poolCode: pool.Code, platform: pool.Platform,
				reason: "available_members", message: message,
				availableMembers: alertIntPtr(available), totalMembers: alertIntPtr(configured),
			})
		}

		medianLatency := medianOperationalProbeLatency(rowsByPool[pool.ID])
		for _, member := range members {
			account := accountsByID[member.AccountID]
			if account == nil || !member.Enabled || member.ManualDrained {
				continue
			}
			base := upstreamHealthAlertObservation{
				poolID: pool.ID, poolName: pool.Name, poolCode: pool.Code,
				accountID: account.ID, accountName: account.Name, platform: pool.Platform,
			}
			if account.RateLimitResetAt != nil && now.Before(*account.RateLimitResetAt) {
				observation := base
				observation.key = fmt.Sprintf("account_rate_limited:pool:%d:account:%d", pool.ID, account.ID)
				observation.alertType = "account_rate_limited"
				observation.severity = "warning"
				observation.reason = "rate_limit_reset_pending"
				observation.message = "上游账号持续处于限流状态"
				observations = append(observations, observation)
			}
			if member.RuntimeErrorRate != nil && *member.RuntimeErrorRate >= upstreamHealthErrorRateThreshold {
				observation := base
				observation.key = fmt.Sprintf("account_error_rate_high:pool:%d:account:%d", pool.ID, account.ID)
				observation.alertType = "account_error_rate_high"
				observation.severity = "warning"
				observation.reason = "runtime_error_rate"
				observation.message = "上游账号错误率持续偏高"
				observation.errorRate = float64PtrValue(*member.RuntimeErrorRate)
				observations = append(observations, observation)
			}

			row := rowsByPool[pool.ID][account.ID]
			if row != nil {
				base.model = row.Model
				switch row.Status {
				case MonitorStatusFailed, MonitorStatusError:
					observation := base
					observation.key = fmt.Sprintf("account_probe_failed:pool:%d:account:%d", pool.ID, account.ID)
					observation.alertType = "account_probe_failed"
					observation.severity = "critical"
					observation.reason = row.Status
					observation.message = "上游账号连续探测失败"
					observation.latencyMs = row.LatencyMs
					observations = append(observations, observation)
				case MonitorStatusOperational:
					if row.LatencyMs != nil && *row.LatencyMs > 0 && medianLatency > 0 && float64(*row.LatencyMs)/float64(medianLatency) >= upstreamHealthLatencyRatioThreshold {
						observation := base
						observation.key = fmt.Sprintf("account_latency_degraded:pool:%d:account:%d", pool.ID, account.ID)
						observation.alertType = "account_latency_degraded"
						observation.severity = "warning"
						observation.reason = "pool_relative_latency"
						observation.message = "上游账号延迟持续明显高于同池账号"
						observation.latencyMs = row.LatencyMs
						observation.medianLatencyMs = alertIntPtr(medianLatency)
						observations = append(observations, observation)
					}
				}
			}

			if state := statesByPool[pool.ID][account.ID]; state != nil && now.Sub(state.LastObservedAt) <= poolRuntimeWeightTTL && state.Factor <= 0.5+1e-9 {
				observation := base
				observation.key = fmt.Sprintf("account_runtime_weight_low:pool:%d:account:%d", pool.ID, account.ID)
				observation.alertType = "account_runtime_weight_low"
				observation.severity = "warning"
				observation.reason = state.Reason
				observation.message = "上游账号自动调权持续处于低位"
				observation.weightFactor = float64PtrValue(state.Factor)
				observations = append(observations, observation)
			}
		}
	}

	s.reconcileUpstreamHealthAlerts(observations)
	return nil
}

func (s *AccountMonitorService) loadHealthAlertRuntimeWeights(ctx context.Context, pools []UpstreamPool) (map[int64]map[int64]*PoolRuntimeWeightState, error) {
	poolIDs := make([]int64, 0, len(pools))
	for _, pool := range pools {
		if pool.AutoWeightEnabled && strings.EqualFold(pool.Platform, PlatformOpenAI) {
			poolIDs = append(poolIDs, pool.ID)
		}
	}
	if len(poolIDs) == 0 {
		return map[int64]map[int64]*PoolRuntimeWeightState{}, nil
	}
	states, err := s.repo.ListPoolRuntimeWeightStates(ctx, poolIDs)
	if err != nil {
		return nil, fmt.Errorf("list runtime weights for health alerts: %w", err)
	}
	return states, nil
}

func (s *AccountMonitorService) reconcileUpstreamHealthAlerts(observations []upstreamHealthAlertObservation) {
	if s == nil {
		return
	}
	now := s.currentTime().UTC()
	present := make(map[string]upstreamHealthAlertObservation, len(observations))
	for _, observation := range observations {
		if strings.TrimSpace(observation.key) != "" {
			present[observation.key] = observation
		}
	}

	events := make([]upstreamHealthAlertEvent, 0)
	s.healthAlertMu.Lock()
	if s.healthAlertStates == nil {
		s.healthAlertStates = make(map[string]*upstreamHealthAlertState)
	}
	keys := make([]string, 0, len(present))
	for key := range present {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		observation := present[key]
		state := s.healthAlertStates[key]
		if state == nil {
			state = &upstreamHealthAlertState{}
			s.healthAlertStates[key] = state
		}
		state.breaches++
		state.lastObserved = observation
		if state.breaches < upstreamHealthAlertRequiredObservations {
			continue
		}
		if !state.active {
			state.active = true
			state.lastEmitted = now
			events = append(events, upstreamHealthAlertEvent{status: "firing", observation: observation, at: now})
		} else if now.Sub(state.lastEmitted) >= upstreamHealthAlertReminderCooldown {
			state.lastEmitted = now
			events = append(events, upstreamHealthAlertEvent{status: "reminder", observation: observation, at: now})
		}
	}
	for key, state := range s.healthAlertStates {
		if _, ok := present[key]; ok {
			continue
		}
		if state.active {
			events = append(events, upstreamHealthAlertEvent{status: "resolved", observation: state.lastObserved, at: now})
		}
		delete(s.healthAlertStates, key)
	}
	s.healthAlertMu.Unlock()

	emitter := s.healthAlertEmitter
	if emitter == nil {
		emitter = emitUpstreamHealthAlert
	}
	for _, event := range events {
		emitter(event)
	}
}

func groupHealthAlertProbeRows(rows []*AccountMonitorHistoryRow) map[int64]map[int64]*AccountMonitorHistoryRow {
	out := make(map[int64]map[int64]*AccountMonitorHistoryRow)
	for _, row := range rows {
		if row == nil || row.PoolID == nil || *row.PoolID <= 0 || row.AccountID <= 0 {
			continue
		}
		if out[*row.PoolID] == nil {
			out[*row.PoolID] = make(map[int64]*AccountMonitorHistoryRow)
		}
		out[*row.PoolID][row.AccountID] = row
	}
	return out
}

func poolHealthAlertCapacity(members []UpstreamPoolMember, accounts map[int64]*Account) (available, configured int) {
	seen := make(map[int64]struct{})
	for _, member := range members {
		if !member.Enabled || member.ManualDrained || member.AccountID <= 0 {
			continue
		}
		if _, ok := seen[member.AccountID]; ok {
			continue
		}
		seen[member.AccountID] = struct{}{}
		configured++
		if poolMemberSchedulable(member, accounts[member.AccountID]) {
			available++
		}
	}
	return available, configured
}

func emitUpstreamHealthAlert(event upstreamHealthAlertEvent) {
	observation := event.observation
	level := "warn"
	if event.status == "resolved" {
		level = "info"
	}
	fields := map[string]any{
		"alert_key": event.observation.key, "alert_type": observation.alertType,
		"alert_status": event.status, "severity": observation.severity,
		"pool_id": observation.poolID, "pool_name": observation.poolName, "pool_code": observation.poolCode,
		"account_id": observation.accountID, "account_name": observation.accountName,
		"platform": observation.platform, "model": observation.model, "reason": observation.reason,
	}
	if observation.latencyMs != nil {
		fields["latency_ms"] = *observation.latencyMs
	}
	if observation.medianLatencyMs != nil {
		fields["pool_median_latency_ms"] = *observation.medianLatencyMs
	}
	if observation.errorRate != nil {
		fields["error_rate"] = math.Round(*observation.errorRate*10000) / 10000
	}
	if observation.weightFactor != nil {
		fields["runtime_weight_factor"] = *observation.weightFactor
	}
	if observation.availableMembers != nil {
		fields["available_members"] = *observation.availableMembers
	}
	if observation.totalMembers != nil {
		fields["total_members"] = *observation.totalMembers
	}
	message := observation.message
	if event.status == "resolved" {
		message += "（已恢复）"
	} else if event.status == "reminder" {
		message += "（仍在持续）"
	}
	logger.WriteSinkEvent(level, "upstream.health_alert", message, fields)
}

func alertIntPtr(value int) *int { return &value }

func float64PtrValue(value float64) *float64 { return &value }

package service

import "time"

type AccountMonitorHistoryRow struct {
	AccountID     int64
	PoolID        *int64
	GroupID       *int64
	Provider      string
	Model         string
	Status        string
	LatencyMs     *int
	PingLatencyMs *int
	Message       string
	CheckedAt     time.Time
}

type AccountMonitorHistoryEntry struct {
	ID            int64
	AccountID     int64
	PoolID        *int64
	GroupID       *int64
	Provider      string
	Model         string
	Status        string
	LatencyMs     *int
	PingLatencyMs *int
	Message       string
	CheckedAt     time.Time
}

type AccountMonitorLatest struct {
	AccountID     int64
	Provider      string
	Model         string
	Status        string
	LatencyMs     *int
	PingLatencyMs *int
	CheckedAt     time.Time
}

type AccountMonitorAvailability struct {
	AccountID         int64
	Model             string
	WindowDays        int
	TotalChecks       int
	OperationalChecks int
	AvailabilityPct   float64
	AvgLatencyMs      *int
}

type PoolAvailabilitySnapshotRow struct {
	PoolID           int64
	Status           string
	TotalMembers     int
	AvailableMembers int
	CheckedAt        time.Time
}

type PoolAvailabilitySnapshotEntry struct {
	ID               int64
	PoolID           int64
	Status           string
	TotalMembers     int
	AvailableMembers int
	CheckedAt        time.Time
}

type PoolRuntimeWeightState struct {
	PoolID          int64
	AccountID       int64
	Factor          float64
	TargetFactor    float64
	HealthyStreak   int
	UnhealthyStreak int
	Reason          string
	LastObservedAt  time.Time
	UpdatedAt       time.Time
}

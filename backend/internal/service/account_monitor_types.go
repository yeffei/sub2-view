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

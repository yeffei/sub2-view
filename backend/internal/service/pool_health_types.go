package service

import (
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var ErrPoolHealthNotFound = infraerrors.NotFound("UPSTREAM_POOL_HEALTH_NOT_FOUND", "upstream pool health not found")

type PoolHealthView struct {
	ID                  int64
	Name                string
	Provider            string
	GroupID             *int64
	GroupName           string
	Status              string
	Availability7d      float64
	BestLatencyMs       *int
	BestPingLatencyMs   *int
	HealthyMemberCount  int
	DegradedMemberCount int
	FailedMemberCount   int
	TotalMemberCount    int
	Timeline            []PoolHealthTimelinePoint
}

type PoolHealthTimelinePoint struct {
	Status        string
	LatencyMs     *int
	PingLatencyMs *int
	CheckedAt     time.Time
}

type PoolHealthDetail struct {
	ID                  int64
	Name                string
	Provider            string
	GroupID             *int64
	GroupName           string
	Status              string
	Availability7d      float64
	Availability15d     float64
	Availability30d     float64
	BestLatencyMs       *int
	BestPingLatencyMs   *int
	HealthyMemberCount  int
	DegradedMemberCount int
	FailedMemberCount   int
	TotalMemberCount    int
	Timeline            []PoolHealthTimelinePoint
	Members             []PoolHealthMemberItem
	Lines               []PoolHealthLineItem
}

type PoolHealthMemberItem struct {
	AccountID     int64
	AccountName   string
	Platform      string
	AccountStatus string
	HealthStatus  string
	RuntimeStatus string
	RuntimeReason string
	Schedulable   bool
	Enabled       bool
	ManualDrained bool
	Weight        int
	SourceType    string
	SourceSetName string
}

type PoolHealthLineItem struct {
	AccountID           int64
	AccountName         string
	GroupID             *int64
	GroupName           string
	ProbeModel          string
	LatestStatus        string
	LatestLatencyMs     *int
	LatestPingLatencyMs *int
	Availability7d      float64
	Availability15d     float64
	Availability30d     float64
	LastCheckedAt       *time.Time
}

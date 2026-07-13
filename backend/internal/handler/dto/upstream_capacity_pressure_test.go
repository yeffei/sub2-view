package dto

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpstreamCapacityPressureFromService(t *testing.T) {
	hardLimit := 1000
	softShare := 800
	source := &service.UpstreamCapacityPressure{
		SetID:              7,
		SetName:            "聪明共享容量",
		SetCode:            "smart-capacity",
		Platform:           "openai",
		Enabled:            true,
		CapacityLimit:      3000,
		CurrentConcurrency: 2400,
		AvailableCapacity:  600,
		WaitingCount:       2,
		GroupFullCount:     3,
		MemberFullCount:    4,
		BorrowedSlotCount:  5,
		Members: []service.UpstreamCapacityMemberPressure{{
			AccountID:            11,
			AccountName:          "smart-1",
			HardConcurrencyLimit: &hardLimit,
			SoftConcurrencyShare: &softShare,
			CurrentConcurrency:   900,
			WaitingCount:         1,
			LoadRate:             90,
		}},
	}

	out := UpstreamCapacityPressureFromService(source)
	require.NotNil(t, out)
	require.Equal(t, int64(7), out.SetID)
	require.Equal(t, 3000, out.CapacityLimit)
	require.Equal(t, 5, out.BorrowedSlotCount)
	require.Len(t, out.Members, 1)
	require.Equal(t, int64(11), out.Members[0].AccountID)
	require.Equal(t, 90, out.Members[0].LoadRate)
	require.NotSame(t, source.Members[0].HardConcurrencyLimit, out.Members[0].HardConcurrencyLimit)
	require.NotSame(t, source.Members[0].SoftConcurrencyShare, out.Members[0].SoftConcurrencyShare)
}

func TestUpstreamCapacityPressureFromServiceNil(t *testing.T) {
	require.Nil(t, UpstreamCapacityPressureFromService(nil))
}

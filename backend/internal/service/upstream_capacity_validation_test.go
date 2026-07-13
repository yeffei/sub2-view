//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeUpstreamAccountSetForCreate_SharedConcurrencyLimit(t *testing.T) {
	validLimit := 3000
	valid := &UpstreamAccountSet{
		Name:                   "聪明共享容量",
		Code:                   "smart-shared-capacity",
		Platform:               PlatformOpenAI,
		Enabled:                true,
		SharedConcurrencyLimit: &validLimit,
	}
	require.NoError(t, normalizeUpstreamAccountSetForCreate(valid))
	require.NotNil(t, valid.SharedConcurrencyLimit)
	require.Equal(t, 3000, *valid.SharedConcurrencyLimit)

	for _, invalid := range []int{0, -1} {
		t.Run("invalid", func(t *testing.T) {
			value := invalid
			item := &UpstreamAccountSet{
				Name:                   "invalid",
				Code:                   "invalid",
				Platform:               PlatformOpenAI,
				SharedConcurrencyLimit: &value,
			}
			err := normalizeUpstreamAccountSetForCreate(item)
			require.EqualError(t, err, "shared_concurrency_limit must be greater than 0")
		})
	}
}

func TestClonePositiveIntPointer_CopiesValue(t *testing.T) {
	value := 1000
	cloned := clonePositiveIntPointer(&value)
	require.NotNil(t, cloned)
	require.Equal(t, value, *cloned)
	require.NotSame(t, &value, cloned)
	require.Nil(t, clonePositiveIntPointer(nil))
}

func TestValidateUpstreamAccountSetCapacityConfig(t *testing.T) {
	groupLimit := 3000
	hardLimit := 1000
	softShare := 1000
	require.NoError(t, ValidateUpstreamAccountSetCapacityConfig(
		&groupLimit,
		[]UpstreamAccountSetCapacityMemberConfig{
			{AccountID: 1, HardConcurrencyLimit: &hardLimit, SoftConcurrencyShare: &softShare},
			{AccountID: 2, HardConcurrencyLimit: &hardLimit, SoftConcurrencyShare: &softShare},
			{AccountID: 3, HardConcurrencyLimit: &hardLimit, SoftConcurrencyShare: &softShare},
		},
	))

	t.Run("members require group limit", func(t *testing.T) {
		err := ValidateUpstreamAccountSetCapacityConfig(nil, []UpstreamAccountSetCapacityMemberConfig{{AccountID: 1}})
		require.EqualError(t, err, "shared_concurrency_limit is required when capacity members are configured")
	})

	t.Run("member hard limit cannot exceed group", func(t *testing.T) {
		tooHigh := 3001
		err := ValidateUpstreamAccountSetCapacityConfig(
			&groupLimit,
			[]UpstreamAccountSetCapacityMemberConfig{{AccountID: 1, HardConcurrencyLimit: &tooHigh}},
		)
		require.EqualError(t, err, "capacity member 1 hard_concurrency_limit must not exceed shared_concurrency_limit")
	})

	t.Run("duplicate member rejected", func(t *testing.T) {
		err := ValidateUpstreamAccountSetCapacityConfig(
			&groupLimit,
			[]UpstreamAccountSetCapacityMemberConfig{{AccountID: 1}, {AccountID: 1}},
		)
		require.EqualError(t, err, "capacity member account_id 1 is duplicated")
	})

	t.Run("soft share must be positive", func(t *testing.T) {
		invalid := 0
		err := ValidateUpstreamAccountSetCapacityConfig(
			&groupLimit,
			[]UpstreamAccountSetCapacityMemberConfig{{AccountID: 1, SoftConcurrencyShare: &invalid}},
		)
		require.EqualError(t, err, "capacity member 1 soft_concurrency_share must be greater than 0")
	})
}

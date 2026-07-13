//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpstreamAccountSetRepository_PersistsSharedCapacityConfig(t *testing.T) {
	ctx := context.Background()
	repo := NewUpstreamPoolRepository(nil, integrationDB)
	suffix := time.Now().UnixNano()
	limit := 3000
	set := &service.UpstreamAccountSet{
		Name:                   fmt.Sprintf("capacity-%d", suffix),
		Code:                   fmt.Sprintf("capacity-%d", suffix),
		Platform:               service.PlatformOpenAI,
		Enabled:                true,
		SharedConcurrencyLimit: &limit,
	}

	created, err := repo.CreateUpstreamAccountSet(ctx, set)
	require.NoError(t, err)
	require.Positive(t, created.ID)
	t.Cleanup(func() {
		_ = repo.DeleteUpstreamAccountSet(context.Background(), created.ID)
	})

	var accountID int64
	err = integrationDB.QueryRowContext(
		ctx,
		"INSERT INTO accounts (name, platform, type) VALUES ($1, $2, $3) RETURNING id",
		fmt.Sprintf("capacity-account-%d", suffix),
		service.PlatformOpenAI,
		service.AccountTypeAPIKey,
	).Scan(&accountID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = integrationDB.ExecContext(context.Background(), "DELETE FROM accounts WHERE id = $1", accountID)
	})

	require.NoError(t, repo.AddUpstreamAccountSetMembers(ctx, created.ID, []int64{accountID}))
	_, err = integrationDB.ExecContext(
		ctx,
		"INSERT INTO upstream_account_set_capacity_members (account_id, set_id, hard_concurrency_limit, soft_concurrency_share) VALUES ($1, $2, $3, $4)",
		accountID,
		created.ID,
		1000,
		1000,
	)
	require.NoError(t, err)

	loaded, err := repo.GetUpstreamAccountSetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, loaded.SharedConcurrencyLimit)
	require.Equal(t, 3000, *loaded.SharedConcurrencyLimit)

	members, err := repo.ListUpstreamAccountSetMembers(ctx, created.ID)
	require.NoError(t, err)
	require.Len(t, members, 1)
	require.NotNil(t, members[0].CapacityHardLimit)
	require.NotNil(t, members[0].CapacitySoftShare)
	require.Equal(t, 1000, *members[0].CapacityHardLimit)
	require.Equal(t, 1000, *members[0].CapacitySoftShare)

	accountRepo := NewAccountRepository(integrationEntClient, integrationDB, nil)
	loadedAccount, err := accountRepo.GetByID(ctx, accountID)
	require.NoError(t, err)
	require.NotNil(t, loadedAccount.CapacityScope)
	require.Equal(t, created.ID, loadedAccount.CapacityScope.GroupID)
	require.Equal(t, 3000, loadedAccount.CapacityScope.GroupLimit)
	require.Equal(t, 1000, loadedAccount.CapacityScope.MemberHardLimit)
	require.Equal(t, 1000, loadedAccount.CapacityScope.MemberSoftShare)

	_, err = integrationDB.ExecContext(
		ctx,
		"UPDATE upstream_account_set_capacity_members SET hard_concurrency_limit = $1 WHERE account_id = $2",
		3001,
		accountID,
	)
	require.ErrorContains(t, err, "hard_concurrency_limit must not exceed shared_concurrency_limit")

	loaded.SharedConcurrencyLimit = nil
	_, err = repo.UpdateUpstreamAccountSet(ctx, loaded)
	require.ErrorContains(t, err, "shared_concurrency_limit cannot be cleared while capacity members exist")

	updatedLimit := 4000
	loaded.SharedConcurrencyLimit = &updatedLimit
	updated, err := repo.UpdateUpstreamAccountSet(ctx, loaded)
	require.NoError(t, err)
	require.NotNil(t, updated.SharedConcurrencyLimit)
	require.Equal(t, 4000, *updated.SharedConcurrencyLimit)

	reloaded, err := repo.GetUpstreamAccountSetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, reloaded.SharedConcurrencyLimit)
	require.Equal(t, 4000, *reloaded.SharedConcurrencyLimit)
}

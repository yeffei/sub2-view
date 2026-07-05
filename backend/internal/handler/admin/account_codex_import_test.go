package admin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportCodexSessionsDryRunDoesNotCreateAccounts(t *testing.T) {
	adminSvc := newStubAdminService()
	handler := NewAccountHandler(adminSvc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	dryRun := true

	result, err := handler.importCodexSessions(context.Background(), CodexSessionImportRequest{
		DryRun: &dryRun,
	}, []codexImportEntry{{
		Index: 1,
		Value: map[string]any{
			"access_token":  "access-token-1",
			"refresh_token": "refresh-token-1",
			"email":         "codex@example.com",
		},
	}})

	require.NoError(t, err)
	require.True(t, result.DryRun)
	require.NotEmpty(t, result.BatchID)
	require.Equal(t, 1, result.Total)
	require.Equal(t, 1, result.Created)
	require.Zero(t, result.Updated)
	require.Zero(t, result.Failed)
	require.Empty(t, adminSvc.createdAccounts)
}

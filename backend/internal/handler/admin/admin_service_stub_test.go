package admin

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type stubAdminService struct {
	users                  []service.User
	apiKeys                []service.APIKey
	groups                 []service.Group
	accounts               []service.Account
	upstreamPools          []service.UpstreamPool
	upstreamPoolMembers    []service.UpstreamPoolMember
	accountSets            []service.UpstreamAccountSet
	accountSetMembers      []service.UpstreamAccountSetMember
	upstreamPoolMemberSets []service.UpstreamPoolMemberSet
	upstreamPoolBindings   []service.UpstreamPoolBinding
	capacityPressures      []service.UpstreamCapacityPressure
	proxies                []service.Proxy
	proxyCounts            []service.ProxyWithAccountCount
	redeems                []service.RedeemCode
	boundAuthIdentity      *service.AdminBindAuthIdentityInput
	boundAuthIdentityFor   int64
	createdAccounts        []*service.CreateAccountInput
	createdProxies         []*service.CreateProxyInput
	updatedProxyIDs        []int64
	updatedProxies         []*service.UpdateProxyInput
	testedProxyIDs         []int64
	getUserErr             error
	createAccountErr       error
	updateAccountErr       error
	getAccountsByIDsFunc   func(context.Context, []int64) ([]*service.Account, error)
	bulkUpdateAccountErr   error
	checkMixedErr          error
	lastMixedCheck         struct {
		accountID int64
		platform  string
		groupIDs  []int64
	}
	lastListAccounts struct {
		platform      string
		accountType   string
		status        string
		anomalyReason string
		search        string
		groupID       int64
		privacyMode   string
		sortBy        string
		sortOrder     string
		calls         int
	}
	lastListUsers struct {
		page      int
		pageSize  int
		filters   service.UserListFilters
		sortBy    string
		sortOrder string
		calls     int
	}
	lastListProxies struct {
		protocol  string
		status    string
		search    string
		sortBy    string
		sortOrder string
		calls     int
	}
	lastListRedeemCodes struct {
		codeType  string
		status    string
		search    string
		sortBy    string
		sortOrder string
		calls     int
	}
	mu sync.Mutex
}

func newStubAdminService() *stubAdminService {
	now := time.Now().UTC()
	user := service.User{
		ID:        1,
		Email:     "user@example.com",
		Role:      service.RoleUser,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	apiKey := service.APIKey{
		ID:        10,
		UserID:    user.ID,
		Key:       "sk-test",
		Name:      "test",
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	group := service.Group{
		ID:        2,
		Name:      "group",
		Platform:  service.PlatformAnthropic,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	account := service.Account{
		ID:        3,
		Name:      "account",
		Platform:  service.PlatformAnthropic,
		Type:      service.AccountTypeOAuth,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	openAIAccount := service.Account{
		ID:        13,
		Name:      "openai-account",
		Platform:  service.PlatformOpenAI,
		Type:      service.AccountTypeOAuth,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	proxy := service.Proxy{
		ID:        4,
		Name:      "proxy",
		Protocol:  "http",
		Host:      "127.0.0.1",
		Port:      8080,
		Status:    service.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	redeem := service.RedeemCode{
		ID:        5,
		Code:      "R-TEST",
		Type:      service.RedeemTypeBalance,
		Value:     10,
		Status:    service.StatusUnused,
		CreatedAt: now,
	}
	return &stubAdminService{
		users:    []service.User{user},
		apiKeys:  []service.APIKey{apiKey},
		groups:   []service.Group{group},
		accounts: []service.Account{account, openAIAccount},
		upstreamPools: []service.UpstreamPool{{
			ID:        11,
			Name:      "pool",
			Code:      "openai-main",
			Platform:  service.PlatformOpenAI,
			Enabled:   true,
			CreatedAt: now,
			UpdatedAt: now,
		}},
		upstreamPoolMembers: []service.UpstreamPoolMember{{
			ID:              21,
			PoolID:          11,
			AccountID:       openAIAccount.ID,
			AccountName:     openAIAccount.Name,
			AccountPlatform: openAIAccount.Platform,
			AccountType:     openAIAccount.Type,
			AccountStatus:   openAIAccount.Status,
			Enabled:         true,
			Weight:          100,
			JoinedAt:        now,
			UpdatedAt:       now,
		}},
		accountSets: []service.UpstreamAccountSet{{
			ID:           41,
			Name:         "set",
			Code:         "openai-set",
			Platform:     service.PlatformOpenAI,
			Enabled:      true,
			AccountCount: 1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
		accountSetMembers: []service.UpstreamAccountSetMember{{
			SetID:           41,
			AccountID:       openAIAccount.ID,
			AccountName:     openAIAccount.Name,
			AccountPlatform: openAIAccount.Platform,
			AccountType:     openAIAccount.Type,
			AccountStatus:   openAIAccount.Status,
			AddedAt:         now,
		}},
		upstreamPoolMemberSets: []service.UpstreamPoolMemberSet{{
			ID:          51,
			PoolID:      11,
			SetID:       41,
			SetName:     "set",
			SetCode:     "openai-set",
			SetPlatform: service.PlatformOpenAI,
			Enabled:     true,
			JoinedAt:    now,
			UpdatedAt:   now,
		}},
		upstreamPoolBindings: []service.UpstreamPoolBinding{{
			ID:        31,
			GroupID:   group.ID,
			PoolID:    11,
			Platform:  service.PlatformOpenAI,
			Enabled:   true,
			CreatedAt: now,
			UpdatedAt: now,
		}},
		proxies:     []service.Proxy{proxy},
		proxyCounts: []service.ProxyWithAccountCount{{Proxy: proxy, AccountCount: 1}},
		redeems:     []service.RedeemCode{redeem},
	}
}

func (s *stubAdminService) ListUsers(ctx context.Context, page, pageSize int, filters service.UserListFilters, sortBy, sortOrder string) ([]service.User, int64, error) {
	s.lastListUsers.page = page
	s.lastListUsers.pageSize = pageSize
	s.lastListUsers.filters = filters
	s.lastListUsers.sortBy = sortBy
	s.lastListUsers.sortOrder = sortOrder
	s.lastListUsers.calls++
	return s.users, int64(len(s.users)), nil
}

func (s *stubAdminService) GetUser(ctx context.Context, id int64) (*service.User, error) {
	if s.getUserErr != nil {
		return nil, s.getUserErr
	}
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i], nil
		}
	}
	user := service.User{ID: id, Email: "user@example.com", Status: service.StatusActive}
	return &user, nil
}

func (s *stubAdminService) GetUserIncludeDeleted(ctx context.Context, id int64) (*service.User, error) {
	return s.GetUser(ctx, id)
}

func (s *stubAdminService) CreateUser(ctx context.Context, input *service.CreateUserInput) (*service.User, error) {
	user := service.User{ID: 100, Email: input.Email, Status: service.StatusActive}
	return &user, nil
}

func (s *stubAdminService) UpdateUser(ctx context.Context, id int64, input *service.UpdateUserInput) (*service.User, error) {
	user := service.User{ID: id, Email: "updated@example.com", Status: service.StatusActive}
	return &user, nil
}

func (s *stubAdminService) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) UpdateUserBalance(ctx context.Context, userID int64, balance float64, operation string, notes string) (*service.User, error) {
	user := service.User{ID: userID, Balance: balance, Status: service.StatusActive}
	return &user, nil
}

func (s *stubAdminService) BatchUpdateConcurrency(ctx context.Context, userIDs []int64, value int, mode string) (int, error) {
	return len(userIDs), nil
}

func (s *stubAdminService) GetUserAPIKeys(ctx context.Context, userID int64, page, pageSize int, sortBy, sortOrder string) ([]service.APIKey, int64, error) {
	return s.apiKeys, int64(len(s.apiKeys)), nil
}

func (s *stubAdminService) GetUserUsageStats(ctx context.Context, userID int64, period string) (any, error) {
	return map[string]any{"user_id": userID}, nil
}

func (s *stubAdminService) GetUserRPMStatus(ctx context.Context, userID int64) (*service.UserRPMStatus, error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &service.UserRPMStatus{
		UserRPMUsed:  0,
		UserRPMLimit: user.RPMLimit,
	}, nil
}

func (s *stubAdminService) BindUserAuthIdentity(ctx context.Context, userID int64, input service.AdminBindAuthIdentityInput) (*service.AdminBoundAuthIdentity, error) {
	s.boundAuthIdentityFor = userID
	copied := input
	if input.Metadata != nil {
		copied.Metadata = map[string]any{}
		for key, value := range input.Metadata {
			copied.Metadata[key] = value
		}
	}
	if input.Channel != nil {
		channel := *input.Channel
		if input.Channel.Metadata != nil {
			channel.Metadata = map[string]any{}
			for key, value := range input.Channel.Metadata {
				channel.Metadata[key] = value
			}
		}
		copied.Channel = &channel
	}
	s.boundAuthIdentity = &copied

	now := time.Now().UTC()
	result := &service.AdminBoundAuthIdentity{
		UserID:          userID,
		ProviderType:    input.ProviderType,
		ProviderKey:     input.ProviderKey,
		ProviderSubject: input.ProviderSubject,
		VerifiedAt:      &now,
		Issuer:          input.Issuer,
		Metadata:        input.Metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if input.Channel != nil {
		result.Channel = &service.AdminBoundAuthIdentityChannel{
			Channel:        input.Channel.Channel,
			ChannelAppID:   input.Channel.ChannelAppID,
			ChannelSubject: input.Channel.ChannelSubject,
			Metadata:       input.Channel.Metadata,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	return result, nil
}

func (s *stubAdminService) ListGroups(ctx context.Context, page, pageSize int, platform, status, search string, isExclusive *bool, sortBy, sortOrder string) ([]service.Group, int64, error) {
	return s.groups, int64(len(s.groups)), nil
}

func (s *stubAdminService) GetAllGroups(ctx context.Context) ([]service.Group, error) {
	return s.groups, nil
}

func (s *stubAdminService) GetAllGroupsByPlatform(ctx context.Context, platform string) ([]service.Group, error) {
	return s.groups, nil
}

func (s *stubAdminService) GetAllGroupsIncludingInactive(ctx context.Context) ([]service.Group, error) {
	return s.groups, nil
}

func (s *stubAdminService) GetGroup(ctx context.Context, id int64) (*service.Group, error) {
	group := service.Group{ID: id, Name: "group", Status: service.StatusActive}
	return &group, nil
}

func (s *stubAdminService) GetGroupModelsListCandidates(ctx context.Context, id int64, platform string) ([]string, error) {
	if platform == service.PlatformOpenAI {
		return []string{"gpt-5.5", "gpt-5.4"}, nil
	}
	return []string{"claude-sonnet-4-6"}, nil
}

func (s *stubAdminService) CreateGroup(ctx context.Context, input *service.CreateGroupInput) (*service.Group, error) {
	group := service.Group{ID: 200, Name: input.Name, Status: service.StatusActive}
	return &group, nil
}

func (s *stubAdminService) UpdateGroup(ctx context.Context, id int64, input *service.UpdateGroupInput) (*service.Group, error) {
	group := service.Group{ID: id, Name: input.Name, Status: service.StatusActive}
	return &group, nil
}

func (s *stubAdminService) DeleteGroup(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) GetGroupAPIKeys(ctx context.Context, groupID int64, page, pageSize int) ([]service.APIKey, int64, error) {
	return s.apiKeys, int64(len(s.apiKeys)), nil
}

func (s *stubAdminService) ListUpstreamPools(ctx context.Context) ([]service.UpstreamPool, error) {
	return s.upstreamPools, nil
}

func (s *stubAdminService) ListUpstreamCapacityPressures(ctx context.Context) ([]service.UpstreamCapacityPressure, error) {
	return s.capacityPressures, nil
}

func (s *stubAdminService) GetUpstreamPoolByID(ctx context.Context, id int64) (*service.UpstreamPool, error) {
	for i := range s.upstreamPools {
		if s.upstreamPools[i].ID == id {
			return &s.upstreamPools[i], nil
		}
	}
	return nil, service.ErrUpstreamPoolNotFound
}

func (s *stubAdminService) CreateUpstreamPool(ctx context.Context, input *service.CreateUpstreamPoolInput) (*service.UpstreamPool, error) {
	pool := service.UpstreamPool{ID: 1001, Name: input.Name, Code: input.Code, Platform: input.Platform, Enabled: true}
	return &pool, nil
}

func (s *stubAdminService) UpdateUpstreamPool(ctx context.Context, id int64, input *service.UpdateUpstreamPoolInput) (*service.UpstreamPool, error) {
	pool := service.UpstreamPool{ID: id, Name: "pool", Code: "pool-code", Platform: service.PlatformOpenAI, Enabled: true}
	return &pool, nil
}

func (s *stubAdminService) DeleteUpstreamPool(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) GetUpstreamPoolMemberByID(ctx context.Context, id int64) (*service.UpstreamPoolMember, error) {
	for i := range s.upstreamPoolMembers {
		if s.upstreamPoolMembers[i].ID == id {
			return &s.upstreamPoolMembers[i], nil
		}
	}
	return nil, service.ErrUpstreamPoolNotFound
}

func (s *stubAdminService) ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]service.UpstreamPoolMember, error) {
	out := make([]service.UpstreamPoolMember, 0)
	for _, member := range s.upstreamPoolMembers {
		if member.PoolID == poolID {
			out = append(out, member)
		}
	}
	return out, nil
}

func (s *stubAdminService) PreviewUpstreamPoolMemberSync(ctx context.Context, poolID int64, input *service.UpstreamPoolMemberSyncPreviewInput) (*service.UpstreamPoolMemberSyncResult, error) {
	pool, err := s.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	return &service.UpstreamPoolMemberSyncResult{PoolID: pool.ID, Platform: pool.Platform, Mode: service.NormalizeUpstreamPoolMemberSyncMode(input.Mode)}, nil
}

func (s *stubAdminService) ApplyUpstreamPoolMemberSync(ctx context.Context, poolID int64, input *service.UpstreamPoolMemberSyncApplyInput) (*service.UpstreamPoolMemberSyncResult, error) {
	pool, err := s.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	return &service.UpstreamPoolMemberSyncResult{PoolID: pool.ID, Platform: pool.Platform, Mode: service.NormalizeUpstreamPoolMemberSyncMode(input.Mode)}, nil
}

func (s *stubAdminService) CreateUpstreamPoolMember(ctx context.Context, poolID int64, input *service.CreateUpstreamPoolMemberInput) (*service.UpstreamPoolMember, error) {
	member := service.UpstreamPoolMember{ID: 201, PoolID: poolID, AccountID: input.AccountID, Enabled: true, Weight: 100}
	return &member, nil
}

func (s *stubAdminService) UpdateUpstreamPoolMember(ctx context.Context, id int64, input *service.UpdateUpstreamPoolMemberInput) (*service.UpstreamPoolMember, error) {
	member := service.UpstreamPoolMember{ID: id, PoolID: 11, AccountID: 3, Enabled: true, Weight: 100}
	return &member, nil
}

func (s *stubAdminService) DeleteUpstreamPoolMember(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) ListUpstreamAccountSets(ctx context.Context) ([]service.UpstreamAccountSet, error) {
	return s.accountSets, nil
}

func (s *stubAdminService) CreateUpstreamAccountSet(ctx context.Context, input *service.CreateUpstreamAccountSetInput) (*service.UpstreamAccountSet, error) {
	item := service.UpstreamAccountSet{
		ID:          401,
		Name:        input.Name,
		Code:        input.Code,
		Platform:    input.Platform,
		Description: input.Description,
		Enabled:     input.Enabled,
	}
	return &item, nil
}

func (s *stubAdminService) UpdateUpstreamAccountSet(ctx context.Context, id int64, input *service.UpdateUpstreamAccountSetInput) (*service.UpstreamAccountSet, error) {
	item := service.UpstreamAccountSet{
		ID:           id,
		Name:         "set",
		Code:         "openai-set",
		Platform:     service.PlatformOpenAI,
		Enabled:      true,
		AccountCount: len(s.accountSetMembers),
	}
	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Code != nil {
		item.Code = *input.Code
	}
	if input.Platform != nil {
		item.Platform = *input.Platform
	}
	if input.Description != nil {
		item.Description = *input.Description
	}
	if input.Enabled != nil {
		item.Enabled = *input.Enabled
	}
	return &item, nil
}

func (s *stubAdminService) DeleteUpstreamAccountSet(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) ListUpstreamAccountSetMembers(ctx context.Context, setID int64) ([]service.UpstreamAccountSetMember, error) {
	out := make([]service.UpstreamAccountSetMember, 0)
	for _, member := range s.accountSetMembers {
		if member.SetID == setID {
			out = append(out, member)
		}
	}
	return out, nil
}

func (s *stubAdminService) AddUpstreamAccountSetMembers(ctx context.Context, setID int64, input *service.AddUpstreamAccountSetMembersInput) error {
	for _, accountID := range input.AccountIDs {
		s.accountSetMembers = append(s.accountSetMembers, service.UpstreamAccountSetMember{
			SetID:     setID,
			AccountID: accountID,
			AddedAt:   time.Now().UTC(),
		})
	}
	return nil
}

func (s *stubAdminService) DeleteUpstreamAccountSetMember(ctx context.Context, setID, accountID int64) error {
	filtered := s.accountSetMembers[:0]
	for _, member := range s.accountSetMembers {
		if member.SetID == setID && member.AccountID == accountID {
			continue
		}
		filtered = append(filtered, member)
	}
	s.accountSetMembers = filtered
	return nil
}

func (s *stubAdminService) ListUpstreamPoolMemberSets(ctx context.Context, poolID int64) ([]service.UpstreamPoolMemberSet, error) {
	out := make([]service.UpstreamPoolMemberSet, 0)
	for _, item := range s.upstreamPoolMemberSets {
		if item.PoolID == poolID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *stubAdminService) CreateUpstreamPoolMemberSet(ctx context.Context, poolID int64, input *service.CreateUpstreamPoolMemberSetInput) (*service.UpstreamPoolMemberSet, error) {
	item := service.UpstreamPoolMemberSet{
		ID:          501,
		PoolID:      poolID,
		SetID:       input.SetID,
		SetName:     "set",
		SetCode:     "openai-set",
		SetPlatform: service.PlatformOpenAI,
		Enabled:     input.Enabled,
		Notes:       input.Notes,
	}
	return &item, nil
}

func (s *stubAdminService) UpdateUpstreamPoolMemberSet(ctx context.Context, id int64, input *service.UpdateUpstreamPoolMemberSetInput) (*service.UpstreamPoolMemberSet, error) {
	item := service.UpstreamPoolMemberSet{
		ID:          id,
		PoolID:      11,
		SetID:       41,
		SetName:     "set",
		SetCode:     "openai-set",
		SetPlatform: service.PlatformOpenAI,
		Enabled:     true,
	}
	if input.Enabled != nil {
		item.Enabled = *input.Enabled
	}
	if input.Notes != nil {
		item.Notes = *input.Notes
	}
	return &item, nil
}

func (s *stubAdminService) DeleteUpstreamPoolMemberSet(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) ListUpstreamPoolBindings(ctx context.Context) ([]service.UpstreamPoolBinding, error) {
	return s.upstreamPoolBindings, nil
}

func (s *stubAdminService) GetUpstreamPoolBindingByID(ctx context.Context, id int64) (*service.UpstreamPoolBinding, error) {
	for i := range s.upstreamPoolBindings {
		if s.upstreamPoolBindings[i].ID == id {
			return &s.upstreamPoolBindings[i], nil
		}
	}
	return nil, service.ErrUpstreamPoolNotFound
}

func (s *stubAdminService) CreateUpstreamPoolBinding(ctx context.Context, input *service.CreateUpstreamPoolBindingInput) (*service.UpstreamPoolBinding, error) {
	binding := service.UpstreamPoolBinding{ID: 301, GroupID: input.GroupID, PoolID: input.PoolID, Platform: input.Platform, Enabled: true}
	return &binding, nil
}

func (s *stubAdminService) UpdateUpstreamPoolBinding(ctx context.Context, id int64, input *service.UpdateUpstreamPoolBindingInput) (*service.UpstreamPoolBinding, error) {
	binding := service.UpstreamPoolBinding{ID: id, GroupID: 2, PoolID: 11, Platform: service.PlatformOpenAI, Enabled: true}
	return &binding, nil
}

func (s *stubAdminService) DeleteUpstreamPoolBinding(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) GetGroupRateMultipliers(_ context.Context, _ int64) ([]service.UserGroupRateEntry, error) {
	return nil, nil
}

func (s *stubAdminService) ClearGroupRateMultipliers(_ context.Context, _ int64) error {
	return nil
}

func (s *stubAdminService) BatchSetGroupRateMultipliers(_ context.Context, _ int64, _ []service.GroupRateMultiplierInput) error {
	return nil
}

func (s *stubAdminService) ClearGroupRPMOverrides(_ context.Context, _ int64) error {
	return nil
}

func (s *stubAdminService) BatchSetGroupRPMOverrides(_ context.Context, _ int64, _ []service.GroupRPMOverrideInput) error {
	return nil
}

func (s *stubAdminService) ListAccounts(ctx context.Context, page, pageSize int, platform, accountType, status, anomalyReason, search string, groupID int64, privacyMode string, sortBy, sortOrder string) ([]service.Account, int64, error) {
	s.lastListAccounts.platform = platform
	s.lastListAccounts.accountType = accountType
	s.lastListAccounts.status = status
	s.lastListAccounts.anomalyReason = anomalyReason
	s.lastListAccounts.search = search
	s.lastListAccounts.groupID = groupID
	s.lastListAccounts.privacyMode = privacyMode
	s.lastListAccounts.sortBy = sortBy
	s.lastListAccounts.sortOrder = sortOrder
	s.lastListAccounts.calls++
	if pageSize <= 0 {
		return s.accounts, int64(len(s.accounts)), nil
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * pageSize
	if start >= len(s.accounts) {
		return []service.Account{}, int64(len(s.accounts)), nil
	}
	end := start + pageSize
	if end > len(s.accounts) {
		end = len(s.accounts)
	}
	return s.accounts[start:end], int64(len(s.accounts)), nil
}

func (s *stubAdminService) GetAccount(ctx context.Context, id int64) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActive}
	return &account, nil
}

func (s *stubAdminService) GetAccountsByIDs(ctx context.Context, ids []int64) ([]*service.Account, error) {
	if s.getAccountsByIDsFunc != nil {
		return s.getAccountsByIDsFunc(ctx, ids)
	}
	out := make([]*service.Account, 0, len(ids))
	for _, id := range ids {
		account := service.Account{ID: id, Name: "account", Status: service.StatusActive}
		out = append(out, &account)
	}
	return out, nil
}

func (s *stubAdminService) CreateAccount(ctx context.Context, input *service.CreateAccountInput) (*service.Account, error) {
	s.mu.Lock()
	s.createdAccounts = append(s.createdAccounts, input)
	s.mu.Unlock()
	if s.createAccountErr != nil {
		return nil, s.createAccountErr
	}
	account := service.Account{ID: 300, Name: input.Name, Status: service.StatusActive}
	return &account, nil
}

func (s *stubAdminService) UpdateAccount(ctx context.Context, id int64, input *service.UpdateAccountInput) (*service.Account, error) {
	if s.updateAccountErr != nil {
		return nil, s.updateAccountErr
	}
	account := service.Account{ID: id, Name: input.Name, Status: service.StatusActive}
	return &account, nil
}

func (s *stubAdminService) UpdateAccountExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
}

func (s *stubAdminService) DeleteAccount(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) RefreshAccountCredentials(ctx context.Context, id int64) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActive}
	return &account, nil
}

func (s *stubAdminService) ClearAccountError(ctx context.Context, id int64) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActive}
	return &account, nil
}

func (s *stubAdminService) SetAccountError(ctx context.Context, id int64, errorMsg string) error {
	return nil
}

func (s *stubAdminService) SetAccountSchedulable(ctx context.Context, id int64, schedulable bool) (*service.Account, error) {
	account := service.Account{ID: id, Name: "account", Status: service.StatusActive, Schedulable: schedulable}
	return &account, nil
}

func (s *stubAdminService) BulkUpdateAccounts(ctx context.Context, input *service.BulkUpdateAccountsInput) (*service.BulkUpdateAccountsResult, error) {
	if s.bulkUpdateAccountErr != nil {
		return nil, s.bulkUpdateAccountErr
	}
	return &service.BulkUpdateAccountsResult{Success: len(input.AccountIDs), Failed: 0, SuccessIDs: input.AccountIDs}, nil
}

func (s *stubAdminService) CheckMixedChannelRisk(ctx context.Context, currentAccountID int64, currentAccountPlatform string, groupIDs []int64) error {
	s.lastMixedCheck.accountID = currentAccountID
	s.lastMixedCheck.platform = currentAccountPlatform
	s.lastMixedCheck.groupIDs = append([]int64(nil), groupIDs...)
	return s.checkMixedErr
}

func (s *stubAdminService) ListProxies(ctx context.Context, page, pageSize int, protocol, status, search string, sortBy, sortOrder string) ([]service.Proxy, int64, error) {
	s.lastListProxies.protocol = protocol
	s.lastListProxies.status = status
	s.lastListProxies.search = search
	s.lastListProxies.sortBy = sortBy
	s.lastListProxies.sortOrder = sortOrder
	s.lastListProxies.calls++
	search = strings.TrimSpace(strings.ToLower(search))
	filtered := make([]service.Proxy, 0, len(s.proxies))
	for _, proxy := range s.proxies {
		if protocol != "" && proxy.Protocol != protocol {
			continue
		}
		if status != "" && proxy.Status != status {
			continue
		}
		if search != "" {
			name := strings.ToLower(proxy.Name)
			host := strings.ToLower(proxy.Host)
			if !strings.Contains(name, search) && !strings.Contains(host, search) {
				continue
			}
		}
		filtered = append(filtered, proxy)
	}
	return filtered, int64(len(filtered)), nil
}

func (s *stubAdminService) ListProxiesWithAccountCount(ctx context.Context, page, pageSize int, protocol, status, search string, sortBy, sortOrder string) ([]service.ProxyWithAccountCount, int64, error) {
	return s.proxyCounts, int64(len(s.proxyCounts)), nil
}

func (s *stubAdminService) GetAllProxies(ctx context.Context) ([]service.Proxy, error) {
	return s.proxies, nil
}

func (s *stubAdminService) GetAllProxiesWithAccountCount(ctx context.Context) ([]service.ProxyWithAccountCount, error) {
	return s.proxyCounts, nil
}

func (s *stubAdminService) GetProxy(ctx context.Context, id int64) (*service.Proxy, error) {
	for i := range s.proxies {
		proxy := s.proxies[i]
		if proxy.ID == id {
			return &proxy, nil
		}
	}
	proxy := service.Proxy{ID: id, Name: "proxy", Status: service.StatusActive}
	return &proxy, nil
}

func (s *stubAdminService) GetProxiesByIDs(ctx context.Context, ids []int64) ([]service.Proxy, error) {
	if len(ids) == 0 {
		return []service.Proxy{}, nil
	}
	out := make([]service.Proxy, 0, len(ids))
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		seen[id] = struct{}{}
	}
	for i := range s.proxies {
		proxy := s.proxies[i]
		if _, ok := seen[proxy.ID]; ok {
			out = append(out, proxy)
		}
	}
	return out, nil
}

func (s *stubAdminService) CreateProxy(ctx context.Context, input *service.CreateProxyInput) (*service.Proxy, error) {
	s.mu.Lock()
	s.createdProxies = append(s.createdProxies, input)
	s.mu.Unlock()
	proxy := service.Proxy{ID: 400, Name: input.Name, Status: service.StatusActive}
	return &proxy, nil
}

func (s *stubAdminService) UpdateProxy(ctx context.Context, id int64, input *service.UpdateProxyInput) (*service.Proxy, error) {
	s.mu.Lock()
	s.updatedProxyIDs = append(s.updatedProxyIDs, id)
	s.updatedProxies = append(s.updatedProxies, input)
	s.mu.Unlock()
	proxy := service.Proxy{ID: id, Name: input.Name, Status: service.StatusActive}
	return &proxy, nil
}

func (s *stubAdminService) DeleteProxy(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) BatchDeleteProxies(ctx context.Context, ids []int64) (*service.ProxyBatchDeleteResult, error) {
	return &service.ProxyBatchDeleteResult{DeletedIDs: ids}, nil
}

func (s *stubAdminService) GetProxyAccounts(ctx context.Context, proxyID int64) ([]service.ProxyAccountSummary, error) {
	return []service.ProxyAccountSummary{{ID: 1, Name: "account"}}, nil
}

func (s *stubAdminService) CheckProxyExists(ctx context.Context, host string, port int, username, password string) (bool, error) {
	return false, nil
}

func (s *stubAdminService) TestProxy(ctx context.Context, id int64) (*service.ProxyTestResult, error) {
	s.mu.Lock()
	s.testedProxyIDs = append(s.testedProxyIDs, id)
	s.mu.Unlock()
	return &service.ProxyTestResult{Success: true, Message: "ok"}, nil
}

func (s *stubAdminService) CheckProxyQuality(ctx context.Context, id int64) (*service.ProxyQualityCheckResult, error) {
	return &service.ProxyQualityCheckResult{
		ProxyID:        id,
		Score:          95,
		Grade:          "A",
		Summary:        "通过 5 项，告警 0 项，失败 0 项，挑战 0 项",
		PassedCount:    5,
		WarnCount:      0,
		FailedCount:    0,
		ChallengeCount: 0,
		CheckedAt:      time.Now().Unix(),
		Items: []service.ProxyQualityCheckItem{
			{Target: "base_connectivity", Status: "pass", Message: "ok"},
			{Target: "openai", Status: "pass", HTTPStatus: 401},
			{Target: "anthropic", Status: "pass", HTTPStatus: 401},
			{Target: "gemini", Status: "pass", HTTPStatus: 200},
		},
	}, nil
}

func (s *stubAdminService) ListRedeemCodes(ctx context.Context, page, pageSize int, codeType, status, search string, sortBy, sortOrder string) ([]service.RedeemCode, int64, error) {
	s.lastListRedeemCodes.codeType = codeType
	s.lastListRedeemCodes.status = status
	s.lastListRedeemCodes.search = search
	s.lastListRedeemCodes.sortBy = sortBy
	s.lastListRedeemCodes.sortOrder = sortOrder
	s.lastListRedeemCodes.calls++
	return s.redeems, int64(len(s.redeems)), nil
}

func (s *stubAdminService) GetRedeemCode(ctx context.Context, id int64) (*service.RedeemCode, error) {
	code := service.RedeemCode{ID: id, Code: "R-TEST", Status: service.StatusUnused}
	return &code, nil
}

func (s *stubAdminService) GenerateRedeemCodes(ctx context.Context, input *service.GenerateRedeemCodesInput) ([]service.RedeemCode, error) {
	return s.redeems, nil
}

func (s *stubAdminService) DeleteRedeemCode(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) BatchDeleteRedeemCodes(ctx context.Context, ids []int64) (int64, error) {
	return int64(len(ids)), nil
}

func (s *stubAdminService) ExpireRedeemCode(ctx context.Context, id int64) (*service.RedeemCode, error) {
	code := service.RedeemCode{ID: id, Code: "R-TEST", Status: service.StatusUsed}
	return &code, nil
}

func (s *stubAdminService) GetUserBalanceHistory(ctx context.Context, userID int64, page, pageSize int, codeType string) ([]service.RedeemCode, int64, float64, error) {
	return s.redeems, int64(len(s.redeems)), 100.0, nil
}

func (s *stubAdminService) UpdateGroupSortOrders(ctx context.Context, updates []service.GroupSortOrderUpdate) error {
	return nil
}

func (s *stubAdminService) AdminUpdateAPIKeyGroupID(ctx context.Context, keyID int64, groupID *int64) (*service.AdminUpdateAPIKeyGroupIDResult, error) {
	for i := range s.apiKeys {
		if s.apiKeys[i].ID == keyID {
			k := s.apiKeys[i]
			if groupID != nil {
				if *groupID == 0 {
					k.GroupID = nil
				} else {
					gid := *groupID
					k.GroupID = &gid
				}
			}
			return &service.AdminUpdateAPIKeyGroupIDResult{APIKey: &k}, nil
		}
	}
	return nil, service.ErrAPIKeyNotFound
}

func (s *stubAdminService) AdminResetAPIKeyRateLimitUsage(ctx context.Context, keyID int64) (*service.APIKey, error) {
	for i := range s.apiKeys {
		if s.apiKeys[i].ID == keyID {
			s.apiKeys[i].Usage5h = 0
			s.apiKeys[i].Usage1d = 0
			s.apiKeys[i].Usage7d = 0
			s.apiKeys[i].Window5hStart = nil
			s.apiKeys[i].Window1dStart = nil
			s.apiKeys[i].Window7dStart = nil
			k := s.apiKeys[i]
			return &k, nil
		}
	}
	return nil, service.ErrAPIKeyNotFound
}

func (s *stubAdminService) ResetAccountQuota(ctx context.Context, id int64) error {
	return nil
}

func (s *stubAdminService) EnsureOpenAIPrivacy(ctx context.Context, account *service.Account) string {
	return ""
}

func (s *stubAdminService) EnsureAntigravityPrivacy(ctx context.Context, account *service.Account) string {
	return ""
}

func (s *stubAdminService) ForceOpenAIPrivacy(ctx context.Context, account *service.Account) string {
	return ""
}

func (s *stubAdminService) ForceAntigravityPrivacy(ctx context.Context, account *service.Account) string {
	return ""
}

func (s *stubAdminService) ReplaceUserGroup(ctx context.Context, userID, oldGroupID, newGroupID int64) (*service.ReplaceUserGroupResult, error) {
	return &service.ReplaceUserGroupResult{MigratedKeys: 0}, nil
}

func (s *stubAdminService) RevertAccountProxyFallback(ctx context.Context, id int64) error {
	return nil
}

// Ensure stub implements interface.
var _ service.AdminService = (*stubAdminService)(nil)

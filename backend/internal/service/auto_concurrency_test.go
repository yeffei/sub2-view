package service

import (
	"context"
	"testing"
)

func TestAutoConcurrencyController_UsesBootstrapAndContractsOnSlowUpstream(t *testing.T) {
	t.Parallel()
	controller := newAutoConcurrencyController(nil)
	scope := AccountConcurrencyScope{AccountID: 101, AccountLimit: 1000}

	resolved, err := controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve scope: %v", err)
	}
	if resolved.AccountLimit != autoConcurrencyBootstrapAccount {
		t.Fatalf("bootstrap account limit = %d, want %d", resolved.AccountLimit, autoConcurrencyBootstrapAccount)
	}

	if err := controller.observe(context.Background(), scope, autoConcurrencySlowUpstreamMs); err != nil {
		t.Fatalf("observe slow upstream: %v", err)
	}
	resolved, err = controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve contracted scope: %v", err)
	}
	if resolved.AccountLimit != autoConcurrencyBootstrapAccount/2 {
		t.Fatalf("contracted account limit = %d, want %d", resolved.AccountLimit, autoConcurrencyBootstrapAccount/2)
	}
}

func TestAutoConcurrencyController_RecoversOnlyAfterHealthySamples(t *testing.T) {
	t.Parallel()
	controller := newAutoConcurrencyController(nil)
	scope := AccountConcurrencyScope{AccountID: 102, AccountLimit: 100}

	if err := controller.observe(context.Background(), scope, autoConcurrencySlowUpstreamMs); err != nil {
		t.Fatalf("observe slow upstream: %v", err)
	}
	for i := 0; i < autoConcurrencyRecoverSamples-1; i++ {
		if err := controller.observe(context.Background(), scope, 100); err != nil {
			t.Fatalf("observe healthy sample %d: %v", i, err)
		}
	}
	resolved, err := controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve pre-recovery scope: %v", err)
	}
	if resolved.AccountLimit != 10 {
		t.Fatalf("limit before recovery = %d, want 10", resolved.AccountLimit)
	}

	if err := controller.observe(context.Background(), scope, 100); err != nil {
		t.Fatalf("observe recovery sample: %v", err)
	}
	resolved, err = controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve recovered scope: %v", err)
	}
	if resolved.AccountLimit != 11 {
		t.Fatalf("limit after recovery = %d, want 11", resolved.AccountLimit)
	}
}

func TestAutoConcurrencyController_ContractsCapacityDomainWithAccount(t *testing.T) {
	t.Parallel()
	controller := newAutoConcurrencyController(nil)
	scope := AccountConcurrencyScope{
		AccountID:    103,
		AccountLimit: 1000,
		Capacity: &AccountCapacityScope{
			GroupID:         31,
			GroupLimit:      500,
			MemberHardLimit: 1000,
			MemberSoftShare: 1000,
		},
	}

	resolved, err := controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve capacity scope: %v", err)
	}
	if resolved.AccountLimit != 20 || resolved.Capacity.GroupLimit != 50 || resolved.Capacity.MemberHardLimit != 20 {
		t.Fatalf("bootstrap scope = %+v, want account=20 group=50 member=20", resolved)
	}

	if err := controller.observe(context.Background(), scope, autoConcurrencySlowUpstreamMs); err != nil {
		t.Fatalf("observe slow capacity upstream: %v", err)
	}
	resolved, err = controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve contracted capacity scope: %v", err)
	}
	if resolved.AccountLimit != 10 || resolved.Capacity.GroupLimit != 25 || resolved.Capacity.MemberHardLimit != 10 {
		t.Fatalf("contracted scope = %+v, want account=10 group=25 member=10", resolved)
	}
}

func TestAutoConcurrencyController_ContractsImmediatelyOnUpstreamFailure(t *testing.T) {
	t.Parallel()
	controller := newAutoConcurrencyController(nil)
	scope := AccountConcurrencyScope{AccountID: 104, AccountLimit: 100}

	if err := controller.contract(context.Background(), scope, "upstream_http_failure"); err != nil {
		t.Fatalf("contract after upstream failure: %v", err)
	}
	resolved, err := controller.resolveScope(context.Background(), scope)
	if err != nil {
		t.Fatalf("resolve contracted scope: %v", err)
	}
	if resolved.AccountLimit != 10 {
		t.Fatalf("failure-contracted account limit = %d, want 10", resolved.AccountLimit)
	}
}

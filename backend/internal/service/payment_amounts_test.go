package service

import "testing"

func TestCalculateCreditedBalanceWithCampaign(t *testing.T) {
	t.Parallel()

	t.Run("campaign disabled keeps base credited amount", func(t *testing.T) {
		t.Parallel()
		got, bonus := calculateCreditedBalanceWithCampaign(100, 1.2, false, 100, 10)
		if got != 120 {
			t.Fatalf("credited amount = %v, want 120", got)
		}
		if bonus != 0 {
			t.Fatalf("bonus amount = %v, want 0", bonus)
		}
	})

	t.Run("below threshold no bonus", func(t *testing.T) {
		t.Parallel()
		got, bonus := calculateCreditedBalanceWithCampaign(99.99, 1.2, true, 100, 10)
		if got != 119.99 {
			t.Fatalf("credited amount = %v, want 119.99", got)
		}
		if bonus != 0 {
			t.Fatalf("bonus amount = %v, want 0", bonus)
		}
	})

	t.Run("threshold reached adds bonus on credited amount", func(t *testing.T) {
		t.Parallel()
		got, bonus := calculateCreditedBalanceWithCampaign(100, 1.2, true, 100, 10)
		if got != 132 {
			t.Fatalf("credited amount = %v, want 132", got)
		}
		if bonus != 12 {
			t.Fatalf("bonus amount = %v, want 12", bonus)
		}
	})
}

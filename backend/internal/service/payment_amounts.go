package service

import (
	"math"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/shopspring/decimal"
)

const defaultBalanceRechargeMultiplier = 1.0
const defaultRechargeCampaignAmount = 100.0
const defaultRechargeCampaignBonusRate = 10.0

func normalizeBalanceRechargeMultiplier(multiplier float64) float64 {
	if math.IsNaN(multiplier) || math.IsInf(multiplier, 0) || multiplier <= 0 {
		return defaultBalanceRechargeMultiplier
	}
	return multiplier
}

func calculateCreditedBalance(paymentAmount, multiplier float64) float64 {
	return decimal.NewFromFloat(paymentAmount).
		Mul(decimal.NewFromFloat(normalizeBalanceRechargeMultiplier(multiplier))).
		Round(2).
		InexactFloat64()
}

func normalizeRechargeCampaignAmount(amount float64) float64 {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 {
		return defaultRechargeCampaignAmount
	}
	return amount
}

func normalizeRechargeCampaignBonusRate(rate float64) float64 {
	if math.IsNaN(rate) || math.IsInf(rate, 0) || rate < 0 {
		return defaultRechargeCampaignBonusRate
	}
	return rate
}

func calculateCreditedBalanceWithCampaign(paymentAmount, multiplier float64, campaignEnabled bool, campaignAmount, campaignBonusRate float64) (float64, float64) {
	baseCredited := calculateCreditedBalance(paymentAmount, multiplier)
	if !campaignEnabled {
		return baseCredited, 0
	}

	threshold := normalizeRechargeCampaignAmount(campaignAmount)
	rate := normalizeRechargeCampaignBonusRate(campaignBonusRate)
	if paymentAmount < threshold || rate <= 0 {
		return baseCredited, 0
	}

	bonusAmount := decimal.NewFromFloat(baseCredited).
		Mul(decimal.NewFromFloat(rate)).
		Div(decimal.NewFromFloat(100)).
		Round(2).
		InexactFloat64()
	return decimal.NewFromFloat(baseCredited).
		Add(decimal.NewFromFloat(bonusAmount)).
		Round(2).
		InexactFloat64(), bonusAmount
}

func calculateGatewayRefundAmount(orderAmount, payAmount, refundAmount float64, currency string) float64 {
	if orderAmount <= 0 || payAmount <= 0 || refundAmount <= 0 {
		return 0
	}
	fractionDigits := int32(payment.CurrencyMaxFractionDigits(currency))
	if math.Abs(refundAmount-orderAmount) <= paymentAmountToleranceForCurrency(currency) {
		return decimal.NewFromFloat(payAmount).Round(fractionDigits).InexactFloat64()
	}
	return decimal.NewFromFloat(payAmount).
		Mul(decimal.NewFromFloat(refundAmount)).
		Div(decimal.NewFromFloat(orderAmount)).
		Round(fractionDigits).
		InexactFloat64()
}

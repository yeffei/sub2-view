export interface UsageProfitSummary {
  revenue: number
  cost: number
  grossProfit: number
  grossMargin: number | null
  negative: boolean
}

const finiteAmount = (value: number | null | undefined) => (
  typeof value === 'number' && Number.isFinite(value) ? value : 0
)

export const calculateUsageProfit = (
  revenueValue: number | null | undefined,
  costValue: number | null | undefined,
): UsageProfitSummary => {
  const revenue = finiteAmount(revenueValue)
  const cost = finiteAmount(costValue)
  const grossProfit = revenue - cost
  return {
    revenue,
    cost,
    grossProfit,
    grossMargin: revenue > 0 ? grossProfit / revenue : null,
    negative: grossProfit < 0,
  }
}

export const formatGrossMargin = (margin: number | null, digits = 1) => (
  margin == null || !Number.isFinite(margin) ? '—' : `${(margin * 100).toFixed(digits)}%`
)

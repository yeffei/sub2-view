import { describe, expect, it } from 'vitest'

import { calculateUsageProfit, formatGrossMargin } from '../usageProfit'

describe('usageProfit', () => {
  it('calculates usage gross profit from sale value and account cost', () => {
    expect(calculateUsageProfit(12, 7)).toEqual({
      revenue: 12,
      cost: 7,
      grossProfit: 5,
      grossMargin: 5 / 12,
      negative: false,
    })
  })

  it('marks negative margin and avoids dividing by zero', () => {
    expect(calculateUsageProfit(2, 3).negative).toBe(true)
    expect(calculateUsageProfit(0, 3).grossMargin).toBeNull()
    expect(formatGrossMargin(null)).toBe('—')
  })
})

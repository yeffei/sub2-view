import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import GroupDistributionChart from '../GroupDistributionChart.vue'

const messages: Record<string, string> = {
  'admin.dashboard.groupDistribution': 'Group Distribution',
  'admin.dashboard.group': 'Group',
  'admin.dashboard.noGroup': 'No Group',
  'admin.dashboard.requests': 'Requests',
  'admin.dashboard.tokens': 'Tokens',
  'admin.dashboard.actual': 'Actual',
  'admin.dashboard.standard': 'Standard',
  'admin.dashboard.metricTokens': 'By Tokens',
  'admin.dashboard.metricActualCost': 'By Actual Cost',
  'admin.dashboard.metricGrossProfit': 'By Gross Profit',
  'admin.dashboard.inspectNegativeMargin': 'Inspect',
  'admin.dashboard.inspectNegativeMarginHint': 'Filter negative margin group',
  'admin.dashboard.accountCost': 'Cost',
  'usage.grossProfit': 'Gross Profit',
  'admin.dashboard.noDataAvailable': 'No data available',
}

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => messages[key] ?? key,
    }),
  }
})

vi.mock('vue-chartjs', () => ({
  Doughnut: {
    props: ['data'],
    template: '<div class="chart-data">{{ JSON.stringify(data) }}</div>',
  },
}))

describe('GroupDistributionChart', () => {
  const groupStats = [
    {
      group_id: 1,
      group_name: 'group-a',
      requests: 9,
      total_tokens: 1200,
      cost: 1.8,
      actual_cost: 0.1,
      account_cost: 0.08,
    },
    {
      group_id: 2,
      group_name: 'group-b',
      requests: 4,
      total_tokens: 600,
      cost: 0.7,
      actual_cost: 0.9,
      account_cost: 0.2,
    },
  ]

  const negativeGroupStats = [
    ...groupStats,
    {
      group_id: 3,
      group_name: 'loss-group',
      requests: 2,
      total_tokens: 300,
      cost: 0.4,
      actual_cost: 0.1,
      account_cost: 0.2,
    },
  ]

  it('uses total_tokens and token ordering by default', () => {
    const wrapper = mount(GroupDistributionChart, {
      props: {
        groupStats,
      },
      global: {
        stubs: {
          LoadingSpinner: true,
        },
      },
    })

    const chartData = JSON.parse(wrapper.find('.chart-data').text())
    expect(chartData.labels).toEqual(['group-a', 'group-b'])
    expect(chartData.datasets[0].data).toEqual([1200, 600])

    const rows = wrapper.findAll('tbody tr')
    expect(rows[0].text()).toContain('group-a')
    expect(rows[1].text()).toContain('group-b')

    const options = (wrapper.vm as any).$?.setupState.doughnutOptions
    const label = options.plugins.tooltip.callbacks.label({
      label: 'group-a',
      raw: 1200,
      dataset: { data: [1200, 600] },
    })
    expect(label).toBe('group-a: 1.20K (66.7%)')
  })

  it('uses actual_cost and reorders rows in actual cost mode', () => {
    const wrapper = mount(GroupDistributionChart, {
      props: {
        groupStats,
        metric: 'actual_cost',
      },
      global: {
        stubs: {
          LoadingSpinner: true,
        },
      },
    })

    const chartData = JSON.parse(wrapper.find('.chart-data').text())
    expect(chartData.labels).toEqual(['group-b', 'group-a'])
    expect(chartData.datasets[0].data).toEqual([0.9, 0.1])

    const rows = wrapper.findAll('tbody tr')
    expect(rows[0].text()).toContain('group-b')
    expect(rows[1].text()).toContain('group-a')

    const options = (wrapper.vm as any).$?.setupState.doughnutOptions
    const label = options.plugins.tooltip.callbacks.label({
      label: 'group-b',
      raw: 0.9,
      dataset: { data: [0.9, 0.1] },
    })
    expect(label).toBe('group-b: $0.900 (90.0%)')
  })

  it('uses gross profit and reorders rows in gross profit mode', () => {
    const wrapper = mount(GroupDistributionChart, {
      props: { groupStats: negativeGroupStats, metric: 'gross_profit' },
      global: { stubs: { LoadingSpinner: true } },
    })

    const chartData = JSON.parse(wrapper.find('.chart-data').text())
    expect(chartData.labels).toEqual(['group-b', 'group-a', 'loss-group'])
    expect(chartData.datasets[0].data[0]).toBeCloseTo(0.7)
    expect(chartData.datasets[0].data[1]).toBeCloseTo(0.02)
    expect(chartData.datasets[0].data[2]).toBe(0)
  })

  it('emits a focused inspection action for negative-margin groups', async () => {
    const wrapper = mount(GroupDistributionChart, {
      props: { groupStats: negativeGroupStats, metric: 'gross_profit' },
      global: { stubs: { LoadingSpinner: true } },
    })

    const inspectButton = wrapper.get('button[title="Filter negative margin group"]')
    await inspectButton.trigger('click')
    expect(wrapper.emitted('inspectGroup')).toEqual([[3]])
  })
})

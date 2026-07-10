import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import EndpointDistributionChart from '../EndpointDistributionChart.vue'

const messages: Record<string, string> = {
  'usage.endpointDistribution': 'Endpoint Distribution',
  'usage.endpoint': 'Endpoint',
  'usage.inbound': 'Inbound',
  'usage.upstream': 'Upstream',
  'usage.path': 'Path',
  'admin.dashboard.requests': 'Requests',
  'admin.dashboard.tokens': 'Tokens',
  'admin.dashboard.actual': 'Actual',
  'admin.dashboard.standard': 'Standard',
  'admin.dashboard.metricTokens': 'By Tokens',
  'admin.dashboard.metricActualCost': 'By Actual Cost',
  'admin.dashboard.metricCacheHitRate': 'By Hit Rate',
  'admin.dashboard.metricCacheReadPerHit': 'By Read/Hit',
  'admin.dashboard.noDataAvailable': 'No data available',
  'admin.dashboard.cacheHitRequestsShort': 'Hits {value}',
  'admin.dashboard.cacheHitRateShort': 'Hit rate {value}',
  'admin.dashboard.cacheReadPerHitShort': 'Read/hit {value}',
  'admin.dashboard.avgInputShort': 'Avg in {value}',
}

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string | number>) => {
        const message = messages[key] ?? key
        if (!params) return message
        return Object.entries(params).reduce(
          (result, [name, value]) => result.replace(`{${name}}`, String(value)),
          message
        )
      },
    }),
  }
})

vi.mock('vue-chartjs', () => ({
  Doughnut: {
    props: ['data'],
    template: '<div class="chart-data">{{ JSON.stringify(data) }}</div>',
  },
}))

describe('EndpointDistributionChart', () => {
  it('shows cache helper metrics for endpoint rows', () => {
    const wrapper = mount(EndpointDistributionChart, {
      props: {
        endpointStats: [
          {
            endpoint: '/v1/messages',
            requests: 10,
            input_tokens: 300,
            output_tokens: 200,
            cache_creation_tokens: 40,
            cache_read_tokens: 120,
            cache_read_hit_requests: 4,
            cache_creation_requests: 2,
            cache_read_hit_ratio: 0.4,
            average_cache_read_tokens_per_hit: 30,
            average_actual_input_tokens: 30,
            total_tokens: 660,
            cost: 2.1,
            actual_cost: 1.8,
          },
        ],
      },
      global: {
        stubs: {
          LoadingSpinner: true,
          UserBreakdownSubTable: true,
        },
      },
    })

    const text = wrapper.text()
    expect(text).toContain('/v1/messages')
    expect(text).toContain('Hits 4')
    expect(text).toContain('Hit rate 40.0%')
    expect(text).toContain('Read/hit 30')
    expect(text).toContain('Avg in 30')
  })

  it('uses cache_hit_ratio metric for chart values', () => {
    const wrapper = mount(EndpointDistributionChart, {
      props: {
        metric: 'cache_hit_ratio',
        endpointStats: [
          {
            endpoint: '/v1/messages',
            requests: 10,
            input_tokens: 300,
            output_tokens: 200,
            cache_creation_tokens: 40,
            cache_read_tokens: 120,
            cache_read_hit_requests: 4,
            cache_creation_requests: 2,
            cache_read_hit_ratio: 0.4,
            average_cache_read_tokens_per_hit: 30,
            average_actual_input_tokens: 30,
            total_tokens: 660,
            cost: 2.1,
            actual_cost: 1.8,
          },
          {
            endpoint: '/v1/responses',
            requests: 8,
            input_tokens: 250,
            output_tokens: 180,
            cache_creation_tokens: 12,
            cache_read_tokens: 20,
            cache_read_hit_requests: 1,
            cache_creation_requests: 1,
            cache_read_hit_ratio: 0.125,
            average_cache_read_tokens_per_hit: 20,
            average_actual_input_tokens: 31.25,
            total_tokens: 462,
            cost: 1.5,
            actual_cost: 1.2,
          },
        ],
      },
      global: {
        stubs: {
          LoadingSpinner: true,
          UserBreakdownSubTable: true,
        },
      },
    })

    const chartData = JSON.parse(wrapper.find('.chart-data').text())
    expect(chartData.labels).toEqual(['/v1/messages', '/v1/responses'])
    expect(chartData.datasets[0].data).toEqual([0.4, 0.125])
  })
})

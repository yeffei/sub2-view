import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import UserErrorDetailModal from '@/components/user/UserErrorDetailModal.vue'

const {
  getMyErrorDetailMock,
  listMyErrorRequestsMock,
  queryUsageMock,
  showSuccessMock,
  showErrorMock,
  routerPushMock,
  writeTextMock,
} = vi.hoisted(() => ({
  getMyErrorDetailMock: vi.fn(),
  listMyErrorRequestsMock: vi.fn(),
  queryUsageMock: vi.fn(),
  showSuccessMock: vi.fn(),
  showErrorMock: vi.fn(),
  routerPushMock: vi.fn(),
  writeTextMock: vi.fn(),
}))

const messages: Record<string, string | string[]> = {
  'usage.errors.detail.title': 'Error Request Detail',
  'usage.errors.detail.loadFailed': 'Failed to load detail',
  'usage.errors.time': 'Time',
  'usage.errors.model': 'Model',
  'usage.errors.endpoint': 'Endpoint',
  'usage.errors.status': 'Status',
  'usage.errors.category': 'Category',
  'usage.errors.platform': 'Platform',
  'usage.errors.message': 'Message',
  'usage.errors.detail.upstreamStatus': 'Upstream Status',
  'usage.errors.detail.explanationTitle': 'What happened',
  'usage.errors.detail.copySummary': 'Copy diagnostic summary',
  'usage.errors.detail.copySummarySuccess': 'Diagnostic summary copied',
  'usage.errors.detail.copySummaryFailed': 'Failed to copy diagnostic summary',
  'usage.errors.detail.thisRequest': 'this request',
  'usage.errors.detail.summaryLabels.title': 'Error request diagnostic summary',
  'usage.errors.detail.summaryLabels.explanation': 'Explanation',
  'usage.errors.detail.summaryLabels.advice': 'Advice',
  'usage.errors.detail.summaryLabels.timeline': 'Timeline conclusion',
  'usage.errors.detail.summaryLabels.nextAction': 'Next action',
  'usage.errors.detail.actions.quota': 'Review account and quota',
  'usage.errors.detail.actions.invalid_request': 'Review this request in usage',
  'usage.errors.detail.recoveryGuide.title': 'Self-service recovery guide',
  'usage.errors.detail.recoveryGuide.badge': 'Three-step retest',
  'usage.errors.detail.recoveryGuide.checksTitle': 'Confirm after recovery',
  'usage.errors.detail.recoveryGuide.anchor': 'Related action: {label}',
  'usage.errors.detail.recoveryGuide.summaries.quota': 'For quota issues, confirm balance or plan state has recovered before retrying the request.',
  'usage.errors.detail.recoveryGuide.summaries.invalid_request': 'For request-shape issues, fix the payload first and then run a narrow retest.',
  'usage.errors.detail.recoveryGuide.steps.quota.checkBalance.title': 'Confirm balance or plan state first',
  'usage.errors.detail.recoveryGuide.steps.quota.checkBalance.detail': 'Check whether your balance, plan, or subscription quota has already recovered to a usable state.',
  'usage.errors.detail.recoveryGuide.steps.quota.confirmRefresh.title': 'Wait for state refresh to finish',
  'usage.errors.detail.recoveryGuide.steps.quota.confirmRefresh.detail': 'If you just recharged, renewed, or changed plans, wait for the new state to propagate before sending repeated attempts.',
  'usage.errors.detail.recoveryGuide.steps.quota.retest.title': 'Run one low-volume retest',
  'usage.errors.detail.recoveryGuide.steps.quota.retest.detail': 'Once the state is back, send one small request first and confirm the quota error has disappeared before restoring normal traffic.',
  'usage.errors.detail.recoveryGuide.steps.invalid_request.reviewPayload.title': 'Review the failing request first',
  'usage.errors.detail.recoveryGuide.steps.invalid_request.reviewPayload.detail': 'Go back to usage and check whether the model, endpoint, and request fields still belong together.',
  'usage.errors.detail.recoveryGuide.steps.invalid_request.matchModel.title': 'Confirm the model matches the request style',
  'usage.errors.detail.recoveryGuide.steps.invalid_request.matchModel.detail': 'If you recently switched models or endpoints, confirm the new route still supports the current payload, tool, or attachment shape.',
  'usage.errors.detail.recoveryGuide.steps.invalid_request.retest.title': 'Retry with the smallest payload',
  'usage.errors.detail.recoveryGuide.steps.invalid_request.retest.detail': 'Remove suspicious fields and retest with the smallest payload so you can isolate which parameter causes the failure.',
  'usage.errors.detail.recoveryGuide.checks.invalid_request.statusNormal': 'The retried request returns a normal status again.',
  'usage.errors.detail.recoveryGuide.checks.invalid_request.noSameCategory': 'The request no longer lands in the invalid-request category.',
  'usage.errors.detail.recoveryGuide.checks.quota.balanceVisible': 'Balance or quota state is visible again.',
  'usage.errors.detail.recoveryGuide.checks.quota.noNewQuotaErrors': 'The retry no longer returns quota or balance exhaustion.',
  'usage.errors.categories.quota': 'Balance/Subscription',
  'usage.errors.categories.rate_limit': 'Rate limited',
  'usage.errors.detail.explanations.quota.summary': '{endpoint} could not continue because balance or quota was exhausted.',
  'usage.errors.detail.explanations.quota.advice': [
    'Check your balance or remaining quota.',
    'Retry after the balance state refreshes.',
  ],
  'usage.errors.detail.reasonExplanations.request_model_not_supported.summary': '{endpoint} looks more like a model-name or model-access mismatch: the requested model is not actually usable on this route as configured now.',
  'usage.errors.detail.reasonExplanations.request_model_not_supported.advice': [
    'Call the models list first and verify which model IDs are really available for the current key and group.',
    'Then check whether the client model name is outdated, misspelled, or not enabled for the current group.',
  ],
  'usage.errors.detail.reasonExplanations.service_model_not_available.summary': '{endpoint} currently has no usable account on this route for the requested model, so this looks more like group or routing coverage missing that model than a short-lived spike.',
  'usage.errors.detail.reasonExplanations.service_model_not_available.advice': [
    'Call the models list first and verify which model IDs are actually open to the current key and group.',
    'If you recently changed group, channel, or model name, confirm that this route truly includes that model.',
  ],
  'usage.errors.detail.modelTrace.title': 'Model path',
  'usage.errors.detail.modelTrace.requested': 'Requested model',
  'usage.errors.detail.modelTrace.upstream': 'Upstream model',
  'usage.errors.detail.modelTrace.mapped': 'This request entered as {requested}, but the upstream call was sent as {upstream}. If that is not what you expected, check the current group mapping and client-side model name first.',
  'usage.errors.detail.modelTrace.unavailable': 'This failure is centered on requested model {requested}, which suggests the issue is on that model route itself rather than on all requests globally.',
  'usage.errors.detail.modelTrace.requestedOnly': 'This record only confirms requested model {requested}; if you expected another model path, check the client configuration and current group first.',
  'usage.errors.detail.timeline.title': 'Error replay timeline',
  'usage.errors.detail.timeline.windowLabel': 'Nearby events on the same day',
  'usage.errors.detail.timeline.loading': 'Building the request timeline around this failure…',
  'usage.errors.detail.timeline.loadFailed': 'Failed to build the timeline for now. Please try again later.',
  'usage.errors.detail.timeline.empty': 'There are not enough nearby records to build a timeline yet.',
  'usage.errors.detail.timeline.currentTitle': 'Current failed request',
  'usage.errors.detail.timeline.successTitle': 'Nearby successful request',
  'usage.errors.detail.timeline.errorTitle': 'Nearby failure: {category}',
  'usage.errors.detail.timeline.badges.current': 'Current error',
  'usage.errors.detail.timeline.badges.success': 'Recovered',
  'usage.errors.detail.timeline.badges.error': 'Nearby failure',
  'usage.errors.detail.timeline.summaries.recovered': 'A successful request appeared shortly after this failure, so this looks more like a transient spike or brief congestion.',
  'usage.errors.detail.timeline.summaries.continuous': 'There are multiple nearby failures around this request, which looks more like a continuing issue than a one-off retry case.',
  'usage.errors.detail.timeline.summaries.isolated': 'There were successful requests before this failure and no dense cluster of nearby failures, so this looks closer to an isolated event.',
  'usage.errors.detail.timeline.summaries.sparse': 'Nearby records are sparse, so the timeline can only confirm this failure itself for now.',
  'usage.errors.detail.timeline.summaries.observe': 'Some nearby events were found, but there is still not enough evidence to make a clear attribution yet.',
}

vi.mock('@/api/usage', () => ({
  getMyErrorDetail: getMyErrorDetailMock,
  listMyErrorRequests: listMyErrorRequestsMock,
  query: queryUsageMock,
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showSuccess: showSuccessMock,
    showError: showErrorMock,
  }),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: routerPushMock,
  }),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string>) => {
        const value = messages[key]
        if (Array.isArray(value)) return value
        const message = value ?? key
        if (!params) return message
        return Object.entries(params).reduce(
          (result, [name, replacement]) => result.replaceAll(`{${name}}`, replacement),
          message,
        )
      },
    }),
  }
})

describe('UserErrorDetailModal', () => {
  beforeEach(() => {
    getMyErrorDetailMock.mockReset()
    listMyErrorRequestsMock.mockReset()
    queryUsageMock.mockReset()
    showSuccessMock.mockReset()
    showErrorMock.mockReset()
    routerPushMock.mockReset()
    writeTextMock.mockReset()
    writeTextMock.mockResolvedValue(undefined)

    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: {
        writeText: writeTextMock,
      },
    })
  })

  it('copies a structured diagnostic summary from the loaded detail', async () => {
    getMyErrorDetailMock.mockResolvedValue({
      id: 8,
      created_at: '2026-06-23T12:00:00Z',
      model: 'gpt-5',
      inbound_endpoint: '/v1/chat/completions',
      status_code: 402,
      category: 'quota',
      platform: 'openai',
      message: 'quota exceeded',
      key_name: 'demo',
      key_deleted: false,
      error_body: '{"error":"quota exceeded"}',
      upstream_status_code: 429,
    })
    queryUsageMock.mockResolvedValue({
      items: [
        {
          id: 101,
          created_at: '2026-06-23T12:03:00Z',
          model: 'gpt-5',
          inbound_endpoint: '/v1/chat/completions',
          input_tokens: 120,
          output_tokens: 30,
          actual_cost: 0.015,
        },
      ],
    })
    listMyErrorRequestsMock.mockResolvedValue({
      items: [
        {
          id: 9,
          created_at: '2026-06-23T11:58:00Z',
          model: 'gpt-5',
          inbound_endpoint: '/v1/chat/completions',
          status_code: 429,
          category: 'rate_limit',
          platform: 'openai',
          message: 'too many requests',
          key_name: 'demo',
          key_deleted: false,
        },
      ],
    })

    const wrapper = mount(UserErrorDetailModal, {
      props: {
        show: false,
        errorId: null,
      },
      global: {
        stubs: {
          BaseDialog: {
            props: ['show', 'title'],
            template: '<div v-if="show"><slot></slot></div>',
          },
        },
      },
    })

    await wrapper.setProps({ show: true, errorId: 8 })
    await flushPromises()

    const buttons = wrapper.findAll('button')
    const copyButton = buttons.find((button) => button.text().includes('Copy diagnostic summary'))
    expect(copyButton).toBeTruthy()

    await copyButton!.trigger('click')

    expect(writeTextMock).toHaveBeenCalledTimes(1)
    const copiedText = writeTextMock.mock.calls[0][0] as string
    expect(copiedText).toContain('Error request diagnostic summary')
    expect(copiedText).toContain('Time: ')
    expect(copiedText).toContain('06/23/2026')
    expect(copiedText).toContain('Model: gpt-5')
    expect(copiedText).toContain('Endpoint: /v1/chat/completions')
    expect(copiedText).toContain('Status: 402')
    expect(copiedText).toContain('Category: Balance/Subscription')
    expect(copiedText).toContain('Platform: openai')
    expect(copiedText).toContain('Upstream Status: 429')
    expect(copiedText).toContain('Message: quota exceeded')
    expect(copiedText).toContain('Explanation: /v1/chat/completions could not continue because balance or quota was exhausted.')
    expect(copiedText).toContain('Advice: Check your balance or remaining quota.；Retry after the balance state refreshes.')
    expect(copiedText).toContain('Timeline conclusion: A successful request appeared shortly after this failure, so this looks more like a transient spike or brief congestion.')
    expect(copiedText).toContain('Next action: Review account and quota')
    expect(showSuccessMock).toHaveBeenCalledWith('Diagnostic summary copied')
    expect(showErrorMock).not.toHaveBeenCalled()
  })

  it('renders a quota recovery guide with ordered steps and checks', async () => {
    getMyErrorDetailMock.mockResolvedValue({
      id: 8,
      created_at: '2026-06-23T12:00:00Z',
      model: 'gpt-5',
      inbound_endpoint: '/v1/chat/completions',
      status_code: 402,
      category: 'quota',
      platform: 'openai',
      message: 'quota exceeded',
      key_name: 'demo',
      key_deleted: false,
      error_body: '{"error":"quota exceeded"}',
      upstream_status_code: 429,
    })
    queryUsageMock.mockResolvedValue({ items: [] })
    listMyErrorRequestsMock.mockResolvedValue({ items: [] })

    const wrapper = mount(UserErrorDetailModal, {
      props: {
        show: false,
        errorId: null,
      },
      global: {
        stubs: {
          BaseDialog: {
            props: ['show', 'title'],
            template: '<div v-if="show"><slot></slot></div>',
          },
        },
      },
    })

    await wrapper.setProps({ show: true, errorId: 8 })
    await flushPromises()

    const text = wrapper.text()
    expect(text).toContain('Self-service recovery guide')
    expect(text).toContain('For quota issues, confirm balance or plan state has recovered before retrying the request.')
    expect(text).toContain('Confirm balance or plan state first')
    expect(text).toContain('Related action: Review account and quota')
    expect(text).toContain('Confirm after recovery')
    expect(text).toContain('The retry no longer returns quota or balance exhaustion.')
  })

  it('prefers diagnosis-specific explanation when the backend returns a reason code', async () => {
    getMyErrorDetailMock.mockResolvedValue({
      id: 18,
      created_at: '2026-06-23T12:00:00Z',
      model: 'gpt-5.5',
      inbound_endpoint: '/v1/chat/completions',
      status_code: 404,
      category: 'invalid_request',
      platform: 'openai',
      message: 'model not found',
      key_name: 'demo',
      key_deleted: false,
      error_body: '{"error":"model not found"}',
      diagnosis: {
        reason_code: 'request_model_not_supported',
        action_code: 'usage_review_payload',
        requested_model: 'gpt-5.5',
        upstream_model: 'gpt-5.5-mini',
      },
    })
    queryUsageMock.mockResolvedValue({ items: [] })
    listMyErrorRequestsMock.mockResolvedValue({ items: [] })

    const wrapper = mount(UserErrorDetailModal, {
      props: {
        show: false,
        errorId: null,
      },
      global: {
        stubs: {
          BaseDialog: {
            props: ['show', 'title'],
            template: '<div v-if="show"><slot></slot></div>',
          },
        },
      },
    })

    await wrapper.setProps({ show: true, errorId: 18 })
    await flushPromises()

    const text = wrapper.text()
    expect(text).toContain('model-name or model-access mismatch')
    expect(text).toContain('Call the models list first')
    expect(text).toContain('Requested model')
    expect(text).toContain('Upstream model')
    expect(text).toContain('This request entered as gpt-5.5, but the upstream call was sent as gpt-5.5-mini')
    expect(text).toContain('Review this request in usage')
  })

  it('explains route-level model unavailability with requested model evidence', async () => {
    getMyErrorDetailMock.mockResolvedValue({
      id: 28,
      created_at: '2026-06-23T12:05:00Z',
      model: 'gpt-5.5',
      inbound_endpoint: '/v1/chat/completions',
      status_code: 503,
      category: 'service_unavailable',
      platform: 'openai',
      message: 'no available accounts supporting model',
      key_name: 'demo',
      key_deleted: false,
      error_body: '{"error":"no available accounts"}',
      diagnosis: {
        reason_code: 'service_model_not_available',
        action_code: 'usage_retry_later',
        requested_model: 'gpt-5.5',
      },
    })
    queryUsageMock.mockResolvedValue({ items: [] })
    listMyErrorRequestsMock.mockResolvedValue({ items: [] })

    const wrapper = mount(UserErrorDetailModal, {
      props: {
        show: false,
        errorId: null,
      },
      global: {
        stubs: {
          BaseDialog: {
            props: ['show', 'title'],
            template: '<div v-if="show"><slot></slot></div>',
          },
        },
      },
    })

    await wrapper.setProps({ show: true, errorId: 28 })
    await flushPromises()

    const text = wrapper.text()
    expect(text).toContain('no usable account on this route for the requested model')
    expect(text).toContain('Requested model')
    expect(text).toContain('gpt-5.5')
    expect(text).toContain('This failure is centered on requested model gpt-5.5')
  })
})

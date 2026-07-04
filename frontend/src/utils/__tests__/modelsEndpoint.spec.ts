import { describe, expect, it } from 'vitest'
import { extractErrorTextFromPayload, extractModelIdsFromModelsPayload, safeParseJson } from '@/utils/modelsEndpoint'

describe('modelsEndpoint helpers', () => {
  it('parses valid json safely', () => {
    expect(safeParseJson('{"data":[{"id":"gpt-5.5"}]}')).toEqual({ data: [{ id: 'gpt-5.5' }] })
    expect(safeParseJson('')).toBeNull()
    expect(safeParseJson('not-json')).toBeNull()
  })

  it('extracts unique model ids from /v1/models payload', () => {
    expect(
      extractModelIdsFromModelsPayload({
        object: 'list',
        data: [
          { id: 'gpt-5.5' },
          { id: 'gpt-5.5-mini' },
          { id: 'gpt-5.5' },
          { id: '' },
        ],
      }),
    ).toEqual(['gpt-5.5', 'gpt-5.5-mini'])
  })

  it('extracts a user-facing error text from common error payloads', () => {
    expect(extractErrorTextFromPayload({ error: { message: 'invalid api key' } })).toBe('invalid api key')
    expect(extractErrorTextFromPayload({ message: 'forbidden' })).toBe('forbidden')
    expect(extractErrorTextFromPayload({ error: { detail: 'temporarily unavailable' } })).toBe('temporarily unavailable')
  })
})

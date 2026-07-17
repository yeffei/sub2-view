package service

import "context"

// HandleOpenAIPoolModeHighTTFT delegates the global high-TTFT circuit to the
// existing rate-limit service, keeping scheduling state in one place.
func (s *OpenAIGatewayService) HandleOpenAIPoolModeHighTTFT(ctx context.Context, account *Account, firstTokenMs *int) bool {
	if s == nil || s.rateLimitService == nil {
		return false
	}
	stateCtx, cancel := openAIAccountStateContext(ctx)
	defer cancel()
	return s.rateLimitService.HandleOpenAIPoolModeHighTTFT(stateCtx, account, firstTokenMs)
}

// ReportOpenAIUpstreamPostWriteWait feeds the adaptive concurrency controller
// with only the interval after the upstream received the request body. This
// keeps local slot wait time and client upload time out of the feedback loop.
func (s *OpenAIGatewayService) ReportOpenAIUpstreamPostWriteWait(ctx context.Context, account *Account, postWriteWaitMs int64) {
	if s == nil || s.concurrencyService == nil || account == nil || postWriteWaitMs <= 0 {
		return
	}
	if err := s.concurrencyService.ReportUpstreamPostWriteWait(ctx, account.ConcurrencyScope(), postWriteWaitMs); err != nil {
		return
	}
}

func (s *OpenAIGatewayService) ReportOpenAIUpstreamFailure(ctx context.Context, account *Account, statusCode int) {
	if s == nil || s.concurrencyService == nil || account == nil || (statusCode != 429 && statusCode < 500) {
		return
	}
	_ = s.concurrencyService.ReportUpstreamFailure(ctx, account.ConcurrencyScope(), "upstream_http_failure")
}

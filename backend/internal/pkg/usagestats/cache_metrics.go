package usagestats

func safeUsageRatio(numerator, denominator int64) float64 {
	if denominator <= 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

func safeUsageAverage(total, count int64) float64 {
	if count <= 0 {
		return 0
	}
	return float64(total) / float64(count)
}

func (s *UsageStats) PopulateCacheDerivedMetrics() {
	if s == nil {
		return
	}
	s.CacheReadHitRatio = safeUsageRatio(s.CacheReadHitRequests, s.TotalRequests)
	s.AverageCacheReadTokensPerHit = safeUsageAverage(s.TotalCacheReadTokens, s.CacheReadHitRequests)
	s.AverageActualInputTokens = safeUsageAverage(s.TotalInputTokens, s.TotalRequests)
}

func (s *ModelStat) PopulateCacheDerivedMetrics() {
	if s == nil {
		return
	}
	s.CacheReadHitRatio = safeUsageRatio(s.CacheReadHitRequests, s.Requests)
	s.AverageCacheReadTokensPerHit = safeUsageAverage(s.CacheReadTokens, s.CacheReadHitRequests)
	s.AverageActualInputTokens = safeUsageAverage(s.InputTokens, s.Requests)
}

func (s *EndpointStat) PopulateCacheDerivedMetrics() {
	if s == nil {
		return
	}
	s.CacheReadHitRatio = safeUsageRatio(s.CacheReadHitRequests, s.Requests)
	s.AverageCacheReadTokensPerHit = safeUsageAverage(s.CacheReadTokens, s.CacheReadHitRequests)
	s.AverageActualInputTokens = safeUsageAverage(s.InputTokens, s.Requests)
}

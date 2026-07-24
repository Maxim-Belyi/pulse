package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"pulse/pkg/logger"
	"time"
)

type AnalyticsUseCase struct {
	logger *logger.Logger
	repo   AnalyticsRepository
	cache  CacheRepository
}

func NewAnalyticsUseCase(logger *logger.Logger, repo AnalyticsRepository, cache CacheRepository) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

func (a *AnalyticsUseCase) GetTopSources(ctx context.Context, since time.Time) ([]SourceStats, error) {
	cacheKey := fmt.Sprintf("top_sources:%d", since.Unix())

	cachedBytes, err := a.cache.Get(ctx, cacheKey)
	if err == nil {
		var result []SourceStats
		json.Unmarshal(cachedBytes, &result)
		return result, nil
	}
	a.logger.Info("кеша GetTopSources нет или он недоступен") //TODO убрать

	result, err := a.repo.GetTopSources(ctx, since)
	if err != nil {
		a.logger.Error(err, "GetTopSources: ошибка бд")
		return nil, err
	}

	bytesToCache, err := json.Marshal(result)
	if err == nil {
		a.cache.Set(ctx, cacheKey, bytesToCache, 5 * time.Minute) 
	}
	return result, nil
}

func (a *AnalyticsUseCase) GetHourlyTrends(ctx context.Context, since time.Time) ([]TrendStat,error) {
	hourlyTrends := fmt.Sprintf("top_hour_sources:%d", since.Unix())

	cachedBytes, err := a.cache.Get(ctx, hourlyTrends)
	if err == nil {
		var result []TrendStat
		json.Unmarshal(cachedBytes, &result)
		return result, nil
	}
	a.logger.Info("кеша GetHourlyTrends нет или он не доступен") //TODO убрать

	result, err := a.repo.GetHourlyTrends(ctx, since)
	if err != nil {
		a.logger.Error(err, "GetHourlyTrends: ошибка бд")
		return nil, err
	}

	bytesCache, err := json.Marshal(result)
	if err == nil {
		a.cache.Set(ctx, hourlyTrends, bytesCache, 5* time.Second)
	}
	return result, nil
}
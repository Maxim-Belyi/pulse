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
		if err := json.Unmarshal(cachedBytes, &result); err == nil {
			return result, nil
		} else {
			a.logger.Error(err, "ошибка десереализации кеша GetTopSources")
		}
	}
	a.logger.Info("кеша GetTopSources нет или он недоступен") 

	result, err := a.repo.GetTopSources(ctx, since)
	if err != nil {
		a.logger.Error(err, "GetTopSources: ошибка бд")
		return nil, fmt.Errorf("AnalyticsUseCase - GetTopSources: %w", err)
	}

	bytesToCache, err := json.Marshal(result)
	if err == nil {
		a.cache.Set(ctx, cacheKey, bytesToCache, time.Hour) 
	}
	return result, nil
}

func (a *AnalyticsUseCase) GetHourlyTrends(ctx context.Context, since time.Time) ([]TrendStat,error) {
	hourlyTrends := fmt.Sprintf("top_hour_sources:%d", since.Unix())

	cachedBytes, err := a.cache.Get(ctx, hourlyTrends)
	if err == nil {
		var result []TrendStat
		if err := json.Unmarshal(cachedBytes, &result); err == nil {
			return result,nil
		} else {
			a.logger.Error(err, "ошибка десереализации кеша GetHourlyTrends")
		}
	}
	a.logger.Info("кеша GetHourlyTrends нет или он не доступен") 

	result, err := a.repo.GetHourlyTrends(ctx, since)
	if err != nil {
		a.logger.Error(err, "GetHourlyTrends: ошибка бд")
		return nil, fmt.Errorf("AnalyticsUseCase - GetHourlyTrends: %w", err)
	}

	bytesCache, err := json.Marshal(result)
	if err == nil {
		a.cache.Set(ctx, hourlyTrends, bytesCache, time.Hour)
	}
	return result, nil
}
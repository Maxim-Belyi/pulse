package clickhouse

import (
	"context"
	"fmt"
	"pulse/internal/usecase"
	"pulse/pkg/logger"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type AnalyticsRepository struct {
	logger *logger.Logger
	dbConn driver.Conn
}

func NewAnalyticsRepository(logger *logger.Logger, dbConn driver.Conn) *AnalyticsRepository {
	return &AnalyticsRepository{
		logger: logger,
		dbConn: dbConn,
	}
}

func (a *AnalyticsRepository) GetTopSources(ctx context.Context, since time.Time) ([]usecase.SourceStats, error) {
	query := `SELECT source,
			  count() FROM events
			  WHERE occurred_at >= ?
			  GROUP BY source`

	rows, err := a.dbConn.Query(ctx, query, since)
	if err != nil {
		a.logger.Error(err, "ошибка подключения к бд GetTopSources")
		return nil, fmt.Errorf("db conn scan: %w", err)
	}
	defer rows.Close()
	result := make([]usecase.SourceStats, 0)

	for rows.Next() {
		var sourceName string
		var total uint64

		err := rows.Scan(&sourceName, &total)
		if err != nil {
			return nil, fmt.Errorf("getTopResources scan: %w", err)
		}
		stat := usecase.SourceStats{
			Source:      sourceName,
			TotalEvents: total,
		}
		result = append(result, stat)
	}

	if err = rows.Err(); err != nil {
		a.logger.Error(err, "ошибка при итерации строк")
		return nil, fmt.Errorf("getTopResources iteration: %w", err)
	}

	return result, nil
}

func (a *AnalyticsRepository) GetHourlyTrends(ctx context.Context, since time.Time) ([]usecase.TrendStat, error) {
	query := `SELECT 
    		  toStartOfHour(occurred_at) AS hour_bucket, 
    		  count() AS total_events
			  FROM events
			  WHERE occurred_at >= ?
			  GROUP BY hour_bucket
			  ORDER BY hour_bucket ASC WITH FILL STEP 3600`

	rows, err := a.dbConn.Query(ctx, query, since)
	if err != nil {
		a.logger.Error(err, "ошибка подключения к бд")
		return nil, fmt.Errorf("db conn scan: %w", err)
	}
	defer rows.Close()
	result := make([]usecase.TrendStat, 0)

	for rows.Next() {
		var hourBucket time.Time
		var totalEvents uint64

		err := rows.Scan(&hourBucket, &totalEvents)
		if err != nil {
			return nil, fmt.Errorf("GetHourlyTrends scan: %w", err)
		}

		stat := usecase.TrendStat{
			HourBucket:  hourBucket,
			TotalEvents: totalEvents,
		}
		result = append(result, stat)
	}

	if err = rows.Err(); err != nil {
		a.logger.Error(err, "ошибка при итерации строк")
		return nil, fmt.Errorf("GetHourlyTrends iteration: %w", err)
	}
	return result, nil
}

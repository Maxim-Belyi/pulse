package grpc

import (
	"context"
	pb "pulse/api/pb/analytics_v1"
	"pulse/internal/usecase"
	"pulse/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsController struct {
	logger  *logger.Logger
	useCase *usecase.AnalyticsUseCase
}

func NewAnalyticsController(logger *logger.Logger, useCase *usecase.AnalyticsUseCase) *AnalyticsController {
	return &AnalyticsController{
		logger:  logger,
		useCase: useCase,
	}
}

func (a *AnalyticsController) GetTopSources(ctx context.Context, req *pb.GetTopSourcesRequest) (*pb.SourceResponse, error) {
	parsedTime := req.Since.AsTime()
	stats, err := a.useCase.GetTopSources(ctx, parsedTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var result []*pb.SourceStats

	for _, stat := range stats {
		pbStat := &pb.SourceStats{
			Source:      stat.Source,
			TotalEvents: stat.TotalEvents,
		}
		result = append(result, pbStat)
	}
	return &pb.SourceResponse{Stats: result}, nil
}
 
func (a *AnalyticsController) GetHourlyTrends(ctx context.Context, req *pb.GetHourlyTrendsRequest) (*pb.GetHourlyTrendsResponse, error){
	parsedTime := req.Since.AsTime()
	stats, err := a.useCase.GetHourlyTrends(ctx, parsedTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var result []*pb.TrendStat
	for _, stat := range stats {
		pbStat := &pb.TrendStat{
			HourBucket: timestamppb.New(stat.HourBucket), 
			TotalEvents: stat.TotalEvents,
		}
		result = append(result, pbStat)
	}
	return &pb.GetHourlyTrendsResponse{Trends:result}, nil 
}
package http

import (
	pb "pulse/api/pb/analytics_v1"
	"strconv"
	"time"

	"pulse/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsHandler struct {
	logger *logger.Logger
	client pb.AnalyticsServiceClient
}

func NewAnalyticsHandler(logger *logger.Logger, client pb.AnalyticsServiceClient) *AnalyticsHandler {
	return &AnalyticsHandler{
		logger: logger,
		client: client,
	}
}

func (h *AnalyticsHandler) GetHourlyTrends(c *fiber.Ctx) error {
	sinceStr := c.Query("since")
	parsedInt, err := strconv.ParseInt(sinceStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid since parameter"})
	}
	req := &pb.GetHourlyTrendsRequest{Since: timestamppb.New(time.Unix(parsedInt, 0))}
	resp, err := h.client.GetHourlyTrends(c.UserContext(), req)
	if err != nil {
		h.logger.Error(err, "ошибка запроса GetHourlyTrends")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "grpc error"})
	}
	return c.JSON(resp)
}

func (h *AnalyticsHandler) GetTopSources(c *fiber.Ctx) error {
	sinceStr := c.Query("since")
	parsedInt, err := strconv.ParseInt(sinceStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid since parameter"})
	}
	req := &pb.GetTopSourcesRequest{Since: timestamppb.New(time.Unix(parsedInt, 0))}
	resp, err := h.client.GetTopSource(c.UserContext(), req)
	if err != nil {
		h.logger.Error(err, "ошибка grpc запроса GetTopSources")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "grpc error"})
	}
	return c.JSON(resp)
}

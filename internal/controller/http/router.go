package http

import "github.com/gofiber/fiber/v2"


func SetupRoutes(app *fiber.App, analyticsHandler *AnalyticsHandler) {
	v1 := app.Group("/api/v1")
	// v1.Use(jwtMiddleware)

	analyticsGroup := v1.Group("/analytics")
	analyticsGroup.Get("/trends", analyticsHandler.GetHourlyTrends)
	analyticsGroup.Get("/topsources", analyticsHandler.GetTopSources)
}
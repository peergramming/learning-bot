package routes

import (
	macaron "gopkg.in/macaron.v1"
)

func APIGetReportStatusHandler(ctx *macaron.Context) {
	// Allowed statuses: failed canceled running pending success success-with-warnings skipped not_found
	ctx.JSON(200, &map[string]string{
		"status": "skipped",
	})
}

package routes

import (
	macaron "gopkg.in/macaron.v1"
)

func APIGetProjectHandler(ctx *macaron.Context) {
	ctx.JSON(200, &map[string]string{
		"error": "to be implemented",
	})
}

func APIGenReportHandler(ctx *macaron.Context) {
	ctx.JSON(200, &map[string]string{
		"error": "to be implemented",
	})
}

func APIGetReportsHandler(ctx *macaron.Context) {
	ctx.JSON(200, &map[string]string{
		"error": "to be implemented",
	})
}

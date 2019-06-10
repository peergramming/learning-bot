package routes

import (
	macaron "gopkg.in/macaron.v1"
)

func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

func ReportPageHandler(ctx *macaron.Context) {
	ctx.PlainText(200, []byte("report page to be implemented"))
}

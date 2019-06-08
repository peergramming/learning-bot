package routes

import (
	macaron "gopkg.in/macaron.v1"
)

func HomepageHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

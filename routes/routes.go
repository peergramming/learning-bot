package routes

import (
	macaron "gopkg.in/macaron.v1"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
)

// ctxInit adds global context data to the page.
func ctxInit(ctx *macaron.Context) {
	ctx.Data["GitLabInst"] = settings.Config.GitLabInstanceURL
}

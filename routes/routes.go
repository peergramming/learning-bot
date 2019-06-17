package routes

import (
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	macaron "gopkg.in/macaron.v1"
)

// ctxInit adds global context data to the page.
func ctxInit(ctx *macaron.Context) {
	ctx.Data["GitLabInst"] = settings.Config.GitLabInstanceURL
}

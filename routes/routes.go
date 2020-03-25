package routes

import (
	"github.com/peergramming/learning-bot/modules/settings"
	macaron "gopkg.in/macaron.v1"
)

// ctxInit adds global context data to the page.
func ctxInit(ctx *macaron.Context) {
	ctx.Data["SiteTitle"] = settings.Config.SiteTitle
	ctx.Data["GitLabInst"] = settings.Config.GitLabInstanceURL
	ctx.Data["LMSTitle"] = settings.Config.LMSTitle
	ctx.Data["LMSUrl"] = settings.Config.LMSURL
}

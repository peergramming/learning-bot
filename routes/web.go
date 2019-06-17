package routes

import (
	"fmt"
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles the index page, which contains a description of the program.
func HomepageHandler(ctx *macaron.Context) {
	ctxInit(ctx)
	ctx.HTML(200, "index")
}

// ReportPageHandler handles rendering a report page.
func ReportPageHandler(ctx *macaron.Context) {
	ctxInit(ctx)
	ctx.Data["Project"] = fmt.Sprintf("%s/%s", ctx.Params("namespace"), ctx.Params("project"))
	commit := ctx.Params("sha")
	ctx.Data["Commit"] = commit
	ctx.Data["CommitShort"] = commit[:8]
	ctx.HTML(200, "report")
}

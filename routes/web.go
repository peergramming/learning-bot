package routes

import (
	"fmt"
	checkstyle "gitlab.com/gitedulab/learning-bot/models/checkstyle"
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler handles the index page, which contains a description of the program.
func HomepageHandler(ctx *macaron.Context) {
	ctxInit(ctx)
	ctx.HTML(200, "index")
}

// HelpCheckHandler handles the page where the check description would be displayed.
func HelpCheckHandler(ctx *macaron.Context) {
	ctxInit(ctx)

	check := ctx.Params("check")
	desc, ok := checkstyle.Checks[check]
	if !ok {
		ctx.Error(400, "Page does not exist")
		return
	}

	ctx.Data["Check"] = check
	ctx.Data["Description"] = desc.Description
	ctx.Data["Rationale"] = desc.Rationale
	ctx.Data["Suggestion"] = desc.Suggestion
	ctx.Data["Example"] = desc.Example

	ctx.HTML(200, "check_help")
}

// ReportPageHandler handles rendering a report page.
func ReportPageHandler(ctx *macaron.Context) {
	ctxInit(ctx)
	ctx.Data["Project"] = fmt.Sprintf("%s/%s", ctx.Params("namespace"),
		ctx.Params("project"))

	commit := ctx.Params("sha")
	ctx.Data["Commit"] = commit
	ctx.Data["CommitShort"] = commit[:8]

	ctx.HTML(200, "report")
}

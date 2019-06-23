package routes

import (
	"fmt"
	"log"
	"gitlab.com/gitedulab/learning-bot/models/checkstyle"
	"gitlab.com/gitedulab/learning-bot/modules/utils"
	"gitlab.com/gitedulab/learning-bot/models"
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
		ctx.Error(404, "Page does not exist")
		return
	}

	ctx.Data["Check"] = utils.Spacify(check)
	ctx.Data["Description"] = desc.Description
	ctx.Data["Rationale"] = desc.Rationale
	ctx.Data["Suggestion"] = desc.Suggestion
	ctx.Data["Example"] = desc.Example

	ctx.HTML(200, "check_help")
}

// ReportPageHandler handles rendering a report page.
func ReportPageHandler(ctx *macaron.Context) {
	ctxInit(ctx)

	project := fmt.Sprintf("%s/%s", ctx.Params("namespace"),
		ctx.Params("project"))
	commit := ctx.Params("sha")

	report, err := models.GetReport(project, commit)

	if err.Error() == "Report does not exist" {
		// Report does not exist, generate one
	} else if err != nil {
		// Some unknown error
		ctx.Error(500, "Server error")
		log.Printf("Failed to get report %s: %s", project, err)
		return
	}

	ctx.Data["Project"] = project
	ctx.Data["Commit"] = commit
	ctx.Data["CommitShort"] = commit[:8]

	fmt.Printf("Report has %s issues", len(report.Issues))

	ctx.HTML(200, "report")
}

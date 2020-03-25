package routes

import (
	"fmt"
	"github.com/peergramming/learning-bot/models"
	"github.com/peergramming/learning-bot/modules/checkstyle"
	"github.com/peergramming/learning-bot/modules/settings"
	"github.com/peergramming/learning-bot/modules/utils"
	macaron "gopkg.in/macaron.v1"
	"log"
	"strings"
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

func RequestDataHandler(ctx *macaron.Context) {
	key := ctx.Params("key")

	project := fmt.Sprintf("%s/%s", ctx.Params("namespace"),
		ctx.Params("project"))
	repo, err := models.GetRepo(project)

	if err != nil && err.Error() == "Repository does not exist" {
		ctx.Error(404, err.Error())
		return
	} else if err != nil {
		// Some unknown error
		//ctx.Error(500, "Server error")
		ctx.JSON(500, &map[string]string{"error": "Server error"})
		log.Printf("Failed to get report %s: %s", project, err)
		return
	}

	if repo.SecretKey != key {
		ctx.JSON(403, &map[string]string{"error": "No permissions to view this page"})
		//ctx.Error(403, "No permissions to view this page")
		return
	}

	// Load all issues
	for _, report := range repo.Reports {
		report.LoadIssues()
	}

	ctx.JSON(200, &repo)
}

// ReportsListPageHandler handles rendering list of reports.
func ReportsListPageHandler(ctx *macaron.Context) {
	ctxInit(ctx)

	project := fmt.Sprintf("%s/%s", ctx.Params("namespace"),
		ctx.Params("project"))
	key := ctx.Params("key")

	repo, err := models.GetRepo(project)

	if err != nil && err.Error() == "Repository does not exist" {
		ctx.Error(404, err.Error())
		return
	} else if err != nil {
		// Some unknown error
		ctx.Error(500, "Server error")
		log.Printf("Failed to get report %s: %s", project, err)
		return
	}

	if repo.SecretKey != key {
		ctx.Error(403, "No permissions to view this page")
		return
	}

	ctx.Data["Project"] = project
	ctx.Data["Reports"] = repo.Reports
	ctx.Data["Key"] = key

	ctx.HTML(200, "reports")
}

// ReportPageHandler handles rendering a report page.
func ReportPageHandler(ctx *macaron.Context) {
	ctxInit(ctx)

	project := fmt.Sprintf("%s/%s", ctx.Params("namespace"),
		ctx.Params("project"))
	commit := ctx.Params("sha")

	repo, err := models.GetRepo(project)

	if err != nil && err.Error() == "Repository does not exist" {
		ctx.Error(404, err.Error())
		return
	} else if err != nil {
		// Some unknown error
		ctx.Error(500, "Server error")
		log.Printf("Failed to get report %s: %s", project, err)
		return
	}

	rep, ok := repo.GetReport(commit)
	if !ok {
		ctx.Error(404, "Report does not exist")
		return
	}
	rep.LoadIssues()

	// Project data
	ctx.Data["Project"] = project
	ctx.Data["Commit"] = commit
	ctx.Data["CommitShort"] = commit[:8]
	ctx.Data["ReportGenDate"] = rep.CreatedUnix
	ctx.Data["Report"] = rep
	ctx.Data["SecretKey"] = repo.SecretKey

	if repo.LatestCommit != commit {
		ctx.Data["IsOldReport"] = true
	}

	// Survey data
	surveyConf := &settings.Config.Survey
	ctx.Data["ShowSurvey"] = surveyConf.ShowSurvey
	ctx.Data["SurveyTitle"] = surveyConf.Title
	ctx.Data["SurveyMessage"] = surveyConf.Message
	ctx.Data["SurveyURL"] = strings.Replace(surveyConf.SurveyURL, "$username",
		ctx.Params("namespace"), -1)

	ctx.HTML(200, "report")
}

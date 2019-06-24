package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"strings"
)

var CmdManage = cli.Command{
	Name:   "manage",
	Usage:  "Manage active projects",
	Action: runManage,
}

func runManage(clx *cli.Context) error {
	settings.LoadActiveProjs(true)

	if len(settings.ActiveProjs.Projects) == 0 {
		fmt.Println("There are no active projects configured.")
	} else {
		fmt.Printf("Current active projects (%d):\n", len(settings.ActiveProjs.Projects))
		for id, proj := range settings.ActiveProjs.Projects {
			fmt.Printf("\t%d: %s/%s\n", id, proj.Namespace, proj.Project)
		}
	}
	fmt.Println()

	fmt.Println("Enter a project path to add/remove (e.g. ha82/dsa-cw-1).")
	fmt.Println("Type 'q' or 'quit' to exit. Or ^C to abort changes.")

	fmt.Println()

	for {
		fmt.Printf("> ")
		var path string
		fmt.Scanln(&path)
		if path == "q" || path == "quit" {
			break
		}

		if path == "" {
			fmt.Println("Invalid project path")
			continue
		}
		paths := strings.Split(path, "/")
		if len(paths) != 2 {
			fmt.Printf("Invalid path, got %d paths, wanted 2!\n", len(paths))
			continue
		}

		if exists, id := settings.IsActiveProject(paths[0], paths[1]); exists {
			// Remove active project
			fmt.Printf("Are you sure you want to remove %s? (y/n) [n] ", path)
			var resp string
			fmt.Scanln(&resp)
			if resp == "y" {
				settings.ActiveProjs.Projects = remove(settings.ActiveProjs.Projects, id)
				fmt.Println("Removed!")
			}
		} else {
			settings.ActiveProjs.Projects = append(settings.ActiveProjs.Projects,
				settings.Project{Namespace: paths[0], Project: paths[1]})
			fmt.Printf("Added %s!\n", path)
		}
	}

	settings.SaveActiveProjs()
	fmt.Println("Saved!")

	return nil
}

// remove Function from The Go Programming Language (Donovan, Kernighan) book, page 93
func remove(slice []settings.Project, i int) []settings.Project {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

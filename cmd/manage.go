package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"os"
	"strings"
)

// CmdManage represents a command-line command
// which manages the list of active projects.
var CmdManage = cli.Command{
	Name:   "manage",
	Usage:  "Manage active projects",
	Action: runManage,
	Subcommands: []cli.Command{
		{
			Name:    "add",
			Usage:   "Add a project to the active projects list.",
			Aliases: []string{"a"},
			Action: func(c *cli.Context) error {
				settings.LoadActiveProjs(true)
				proj, err := getProjectFromString(c.Args().First())
				if err != nil {
					fmt.Printf("Failed to parse project: %s\n", err)
					os.Exit(2)
				}
				if exists, _ := settings.IsActiveProject(proj.Namespace, proj.Project); !exists {
					settings.ActiveProjs.Projects = append(settings.ActiveProjs.Projects, *proj)
					settings.SaveActiveProjs()
					fmt.Printf("Project added to list!\n")
				} else {
					fmt.Printf("Project is already included in list!\n")
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"rm", "del", "delete"},
			Usage:   "Remove a project from the active projects list.",
			Action: func(c *cli.Context) error {
				settings.LoadActiveProjs(true)
				proj, err := getProjectFromString(c.Args().First())
				if err != nil {
					fmt.Printf("Failed to parse project: %s\n", err)
					os.Exit(2)
				}
				if exists, id := settings.IsActiveProject(proj.Namespace, proj.Project); exists {
					settings.ActiveProjs.Projects = remove(settings.ActiveProjs.Projects, id)
					settings.SaveActiveProjs()
					fmt.Printf("Project removed from list!\n")
				} else {
					fmt.Printf("Project is not included in the list!\n")
					os.Exit(1)
				}
				return nil
			},
		},
	},
}

func getProjectFromString(path string) (project *settings.Project, err error) {
	if path == "" {
		return nil, fmt.Errorf("invalid project path")
	}
	paths := strings.Split(path, "/")
	if len(paths) != 2 {
		return nil, fmt.Errorf("invalid path, got %d paths, wanted 2", len(paths))
	}
	return &settings.Project{Namespace: paths[0], Project: paths[1]}, nil
}

func runManage(clx *cli.Context) error {
	settings.LoadActiveProjs(true)

	if len(settings.ActiveProjs.Projects) == 0 {
		fmt.Println("There are no active projects configured.")
	} else {
		fmt.Printf("Current active projects (%d):\n", len(settings.ActiveProjs.Projects))
		for id, proj := range settings.ActiveProjs.Projects {
			fmt.Printf("\t%d: %s\n", id, proj.GetFullPath())
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
		proj, err := getProjectFromString(path)
		if err != nil {
			fmt.Printf("Invalid project path: %s\n", err)
			continue
		}

		if exists, id := settings.IsActiveProject(proj.Namespace, proj.Project); exists {
			// Remove active project
			fmt.Printf("Are you sure you want to remove %s? (y/n) [n] ", path)
			var resp string
			fmt.Scanln(&resp)
			if resp == "y" {
				settings.ActiveProjs.Projects = remove(settings.ActiveProjs.Projects, id)
				fmt.Println("Removed!")
			}
		} else {
			settings.ActiveProjs.Projects = append(settings.ActiveProjs.Projects, *proj)
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

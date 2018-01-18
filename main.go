package main

import (
	"os"
	"strings"

	"github.com/luhring/commitment/commitment"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "commitment"
	app.Usage = "See fun facts about repositories and their commit history!"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "first",
			Aliases: []string{"1st", "1"},
			Usage:   "See the first commit of a repository",
			Action: func(c *cli.Context) {
				repositoryPath := c.Args().Get(0)
				repository := getRepositoryFromPath(repositoryPath)
				commitReport := repository.GetCommitReport(0)
				commitReport.Print()
			},
		},
		{
			Name:    "second",
			Aliases: []string{"2nd", "2"},
			Usage:   "See the second commit of a repository",
			Action: func(c *cli.Context) {
				repositoryPath := c.Args().Get(0)
				repository := getRepositoryFromPath(repositoryPath)
				commitReport := repository.GetCommitReport(1)
				commitReport.Print()
			},
		},
	}

	app.Run(os.Args)
}

func getRepositoryFromPath(path string) commitment.Repository {
	parts := strings.Split(path, "/")
	return commitment.Repository{
		User:           parts[0],
		RepositoryName: parts[1],
	}
}

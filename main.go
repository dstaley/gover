package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gover"
	app.Usage = "Monitor your CGMiner instances in style!"
	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "create configuration files and the database",
			Action: func(c *cli.Context) {
				setup()
			},
		},
		{
			Name:      "server",
			ShortName: "s",
			Usage:     "start the gover server",
			Action: func(c *cli.Context) {
				server()
			},
		},
	}
	app.Run(os.Args)
}

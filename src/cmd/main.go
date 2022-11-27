package main

import (
	"os"

	proxy "github.com/mmohamed/loki-auth-proxy/src/proxy"
	"github.com/urfave/cli"
)

var (
	version = "dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "Loki Authentication Proxy"
	app.Version = version
	app.Authors = []cli.Author{
		{Name: "Marouan MOHAMED", Email: "medmarouen@gmail.com"},
	}
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Runs the Loki multi tenant proxy",
			Action: proxy.Serve,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Usage: "Port to expose this loki proxy",
					Value: 3501,
				}, cli.StringFlag{
					Name:  "loki-server",
					Usage: "Loki server endpoint",
					Value: "http://localhost:3500",
				}, cli.StringFlag{
					Name:  "auth-config",
					Usage: "Auth yaml configuration file path",
					Value: "auth.yaml",
				}, cli.BoolFlag{
					Name:     "org-check",
					Usage:    "Require XOrgId header and match user account",
					Required: false,
				},
			},
		},
	}
	app.Run(os.Args)
}

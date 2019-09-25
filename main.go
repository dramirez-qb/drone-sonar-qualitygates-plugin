package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var build = "1" // build number set at compile time

func main() {
	app := cli.NewApp()
	app.Name = "Drone-Sonar-QualityGates-Plugin"
	app.Usage = "Drone plugin to integrate with SonarQube Quality Gates."
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "token",
			Usage:  "SonarQube token",
			EnvVar: "PLUGIN_SONAR_TOKEN",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	plugin := Plugin{
		Config: Config{
			Token: basicAuth(c.String("token")),
		},
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func basicAuth(token string) string {
	auth := token + ":"
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

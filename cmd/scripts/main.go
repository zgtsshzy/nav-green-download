package main

import (
	"fmt"
	"nav-green-download/cmd/scripts/gfs"
	"nav-green-download/cmd/scripts/mfwam"
	"nav-green-download/cmd/scripts/seaice"
	"nav-green-download/cmd/scripts/smoc"
	"nav-green-download/pkg/conf"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Info struct {
	Type      string
	StartTime string
	EndTime   string
}

func main() {
	info := new(Info)

	app := cli.App{
		Name:   "nav-green-download",
		Usage:  "some scripts",
		Flags:  getFlags(info),
		Action: getAction(info),
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("app run error: %v", err)
	}
}

func getFlags(info *Info) []cli.Flag {
	flags := make([]cli.Flag, 0)

	flags = append(flags, &cli.StringFlag{
		Name:    "type",
		Aliases: []string{"t"},
		EnvVars: []string{
			"smoc", "mfwam", "sea_ice",
		},
		Value:       "",
		Required:    true,
		Hidden:      false,
		Usage:       "what kind task to do",
		Destination: &info.Type,
	})

	flags = append(flags, &cli.StringFlag{
		Name:    "start time",
		Aliases: []string{"s"},
		EnvVars: []string{
			"2006-01-02 15:04:05",
		},
		Value:       "",
		Required:    false,
		Hidden:      false,
		Usage:       "start insert time",
		Destination: &info.StartTime,
	})

	flags = append(flags, &cli.StringFlag{
		Name:    "end time",
		Aliases: []string{"e"},
		EnvVars: []string{
			"2006-01-02 15:04:05",
		},
		Value:       "",
		Required:    false,
		Hidden:      false,
		Usage:       "end insert time",
		Destination: &info.EndTime,
	})

	return flags
}

func getAction(info *Info) func(*cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		c := conf.New()
		c.Show()

		switch info.Type {
		case "smoc":
			return smoc.ExecuteScript(info.StartTime, info.EndTime)
		case "gfs":
			return gfs.ExecuteScript(info.StartTime, info.EndTime)
		case "mfwam":
			return mfwam.ExecuteScript(info.StartTime, info.EndTime)
		case "sea_ice":
			return seaice.ExecuteScript(info.StartTime, info.EndTime)
		default:
			return fmt.Errorf("type: %s is invalid", info.Type)
		}
	}
}

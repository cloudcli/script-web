package main

import (
	"fmt"
	"os"
	"script-web/cmd"
	"script-web/modules/config"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/rifflock/lfshook"
	"path"
)

var (
	Name = "Cloudcli"
	Version = "0.0.1"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	//log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	fmt.Println("hello script management platform")

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "lang",
			Value: "Zh-cn",
			Usage: "Language",
		},
		cli.BoolFlag{
			Name: "debug",
			Usage: "debug ",
		},
		cli.StringFlag{
			Name: "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command {
		cmd.WebCommand,
		cmd.BootstrapCommand,
	}

	app.Before = func(c *cli.Context) error {
		config.InitConfig(c.String("config"))
		log.AddHook(lfshook.NewHook(lfshook.PathMap{
			log.InfoLevel: path.Join(config.LogDir, "info.log"),
			log.FatalLevel: path.Join(config.LogDir, "fatal.log"),
			log.ErrorLevel: path.Join(config.LogDir, "error.log"),
		}))

		return nil
	}

	app.Action = func (c *cli.Context) error {

		//cli.DefaultAppComplete(c)
		//cli.HandleExitCoder(errors.New("not an exit coder, though"))
		cli.ShowAppHelp(c)
		//cli.ShowCommandCompletions(c, "nope")
		//cli.ShowCommandHelp(c, "get")
		//cli.ShowCompletions(c)
		//cli.ShowSubcommandHelp(c)
		//cli.ShowVersion(c)
		return nil
	}

	app.Run(os.Args)
}

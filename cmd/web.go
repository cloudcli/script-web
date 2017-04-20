package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"script-web/modules/config"
	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
	"script-web/modules/controllers"
	"bytes"
	"html/template"
	"script-web/modules/middleware"
)

var WebCommand = cli.Command{
	Name: "web",
	Usage: "web",
	Description: "web sevice management",
	Subcommands: []cli.Command{
		cli.Command{
			Name: "start",
			Usage: "start web service",
			Description: "start web service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "port, p",
					Value: "3001",
					Usage: "web service http port",
				},
			},
			Action: webStartAction,
		},
	},
}

func webStartAction(ctx *cli.Context) error {
	log.WithFields(log.Fields{
		"config": config.ConfigPath,
		"port": ctx.String("port"),
	}).Info("start web http service")

	templateConfig := html.DefaultConfig()
	templateConfig.Layout = "layouts/layout.html"
	templateConfig.Funcs["raw"] = func(input string) string {
		return input
	}

	tpl := html.New(templateConfig)
	templateConfig.Funcs["CallTemplate"] = func(name string, data interface{}) (ret template.HTML, err error) {
		buf := bytes.NewBuffer([]byte{})
		err = tpl.Templates.ExecuteTemplate(buf, name, data)
		ret = template.HTML(buf.String())
		return
	}

	// @todo use go-bindata compiled static file
	iris.UseTemplate(tpl).Directory("./frontend/templates", ".html")
	iris.StaticWeb("/static", "./frontend/public")

	// home page


	//homes := iris.Party("/", middleware.InjectDatabase)
	iris.Get("/", controllers.Home)


	//scripts.Get("/*path", controllers.Scripts)

	iris.Get("/scripts/*path", middleware.RepoMiddleware, controllers.Scripts)
	iris.Get("/scripts", middleware.RepoMiddleware, controllers.Scripts)
	iris.Post("/script/apply/new", controllers.PostNewScript)

	// all scripts
	//script := iris.Party("/s", middleware.AuthMiddleware)
	//{
	//	// new script
	//	script.Get("/new", controllers.NewScript)
	//
	//	// script detail - tab code
	//	script.Get("/:id", controllers.ScriptDetail)
	//
	//	// script detail - tab revisions
	//	script.Get("/:id/revisions", controllers.ScriptDetailRevisions)
	//
	//	// script detail - tab discuss
	//	script.Get("/:id/discuss", controllers.ScriptDetailDiscuss)
	//
	//	// script detail - tab setting
	//	script.Get("/:id/setting", controllers.ScriptDetailSetting)
	//
	//	// script detail - tab update proposals
	//	script.Get("/:id/proposals", controllers.ScriptDetailProposals)
	//}

	//iris.Get("/scripts/", controllers.Scripts)

	// new script
	//iris.Get("/new/script", controllers.NewScript)
	//
	// script detail - tab code
	iris.Get("/s/:id", middleware.RepoMiddleware, controllers.ScriptDetail)

	// script detail - tab revisions
	iris.Get("/s/:id/revisions", controllers.ScriptDetailRevisions)

	// script detail - tab discuss
	iris.Get("/s/:id/discuss", controllers.ScriptDetailDiscuss)

	// script detail - tab setting
	iris.Get("/s/:id/setting", controllers.ScriptDetailSetting)

	// script detail - tab update proposals
	iris.Get("/s/:id/proposals", controllers.ScriptDetailProposals)

	users := iris.Party("/auth", middleware.InjectDatabase)
	{
		users.Get("/u/:name", controllers.Userspace)

		users.Get("/signin", controllers.SignIn)
		users.Post("/signin", controllers.PostSignIn)
		users.Get("/signup", controllers.SignUp)
		users.Post("/signup", controllers.PostSignUp)
	}
	// user space
	//iris.Get("/u/:name", controllers.Userspace)
	//iris.Get("/signin", controllers.SignIn)
	//iris.Post("/signin", controllers.PostSignIn)
	//iris.Get("/signup", controllers.SignUp)
	//iris.Post("/signup", controllers.PostSignUp)


	iris.Listen(":" + ctx.String("port"))

	return nil
}
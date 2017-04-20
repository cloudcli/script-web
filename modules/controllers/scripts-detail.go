package controllers

import (
	"github.com/kataras/iris"
	_ "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus"
)

func ScriptDetail(ctx *iris.Context) {
	id := ctx.Get("id")
	logrus.Debug("current id:", id)



	ctx.MustRender("scripts-detail.html", &PageDT{
		Data: map[string]interface{} {
			"tab": "code",
		},
	})
}

func ScriptDetailRevisions(ctx *iris.Context) {
	ctx.MustRender("scripts-detail.html", &PageDT{
		Data: map[string]interface{} {
			"tab": "revisions",
		},
	})
}

func ScriptDetailDiscuss(ctx *iris.Context) {
	ctx.MustRender("scripts-detail.html", &PageDT{
		Data: map[string]interface{} {
			"tab": "discuss",
		},
	})
}

func ScriptDetailProposals(ctx *iris.Context) {
	ctx.MustRender("scripts-detail.html", &PageDT{
		Data: map[string]interface{} {
			"tab": "proposals",
		},
	})
}

func ScriptDetailSetting(ctx *iris.Context) {
	ctx.MustRender("scripts-detail.html", &PageDT{
		Data: map[string]interface{} {
			"tab": "setting",
		},
	})
}
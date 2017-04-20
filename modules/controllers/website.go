package controllers

import "github.com/kataras/iris"

type PageDT struct{
	Error 		string
	Data 		map[string]interface{}
}

type JsonDT struct {
	Status  	string
	Content 	interface{}
	Meta 		interface{}
}

type Err struct{
	Error error
}

func Err404(ctx *iris.Context) {
	// get flash error
	ctx.MustRender("404.html", nil)
}

func JsonErr500(ctx *iris.Context, err error, code int) {
	ctx.JSON(500, JsonDT{
		Status: "failed",
		Meta: map[string]interface{} {
			"message": err.Error(),
			"code": code,
		},
	})
}

func Err500(ctx *iris.Context, err error) {

	// get flash error
	ctx.MustRender("500.html", &PageDT{
		Error: err.Error(),
	})
}

func Home(ctx *iris.Context) {
	ctx.MustRender("index.html", &PageDT{
		Data: map[string]interface{} {
		},
	})
}


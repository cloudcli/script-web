package controllers

import (
	"script-web/modules/db"
	"script-web/modules/db/models"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	//"script-web/modules/util"
)

type DT struct {
	Name  string
	Error map[string]string
	Data  map[string]interface{}
}

func Userspace(ctx *iris.Context) {
	db := ctx.Session().Get("_db").(*db.MysqlDB)
	// TODO: must check login before
	//email := ctx.Session().Get("UserEmail").(string)
	email := "admin@cloudcli.com"
	user, err := db.FindUserByKey("email", email)
	if err != nil {
		ctx.MustRender("userspace.html", DT{
			Error: map[string]string{
				"message": "服务端错误",
				"code":    "1",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	ctx.MustRender("userspace.html", &PageDT{
		Data: map[string]interface{}{
			"User": user,
		},
	})
}


func UserMessage(ctx *iris.Context) {

}

func SignIn(ctx *iris.Context) {
	type DT struct {
		Name   string
		Error  map[string]string
		Common map[string]interface{}
	}

	UserEmail := ctx.Session().Get("UserEmail")

	if UserEmail != nil {
		ctx.Redirect("/")
		return
	}

	ctx.MustRender("signin.html", PageDT{
		Data: map[string]interface{}{},
	})
}

func PostSignIn(ctx *iris.Context) {
	type DT struct {
		Name  string
		Error map[string]string
		Data  map[string]interface{}
	}

	var (
		user models.User
		err error
	)

	//name := ctx.FormValue("name")
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")

	db := ctx.Session().Get("_db").(*db.MysqlDB)

	//if err != nil {
	//	Err500(ctx, err)
	//	return
	//}

	//if name == "" && email == "" {
	//	ctx.MustRender("signin.html", DT{
	//		Error: map[string]string{
	//			"message": "用户名不能为空",
	//			"code": "0",
	//		},
	//		Data: map[string]interface{}{},
	//	})
	//	return
	//}
	//
	//if name != "" {
	//	user, err = db.FindUserByKey("name", name)
	//} else {
	//	user, err = db.FindUserByKey("email", email)
	//}

	if email == "" {
		ctx.MustRender("signin.html", DT{
			Error: map[string]string{
				"message": "邮箱地址不能为空",
				"code":    "0",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	user, err = db.FindUserByKey("email", email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.MustRender("signin.html", DT{
				Error: map[string]string{
					"message": "找不到用户",
					"code":    "1",
				},
				Data: map[string]interface{}{},
			})
		} else {
			ctx.MustRender("signin.html", DT{
				Error: map[string]string{
					"message": err.Error(),
					"code":    "1",
				},
				Data: map[string]interface{}{},
			})
		}
		return
	}

	if !user.IsValidPassword(password) {
		ctx.MustRender("signin.html", DT{
			Error: map[string]string{
				"message": "用户名或密码错误",
				"code":    "1",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	//enpassword, err := util.EncodePassword(password)

	//if err != nil {
	//	ctx.MustRender("signin.html", DT{
	//		Error: map[string]string{
	//			"message": "服务端错误",
	//			"code": "1",
	//		},
	//		Data: map[string]interface{}{},
	//	})
	//	return
	//}

	// create jwt auth token
	// authorizer := &simple.Authorizer {
	//	SignAlgorithm: settings.SignAlgorithm,
	//	AccessKey: settings.AccessKey,
	//	AccessTimeout: settings.AccessTimeout,
	// }

	ctx.Session().Set("UserEmail", user.Email)
	ctx.Session().Set("UserId", user.ID)

	ctx.Redirect("/auth/u/" + user.Name)
}

func SignUp(ctx *iris.Context) {
	type DT struct {
		Name   string
		Error  map[string]string
		Common map[string]interface{}
	}
	ctx.MustRender("signup.html", PageDT{
		Data: map[string]interface{}{},
	})
}

func PostSignUp(ctx *iris.Context) {
	type DT struct {
		Name  string
		Error map[string]string
		Data  map[string]interface{}
	}

	name := ctx.FormValue("name")
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")

	db := ctx.Session().Get("_db").(*db.MysqlDB)

	if name == "" || email == "" || password == "" {
		// TODO log it
		ctx.MustRender("signup.html", DT{
			Error: map[string]string{
				"message": "用户名或邮箱或密码不能为空",
				"code":    "1",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	isUserExist, err := db.IsUserEmailExist(email)

	if err != nil {
		// TODO log it
		ctx.MustRender("signup.html", DT{
			Error: map[string]string{
				"message": "系统错误",
				"code":    "1",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	if isUserExist {
		ctx.MustRender("signup.html", DT{
			Error: map[string]string{
				"message": "用户已存在",
				"code":    "1",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	user, err := db.CreateUser(email, name, password)

	if err != nil {
		// TODO log it
		ctx.MustRender("signup.html", DT{
			Error: map[string]string{
				"message": "系统错误",
				"code":    "1",
			},
			Data: map[string]interface{}{},
		})
		return
	}

	ctx.Session().Set("UserEmail", user.Email)
	ctx.Session().Set("UserId", user.ID)

	ctx.Redirect("/auth/u/" + user.Name)
}

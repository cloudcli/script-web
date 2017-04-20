package middleware

import (
	"github.com/kataras/iris"
	"script-web/modules/db"
	"fmt"
)

// InjectDatabase put db to context
func InjectDatabase(ctx *iris.Context) {
	db, err := db.GetDb()
	if err != nil {
		// log something and panic
		fmt.Println(err)
		//panic(err)
	}

	ctx.Session().Set("_db", db)
	ctx.Next()
}
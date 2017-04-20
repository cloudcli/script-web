package middleware

import (
	"github.com/kataras/iris"
	"github.com/Sirupsen/logrus"
)


func AuthMiddleware(ctx *iris.Context) {
	logrus.Debug("start auth")
}
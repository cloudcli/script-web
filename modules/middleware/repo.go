package middleware

import (
	"github.com/kataras/iris"
	"github.com/Sirupsen/logrus"
	"github.com/gogits/git-module"
	"script-web/modules/config"
)

// RepoMiddleware
func RepoMiddleware(ctx *iris.Context) {
	repo, err := git.OpenRepository(config.RepoPath)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect("/500")
		return
	}
	ctx.Set("repo", repo)

	commit, err := repo.GetBranchCommit("master")
	if err != nil {
		logrus.Error(err)
		ctx.Redirect("/500")
		return
	}
	ctx.Set("commit", commit)

	ctx.Next()
}

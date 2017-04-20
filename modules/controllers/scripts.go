package controllers

import (
	"github.com/kataras/iris"
	//"script-web/modules/git"
	//"bytes"
	//"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	//"path"
	"script-web/modules/widget"
	"github.com/Sirupsen/logrus"
	"github.com/gogits/git-module"
	"strings"
	"time"
)

type PageMeta struct {
	Name     string
	ObjectId string
	Updated  time.Time
	Type     string
}

func Scripts(ctx *iris.Context) {
	dirpath := ctx.Param("path")
	if dirpath == "" {
		dirpath = "/"
	}

	repo := ctx.Get("repo").(*git.Repository)
	commit := ctx.Get("commit").(*git.Commit)
	entries, err := commit.ListEntries()
	if err != nil {
		logrus.Error(err)
		return
	}

	pages := make([]PageMeta, 0, len(entries))

	for i := range entries {

		c, err := repo.GetCommitByPath(entries[i].Name())
		if err != nil {
			logrus.Error(err)
			return
		}

		name := entries[i].Name()
		pages = append(pages, PageMeta{
			Name:     name,
			ObjectId: entries[i].ID.String(),
			Type:     string(entries[i].Type),
			Updated:  c.Author.When,
		})
	}

	logrus.Debugf("current pages %+v", pages)


	//
	//logrus.Debugf("current commit author: %+v", commit.Author.Name)
	//logrus.Debugf("current commit author: %+v", commit.Author.Email)
	//logrus.Debugf("current commit author: %+v", commit.Author.When)
	//logrus.Debugf("current commit author: %+v", commit.ID)
	////logrus.Debugf("current commit author: %+v", commit.ListEntries())
	//
	//fileStatus, err := commit.FileStatus()
	//if err != nil {
	//	logrus.Error(err)
	//	return
	//}
	//logrus.Debugf("current commit fileStatus: %+v", fileStatus)
	//
	//logrus.Debugf("curent entries %+v", entries)
	//
	//for _, entry := range entries {
	//	logrus.Debugf("blob name %s", entry.Blob().Name())
	//}
	//

	//
	//git.NewTree(repo, )
	//
	//
	//gitRepo, err := git.GetRepo(config.RepoPath)
	//if err != nil {
	//	Err500(ctx, err)
	//	return
	//}
	//
	//out, err := gitRepo.ListTree(path.Join(".", dirpath), "master")
	//if err != nil {
	//	logrus.Errorf("scripts: get git tree error: %v", err.Error())
	//	Err500(ctx, err)
	//	return
	//}
	//
	//buf := bytes.Buffer{}
	//_, err = buf.ReadFrom(out)
	//if err != nil {
	//	logrus.Errorf("scripts: read bytes error: %v", err)
	//	Err500(ctx, err)
	//	return
	//}
	//
	//rawString := buf.String()
	//listFile := git.ListFile{}
	//err = listFile.Init(rawString)
	//if err != nil {
	//	logrus.Errorf("scripts: parse %s error: %v", rawString, err)
	//	Err500(ctx, errors.Errorf("scripts: parse %s error: %v", rawString, err))
	//	return
	//}
	//
	segs := strings.Split(dirpath, "/")
	breadcrumb := widget.NewBreadcrumb()
	baseUrl := "/scripts/"
	breadcrumb.HomeUrl = baseUrl
	breadcrumb.PrevUrl = "/scripts" + strings.Join(segs[:len(segs) - 1], "/")
	for i, seg := range segs {
		breadcrumb.Add(seg, "/scripts" + strings.Join(segs[:i + 1], "/"))
	}

	breadcrumbHtml, err := breadcrumb.Render()

	if err != nil {
		err = errors.Errorf("scripts: render breadcrumb error: %v", err)
		logrus.Error(err.Error())
		Err500(ctx, err)
		return
	}

	logrus.Debugf("current list file is : %+v", commit)

	//logrus.Info(breadcrumbHtml)
	// find all blob files
	// query all scripts meta data from mysql
	ctx.MustRender("scripts.html", &PageDT{
		Data: map[string]interface{}{
			"Pages":      pages,
			"breadcrumb": breadcrumbHtml,
		},
	})
}

package controllers

import (
	"script-web/modules/config"
	"script-web/modules/db"
	"script-web/modules/db/models"
	"script-web/modules/git"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"strings"
)

func NewScript(ctx *iris.Context) {
	ctx.MustRender("scripts-new-script.html", &PageDT{
		Data: map[string]interface{}{},
	})
}

/**
Post Apply New Script

0. Pre Condition:
	- user signed in
1. if user is admin
	- Create Script Record
	- Create Script Version
2. else
	1. Create Apply Record
		- check dir
		- check name repeation
3. Create Message Record
4. Create Event Log Record

*/
func PostNewScript(ctx *iris.Context) {

	// @todo check is user logged in
	// @todo check is user is admin

	user := "admin"

	if user == "admin" {
		PostNewScriptAdmin(ctx)
		return
	}

	path := ctx.FormValue("path")
	assignTo := ctx.FormValue("assignTo")
	description := ctx.FormValue("description")
	tags := ctx.FormValue("tags")
	scriptName := ctx.FormValue("name")
	scriptContent := ctx.FormValue("content")

	mysqldb, err := db.GetDb()

	if err != nil {
		JsonErr500(ctx, err, 0)
		return
	}

	payload, err := json.Marshal(map[string]interface{}{
		"path":    path,
		"tags":    tags,
		"name":    scriptName,
		"content": scriptContent,
	})

	if err != nil {
		JsonErr500(ctx, err, 1)
		return
	}

	title := fmt.Sprintf("%s提交创建新脚本%s请求", user, scriptName)

	apply := &models.Apply{
		Proposer:    user,
		AssignTo:    assignTo,
		Title:       title,
		Description: description,
		Type:        models.ApplyTypeNewScript,
		Payload:     string(payload),
	}

	if err = mysqldb.Db.Create(apply).Error; err != nil {
		JsonErr500(ctx, err, 2)
		return
	}

	ctx.JSON(200, JsonDT{
		Status: "success",
		Content: map[string]interface{}{
			"apply": apply,
		},
	})
}

/**
  PostNewScriptAdmin

  condition: when user is admin , create script directly

  1. create Script record
  2. commit script to git repo
  3. create Script version record
  4. create message record
  5. create event log record

*/
func PostNewScriptAdmin(ctx *iris.Context) {
	user := "admin"
	email := "admin@idcos.com"
	path := ctx.FormValue("path")
	description := ctx.FormValue("description")
	tags := ctx.FormValue("tags")
	scriptName := ctx.FormValue("name")
	scriptContent := ctx.FormValue("content")

	mysqldb, err := db.GetDb()

	if err != nil {
		JsonErr500(ctx, err, 0)
		return
	}


	// 多个tag数组
	arrayTag := strings.Split(tags, ",")
	fmt.Println(arrayTag)


	// create script record
	script := &models.Script{
		Name:          scriptName,
		Path:          path,
		ScriptVersion: 1,
		Description:   description,
		//Tags:          tags,
	}

	if err = mysqldb.Db.Create(script).Error; err != nil {
		JsonErr500(ctx, err, 1)
		return
	}

	gitrepo, err := git.GetRepo(config.RepoPath)
	if err != nil {
		JsonErr500(ctx, err, 2)
		return
	}

	ref, diff, err := gitrepo.CommitFile(scriptName, path, []byte(scriptContent), user, email)
	if err != nil {
		JsonErr500(ctx, err, 3)
		return
	}

	version := models.ScriptVersion{
		ScriptId: script.ID,
		ApplyId:  0,
		Version:  1,
		Hash:     ref,
		Diff:     diff,
	}

	ctx.JSON(200, JsonDT{
		Status: "success",
		Content: map[string]interface{}{
			"version": version,
		},
	})
}

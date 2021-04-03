package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

func dashboard(ctx iris.Context){
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.Redirect("/", iris.StatusPermanentRedirect)
		return
	}
	// need to check redis key is valid
	// get all redis keys
	// check session_id is exist

	DashInfo := createJSON()
	ctx.JSON(DashInfo)
}


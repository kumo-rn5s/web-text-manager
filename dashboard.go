package main

import "github.com/kataras/iris/v12"

func dashboard(ctx iris.Context){
	// need to check redis key is valid
	// get all redis keys
	// check session_id is exist

	DashInfo := createJSON()
	ctx.JSON(DashInfo)
}


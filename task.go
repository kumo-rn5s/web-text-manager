package main

import (
	"github.com/kataras/iris/v12"
	"io/ioutil"
)


func getTask(ctx iris.Context)  {
	mdData, err := ioutil.ReadFile("userdata/task/task.md")
	Message := []byte("Default")
	if err != nil{
		Message = []byte("# No task")
	} else {
		Message = []byte(string(mdData))
	}
	Task := iris.Map{
		"task" : Message,
	}
	ctx.JSON(Task)
}


func SendTask(ctx iris.Context)  {
	var taskStream MDStream
	ResponseMessage := []byte("Default")
	err := ctx.ReadJSON(&taskStream)
	if err != nil{
		ResponseMessage = []byte("JSON Format Error")
	} else {
		err2 := ioutil.WriteFile("userdata/task/task.md", []byte(taskStream.DataStream), 755)
		if err2 != nil {
			ResponseMessage = []byte("error1")
		} else{
			ResponseMessage = []byte("fooooooooo!!!!!")
		}
	}
	Response := iris.Map{
		"response":ResponseMessage,
	}
	ctx.JSON(Response)
}
package main

import (
	"github.com/kataras/iris/v12"
	"io/ioutil"
)

func createJSON()iris.Map{
	// Function Arguments need to contain db access infos
	// Connect to Database

	// store data
	DBUsername := "test1"
	DBProfilepic := "/user/test1/pic/1.jpg"
	DBLastLogin := "2021-3-21 10:34"
	data, err := ioutil.ReadFile("userdata/task/task.md")
	DBTaskList := []byte("Default")
	if err != nil{
		DBTaskList = []byte("# No task")
	} else {
		DBTaskList = []byte(string(data))
	}
	// Return iris Map type
	return iris.Map{
		"username" : DBUsername,
		"profile_pic" : DBProfilepic,
		"last_login" : DBLastLogin,
		"task_list" : DBTaskList,
	}
}

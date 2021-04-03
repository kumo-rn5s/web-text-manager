package main

import (
	"github.com/kataras/iris/v12"
)

func main() {

	app := iris.New()
	app.HandleDir("/", iris.Dir("./dist"))

	sess := ConnectRedis()
	ConnectMySQL()

	app.Use(sess.Handler())

	app.Post("/login", login)
	app.Get("/logout",logout)

	app.Get("/dashboard",dashboard)
	app.Get("/task", getTask)
	app.Post("/task", SendTask)

	app.Get("/filepath",getFilePath)
	app.Get("/fileList",showAllFiles)

	app.Post("/downloadFile", downloadFile)
	app.Post("/deleteFile", deleteFile)
	app.Post("/saveFile", saveFile)
	app.Get("/getFile", getFile)

	//app.Get("/books", FindBooks)
	//app.Get("/books/:id", FindBook)
	//app.Post("/books", CreateBook)
	//app.Patch("/books/:id", UpdateBook)
	//app.Delete("/books/:id", DeleteBook)

	app.Listen(":8080")
}

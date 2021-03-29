package main

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

func main() {

	app := iris.New()
	app.HandleDir("/", iris.Dir("./dist"))

	db := redis.New(redis.Config{
		Network:   "tcp",
		Addr:      getEnv("REDIS_ADDR", "redis:6379"),
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Username:  "",
		Password:  "",
		Database:  "",
		Prefix:    "mds-",
		Driver:    redis.GoRedis(), // defaults.
	})
	defer db.Close() // close the database connection if application errored.

	sess := sessions.New(sessions.Config{
		Cookie:          "md_session",
		Expires:         -1, // defaults to -1: delete by closing browser
		AllowReclaim:    true,
		CookieSecureTLS: true,
	})

	sess.UseDatabase(db)
	app.Use(sess.Handler())

	// set session values.
	app.Get("/set", func(ctx iris.Context) {

	})

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

	app.Listen(":8080")
}

package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"io/ioutil"
	"os"
	"time"
)

var MYSQLDB *gorm.DB
var REDISDB *redis.Database

type dbConfiguration struct {
	DBHost    string	`json:"dbHost"`
	DBPort    string	`json:"dbPort"`
	DBName    string	`json:"dbName"`
	DBUser    string	`json:"dbUser"`
	DBPass    string	`json:"dbPass"`
}

type User struct {
	Username    string   `json:"user" gorm:"primary_key"`
	Password    string   `json:"pass"`
}

//type CreateBookInput struct {
//	Title  string `json:"title" binding:"required"`
//	Author string `json:"author" binding:"required"`
//}
//
//type UpdateBookInput struct {
//	Title  string `json:"title"`
//	Author string `json:"author"`
//}


func ConnectRedis() *sessions.Sessions {

	REDISDB = redis.New(redis.Config{
		Network:   "tcp",
		Addr:      getEnv("REDIS_ADDR", "localhost:6379"),
		//Addr:      getEnv("REDIS_ADDR", "redis:6379"),
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Username:  "",
		Password:  "",
		Database:  "",
		Prefix:    "mds-",
		Driver:    redis.GoRedis(), // defaults.
	})
	defer REDISDB.Close() // close the database connection if application errored.

	sess := sessions.New(sessions.Config{
		Cookie:          "md_session",
		Expires:         -1, // defaults to -1: delete by closing browser
		AllowReclaim:    true,
		CookieSecureTLS: true,
	})
	sess.UseDatabase(REDISDB)
	return sess
}

func ConnectMySQL() {

	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)

	var dbConf dbConfiguration
	_ = decoder.Decode(&dbConf)
	dst := dbConf.DBUser+":"+dbConf.DBPass+"@tcp("+dbConf.DBHost+":"+dbConf.DBPort+")/"+dbConf.DBName

	database, err := gorm.Open("mysql", dst+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{})
	// CREATE Table

	MYSQLDB = database
}


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



//
//// GET /books
//// Find all books
//func FindBooks(ctx iris.Context) {
//	var books []Book
//	MYSQLDB.Find(&books)
//
//	ctx.JSON(http.StatusOK, gin.H{"data": books})
//}
//
//// GET /books/:id
//// Find a book
//func FindBook(ctx iris.Context) {
//	// Get model if exist
//	var book Book
//	if err := MYSQLDB.Where("id = ?", ctx.Param("id")).First(&book).Error; err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, gin.H{"data": book})
//}
//
//// POST /books
//// Create new book
//func CreateBook(ctx iris.Context) {
//	// Validate input
//	var input CreateBookInput
//	if err := ctx.ShouldBindJSON(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	// Create book
//	book :=  Book{Title: input.Title, Author: input.Author}
//	 MYSQLDB.Create(&book)
//
//	ctx.JSON(http.StatusOK, gin.H{"data": book})
//}
//
//// PATCH /books/:id
//// Update a book
//func UpdateBook(ctx iris.Context) {
//	// Get model if exist
//	var book  Book
//	if err :=  MYSQLDB.Where("id = ?", ctx.Param("id")).First(&book).Error; err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
//		return
//	}
//
//	// Validate input
//	var input UpdateBookInput
//	if err := ctx.ShouldBindJSON(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	 MYSQLDB.Model(&book).Updates(input)
//
//	ctx.JSON(http.StatusOK, gin.H{"data": book})
//}
//
//// DELETE /books/:id
//// Delete a book
//func DeleteBook(ctx iris.Context) {
//	// Get model if exist
//	var book  Book
//	if err :=  MYSQLDB.Where("id = ?", ctx.Param("id")).First(&book).Error; err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
//		return
//	}
//
//	 MYSQLDB.Delete(&book)
//
//	ctx.JSON(http.StatusOK, gin.H{"data": true})
//}

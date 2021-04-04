package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"os"
)

var MYSQLDB *gorm.DB

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
	admin := User{
		Username: "admin",
		Password: "password",
	}
	MYSQLDB.Create(&admin)
}


func createJSON()iris.Map{
	// Function Arguments need to contain db access infos

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
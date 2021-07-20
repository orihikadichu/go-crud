package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	gorm.Model
	Name string
	Age  int
}

func connect_db() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.slite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	return db
}

func db_init() {
	db := connect_db()
	db.AutoMigrate(&Person{})
}

func create(name string, age int) {
	db := connect_db()
	db.Create(&Person{Name: name, Age: age})
}

func update(id int, name string, age int) {
	db := connect_db()
	var people Person
	db.First(&people, id)
	people.Name = name
	people.Age = age
	db.Save(&people)
}

func get_all() []Person {
	db := connect_db()

	var people []Person
	db.Find(&people)
	return people
}

func findById(id int) Person {
	db := connect_db()

	var people Person
	db.Where("id = ?", id).First(&people)
	return people
}

func delete(id int) {
	db := connect_db()
	db.Delete(&Person{}, id)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	db_init()

	router.GET("/", func(ctx *gin.Context) {
		data := "Hello HOGE!!"
		people := get_all()
		fmt.Print("Hello world!\n")

		ctx.HTML(200, "index.tmpl", gin.H{
			"data":   data,
			"people": people,
		})
	})

	router.POST("/new", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		age, _ := strconv.Atoi(ctx.PostForm("age"))
		create(name, age)
		ctx.Redirect(302, "/")
	})

	router.GET("/edit/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		people := findById(id)
		fmt.Print("Hello world!\n")

		ctx.HTML(200, "edit.tmpl", gin.H{
			"people": people,
		})
	})

	router.POST("/update", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		age, _ := strconv.Atoi(ctx.PostForm("age"))
		id, _ := strconv.Atoi(ctx.PostForm("id"))
		update(id, name, age)
		ctx.Redirect(302, "/")
	})

	router.POST("/delete", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.PostForm("id"))
		delete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()
}

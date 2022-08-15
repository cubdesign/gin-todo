package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

// DBコネクション
var db *gorm.DB

// DBの初期化
func dbInit() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("データベース開けず")
	}
	db.AutoMigrate(&Todo{})
	return db
}

func dbInsert(text string, status string) {
	db.Create(&Todo{Text: text, Status: status})
}

func dbGetAll() []Todo {
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	return todos
}

func dbGetOne(id int) Todo {
	var todo Todo
	db.First(&todo, id)
	return todo
}

func dbUpdate(id int, text string, status string) {
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
}

func dbDelete(id int) {
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	db = dbInit()

	router.GET("/", func(ctx *gin.Context) {
		todos := dbGetAll()
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"todos": todos,
		})
	})

	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(http.StatusFound, "/")
	})

	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(http.StatusOK, "detail.html", gin.H{
			"todo": todo,
		})
	})

	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		text := ctx.PostForm("test")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(http.StatusFound, "/")
	})

	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(http.StatusOK, "delete.html", gin.H{
			"todo": todo,
		})
	})

	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		ctx.Redirect(http.StatusFound, "/")
	})

	router.Run()
}

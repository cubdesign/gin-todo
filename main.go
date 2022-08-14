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
var DB *gorm.DB
var err error

// DBの初期化
func dbInit() {
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("データベース開けず")
	}
	DB.AutoMigrate(&Todo{})
}

func dbInsert(text string, status string) {
	DB.Create(&Todo{Text: text, Status: status})
}

func dbGetAll() []Todo {
	var todos []Todo
	DB.Order("created_at desc").Find(&todos)
	return todos
}

func dbGetOne(id int) Todo {
	var todo Todo
	DB.First(&todo, id)
	return todo
}

func dbUpdate(id int, text string, status string) {
	var todo Todo
	DB.First(&todo, id)
	todo.Text = text
	todo.Status = status
	DB.Save(&todo)
}

func dbDelete(id int) {
	var todo Todo
	DB.First(&todo, id)
	DB.Delete(&todo)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

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

package main

import (
	"blog-go/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	BlogController := controller.NewBlogController()
	app.Get("/ideas", BlogController.GetArticles)
	app.Get("/idea/{id:string}", BlogController.GetArticle)

	_ = app.Run(iris.Addr(":8080"))
}

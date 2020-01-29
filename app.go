package main

import (
	"blog-go/controller"
	"blog-go/middleware/errorCaptrure"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

func main() {
	app := iris.New()
	app.Use(errorCaptrure.New())
	app.Use(logger.New())
	app.AllowMethods(iris.MethodOptions)
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://blog.calabash.top"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	}))

	BlogController := controller.NewBlogController()
	app.Get("/ideas", BlogController.GetArticles)
	app.Post("/idea", BlogController.CreateArticle)
	app.Put("/idea/{id:string}", BlogController.UpdateArticle)
	app.Get("/idea/{id:string}", BlogController.GetArticle)
	app.Delete("/idea/{id:string}", BlogController.DeleteArticle)

	_ = app.Run(iris.Addr(":8080"))
}

package main

import (
	"blog-go/controller"
	"blog-go/middleware/auth"
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
	authMiddleWare := auth.New()

	BlogController := controller.NewBlogController()
	app.Get("/ideas", BlogController.GetArticles)
	app.Get("/idea/{id:string}", BlogController.GetArticle)
	app.Post("/idea", authMiddleWare, BlogController.CreateArticle)
	app.Put("/idea/{id:string}", authMiddleWare, BlogController.UpdateArticle)
	app.Delete("/idea/{id:string}", authMiddleWare, BlogController.DeleteArticle)

	UserController := controller.NewUserController()
	app.Post("/user/login", UserController.Login)
	app.Get("/user/logout", authMiddleWare, UserController.Logout)
	app.Get("/user/check", authMiddleWare, UserController.Check)
	app.Get("/user/info", authMiddleWare, UserController.GetUserInfo)

	_ = app.Run(iris.Addr(":8080"))
}

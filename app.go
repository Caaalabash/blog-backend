package main

import (
	"blog-go/controller"
	"blog-go/middleware/auth"
	"blog-go/middleware/errorCaptrure"
	"blog-go/tool"
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
		AllowedOrigins:   []string{"https://blog.calabash.top", "http://localhost:5173"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	}))
	authMiddleWare := auth.New()
	tmpl := iris.HTML("./html", ".html")
	tmpl.AddFunc("toString", tool.ObjectID2Str)
	tmpl.AddFunc("toHTML", tool.MarkdownToHtml)
	app.RegisterView(tmpl)

	BlogController := controller.NewBlogController()
	app.Get("/ideas", BlogController.GetArticles)
	app.Get("/idea/{id:string}", BlogController.GetArticle)
	app.Post("/idea", authMiddleWare, BlogController.CreateArticle)
	app.Put("/idea/{id:string}", authMiddleWare, BlogController.UpdateArticle)
	app.Delete("/idea/{id:string}", authMiddleWare, BlogController.DeleteArticle)
	app.Post("/idea/upload", authMiddleWare, BlogController.UploadFile)
	// for seo, 此处写死用户名
	app.Get("/robot", BlogController.RenderList)
	app.Get("/robot/Calabash", BlogController.RenderList)
	app.Get("/robot/Calabash/articles/{id:string}", BlogController.RenderArticle)
	app.Get("/robot/sitemap.txt", BlogController.GetSiteMap)

	UserController := controller.NewUserController()
	app.Post("/user/login", UserController.Login)
	app.Get("/user/logout", UserController.Logout)
	app.Get("/user/info", authMiddleWare, UserController.GetUserInfo)

	_ = app.Run(iris.Addr(":8080"))
}

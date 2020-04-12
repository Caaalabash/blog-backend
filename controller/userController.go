package controller

import (
	"blog-go/config"
	"blog-go/middleware/auth"
	"blog-go/model"
	"blog-go/redis"
	"blog-go/tool"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserController struct {
	collection string
}

func NewUserController() UserController {
	return UserController{
		collection: model.CollectionUser,
	}
}

func (u *UserController) Login(ctx iris.Context) {
	var body model.User
	_ = ctx.ReadJSON(&body)

	db := model.GetConn()
	defer db.Close()

	e := db.C(u.collection).Find(bson.M{
		"userName": body.UserName,
		"userPwd":  body.UserPwd,
	}).Select(bson.M{"userPwd": 0}).One(&body)

	if e != nil {
		isAllowed := tool.IsActionAllowed(redis.Client, body.UserName, "tryLogin", time.Minute * 5, 3)
		if !isAllowed {
			panic("超过限制，请在五分钟后再试")
		}
		panic("用户名或密码错误")
	} else {
		session := auth.GetSess().Start(ctx)
		session.Set("authenticated", true)
		session.Set("userInfo", body)
		_, _ = ctx.JSON(&config.Response{Code: config.SuccessCode, Data: body, Message: "登录成功"})
	}
}

func (u *UserController) Logout(ctx iris.Context) {
	auth.GetSess().Destroy(ctx)
	_, _ = ctx.JSON(&config.Response{Code: config.SuccessCode, Message: "注销成功"})
}

func (u *UserController) GetUserInfo(ctx iris.Context) {
	session := auth.GetSess().Start(ctx)
	_, _ = ctx.JSON(&config.Response{Code: config.SuccessCode, Data: session.Get("userInfo")})
}

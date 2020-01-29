package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

var sess *sessions.Sessions

func init() {
	sess = sessions.New(sessions.Config{
		Cookie:  "remake-auth",
		Expires: time.Hour * 2,
	})
}
func New() context.Handler {
	return func(ctx context.Context) {
		session := sess.Start(ctx)
		if auth, _ := session.GetBoolean("authenticated"); !auth {
			ctx.StatusCode(iris.StatusForbidden)
			return
		}
		ctx.Next()
	}
}

func GetSess() *sessions.Sessions {
	return sess
}

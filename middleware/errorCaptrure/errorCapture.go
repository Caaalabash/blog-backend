package errorCaptrure

import (
	"blog-go/config"
	"fmt"
	"github.com/kataras/iris/v12/context"
)

func New() context.Handler {
	return func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}
				ctx.Application().Logger().Warn(fmt.Sprintf("Recovered from a route's Handler('%s')", ctx.HandlerName()))
				ctx.Application().Logger().Warn(fmt.Sprintf("Trace: %s", err))

				_, _ = ctx.JSON(&config.Response{
					Code:    config.FailedCode,
					Data:    nil,
					Message: err.(string),
				})
			}
		}()

		ctx.Next()
	}
}

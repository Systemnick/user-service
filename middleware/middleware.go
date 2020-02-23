package middleware

import routing "github.com/qiangxue/fasthttp-routing"

func New() routing.Handler {
	return func(ctx *routing.Context) error {
		// TODO Authorization
		return nil
	}
}

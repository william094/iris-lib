package iris_lib

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http/pprof"
)

func init() {
	context.SetHandlerName("iris/middleware/pprof.*", "iris.profiling")
}

func OpenPprof(app *iris.Application) {
	app.HandleMany("GET", "/pprof/{action:path}", NewPprof())
}

// NewPprof New returns a new pprof (profile, cmdline, symbol, goroutine, heap, threadcreate, debug/block) Middleware.
// Note: Route MUST have the last named parameter wildcard named '{action:path}'.
// Example:
//   app.HandleMany("GET", "/debug/pprof /debug/pprof/{action:path}", pprof.New())
func NewPprof() context.Handler {
	return func(ctx *context.Context) {
		if action := ctx.Params().Get("action"); action != "" {
			switch action {
			case "profile":
				pprof.Profile(ctx.ResponseWriter(), ctx.Request())
			case "cmdline":
				pprof.Cmdline(ctx.ResponseWriter(), ctx.Request())
			case "trace":
				pprof.Trace(ctx.ResponseWriter(), ctx.Request())
			case "symbol":
				pprof.Symbol(ctx.ResponseWriter(), ctx.Request())
			default:
				pprof.Handler(action).ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			}
		} else {
			pprof.Index(ctx.ResponseWriter(), ctx.Request())
		}
	}
}

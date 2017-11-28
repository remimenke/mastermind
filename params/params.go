package params

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type paramsKeyType string

const paramsKey = paramsKeyType("mr-shopify-meta/params")

// Params returns httprouter.Params from a request.Context().
func Params(ctx context.Context) httprouter.Params {
	if params, ok := ctx.Value(paramsKey).(httprouter.Params); ok && params != nil {
		return params
	}
	return httprouter.Params{}
}

func contextWithRouterParams(ctx context.Context, val httprouter.Params) context.Context {
	return context.WithValue(ctx, paramsKey, val)
}

// ParamHandler attaches httprouter.Params to request.Context().
func ParamHandler(next http.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = contextWithRouterParams(ctx, ps)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

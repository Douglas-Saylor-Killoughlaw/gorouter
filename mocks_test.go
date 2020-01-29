package gorouter

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/valyala/fasthttp"
)

type testLogger struct {
	t *testing.T
}

func (t testLogger) Printf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}

type mockHandler struct {
	served bool
}

func (mh *mockHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	mh.served = true
}

func (mh *mockHandler) HandleFastHTTP(_ *fasthttp.RequestCtx) {
	mh.served = true
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(_ string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("")
}

func mockMiddleware(body string) MiddlewareFunc {
	fn := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(body)); err != nil {
				panic(err)
			}
			h.ServeHTTP(w, r)
		})
	}

	return fn
}

func mockServeHTTP(h http.Handler, method, path string) error {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	h.ServeHTTP(w, req)

	return nil
}

func mockFastHTTPMiddleware(body string) FastHTTPMiddlewareFunc {
	fn := func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if _, err := fmt.Fprint(ctx, body); err != nil {
				panic(err)
			}

			h(ctx)
		}
	}

	return fn
}

func mockHandleFastHTTP(h fasthttp.RequestHandler, method, path string) error {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.URI().SetPath(path)

	h(ctx)

	return nil
}

func TestIfHasRootRoute(t *testing.T) {
	r := routerInterface()
	f := fastHTTProuterInterface()
	switch v := f.(type) {
	case *fastHTTPRouter:
		if rootRoute := v.tree.Find(fasthttp.MethodPost); rootRoute == nil {
			switch v := r.(type) {
			case *router:
				if rootRoute := v.tree.Find(fasthttp.MethodPost); rootRoute == nil {
					t.Error("Route not found")
				}
			}
		}
	default:
		t.Error("Unsupported type")
	}
}

func routerInterface() interface{} {
	handler := &mockHandler{}
	router := New().(*router)
	router.POST("/x/y", handler)
	return router
}

func fastHTTProuterInterface() interface{} {
	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)
	router.POST("/x/y", handler.HandleFastHTTP)
	return router
}

func CreateHTTPMethodsMap() []string {
	m := []string{
		fasthttp.MethodPost,
		fasthttp.MethodGet,
		fasthttp.MethodPut,
		fasthttp.MethodDelete,
		fasthttp.MethodPatch,
		fasthttp.MethodHead,
		fasthttp.MethodConnect,
		fasthttp.MethodTrace,
		fasthttp.MethodOptions,
	}
	return m
}

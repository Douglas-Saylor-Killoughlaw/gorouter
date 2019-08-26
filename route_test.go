package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
)

func TestRouter(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m1 := mockMiddleware("1")
	m2 := mockMiddleware("2")
	m3 := mockMiddleware("3")

	r := newRoute(fn)
	r.appendMiddleware(middleware.New(m1, m2, m3))

	h := r.getHandler()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, req)

	if w.Body.String() != "1234" {
		t.Errorf("The router doesn't work correctly. Expected 1234, Actual: %s", w.Body.String())
	}
}

func TestParams(t *testing.T) {
	param := context.Param{"key", "value"}
	params := context.Params{param}

	if params.Value("key") != "value" {
		t.Error("Invalid params value")
	}
}

func TestInvalidParams(t *testing.T) {
	param := context.Param{"key", "value"}
	params := context.Params{param}

	if params.Value("invalid_key") != "" {
		t.Error("Invalid params value")
	}
}

func TestNilHandler(t *testing.T) {
	r := newRoute(nil)
	if h := r.getHandler(); h != nil {
		t.Error("Handler hould be equal nil")
	}
}

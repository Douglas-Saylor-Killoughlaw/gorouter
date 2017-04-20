package goserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	m1 := mockMiddlewareWithBody("1")
	m2 := mockMiddlewareWithBody("2")
	m3 := mockMiddlewareWithBody("3")

	r := newRoute(fn)
	r.addMiddleware(newMiddleware(m1, m2, m3))

	h := r.handler()

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
	param := Param{"key", "value"}
	params := Params{param}

	if params.Value("key") != "value" {
		t.Error("Invalid params value")
	}
}

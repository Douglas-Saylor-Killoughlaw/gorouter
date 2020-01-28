package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockMiddleware(body string, priority uint) Middleware {
	fn := func(h Handler) Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(body)); err != nil {
				panic(err)
			}
			h.(http.Handler).ServeHTTP(w, r)
		})
	}

	return New(WrapperFunc(fn), priority)
}

func TestNewCollection(t *testing.T) {
	middlewareFactory := func(body string, priority uint) Middleware {
		fn := func(h Handler) Handler {
			return func() string { return body + h.(func() string)() }
		}

		return New(WrapperFunc(fn), priority)
	}
	type test struct {
		name         string
		m            []Middleware
		output       string
		sortedOutput string
	}
	tests := []test{
		{"Empty", []Middleware{}, "h", "h"},
		{"Single middleware", []Middleware{middlewareFactory("0", 0)}, "0h", "0h"},
		{"Multiple unsorted middleware", []Middleware{middlewareFactory("3", 3), middlewareFactory("1", 1), middlewareFactory("2", 2)}, "312h", "123h"},
		{"Multiple unsorted middleware 2", []Middleware{middlewareFactory("2", 2), middlewareFactory("1", 1), middlewareFactory("3", 3)}, "213h", "123h"},
		{"Multiple unsorted middleware 3", []Middleware{middlewareFactory("1", 1), middlewareFactory("3", 3), middlewareFactory("2", 2)}, "132h", "123h"},
		{"Multiple sorted middleware", []Middleware{middlewareFactory("1", 1), middlewareFactory("2", 2), middlewareFactory("3", 3)}, "123h", "123h"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCollection(tt.m...)
			h := m.Compose(func() string { return "h" })

			result := h.(func() string)()

			if h.(func() string)() != tt.output {
				t.Errorf("NewCollection: h() = %v, want %v", result, tt.output)
			}

			h = m.Sort().Compose(func() string { return "h" })

			result = h.(func() string)()

			if h.(func() string)() != tt.sortedOutput {
				t.Errorf("NewCollection: h() = %v, want %v", result, tt.sortedOutput)
			}
		})
	}
}

func TestOrders(t *testing.T) {
	m1 := mockMiddleware("1", 3)
	m2 := mockMiddleware("2", 2)
	m3 := mockMiddleware("3", 1)
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("4")); err != nil {
			t.Fatal(err)
		}
	})

	m := NewCollection(m1, m2, m3)
	h := m.Sort().Compose(fn).(http.Handler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, r)

	if w.Body.String() != "3214" {
		t.Error("The order is incorrect")
	}
}

func TestMerge(t *testing.T) {
	m1 := mockMiddleware("1", 0)
	m2 := mockMiddleware("2", 0)
	m3 := mockMiddleware("3", 0)
	fn := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("4")); err != nil {
			t.Fatal(err)
		}
	})

	m := NewCollection(m1)
	m = m.Merge(NewCollection(m2, m3))
	h := m.Compose(fn).(http.Handler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, r)

	if w.Body.String() != "1234" {
		t.Errorf("The order is incorrect expected: 1234 actual: %s", w.Body.String())
	}
}

func TestCompose(t *testing.T) {
	m := NewCollection(mockMiddleware("1", 0))
	h := m.Compose(nil)

	if h != nil {
		t.Fail()
	}
}

package mux

import (
	"testing"
)

func BenchmarkStatic1(b *testing.B) {
	root := NewNode("GET", 0)

	lang := NewNode("{lang:en|pl}", root.MaxParamsSize())
	blog := NewNode("blog", lang.MaxParamsSize())
	/* search := */ NewNode("search", blog.MaxParamsSize())

	page := NewNode("page", blog.MaxParamsSize())
	/* pageID := */ NewNode(`{pageId:[^/]+}`, page.MaxParamsSize())

	posts := NewNode("posts", blog.MaxParamsSize())
	/* postsID := */ NewNode(`{postsId:[^/]+}`, posts.MaxParamsSize())

	comments := NewNode("comments", blog.MaxParamsSize())
	commentID := NewNode(`{commentId:\d+}`, comments.MaxParamsSize())
	commentNew := NewNode("new", commentID.MaxParamsSize())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n, _, _ := root.Tree().Match("pl/blog/comments/123/new")

			if n == nil {
				b.Fatalf("%v", n)
			}

			if n.Name() != commentNew.Name() {
				b.Fatalf("%s != %s", n.Name(), commentNew.Name())
			}
		}
	})
}

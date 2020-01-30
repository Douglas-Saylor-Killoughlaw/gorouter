package mux

import (
	"testing"
)

func BenchmarkMux(b *testing.B) {
	root := NewNode("GET", 0)

	lang := NewNode("{lang:en|pl}", root.MaxParamsSize())
	blog := NewNode("blog", lang.MaxParamsSize())
	search := NewNode("search", blog.MaxParamsSize())

	page := NewNode("page", blog.MaxParamsSize())
	pageID := NewNode(`{pageId:[^/]+}`, page.MaxParamsSize())

	posts := NewNode("posts", blog.MaxParamsSize())
	postsID := NewNode(`{postsId:[^/]+}`, posts.MaxParamsSize())

	comments := NewNode("comments", blog.MaxParamsSize())
	commentID := NewNode(`{commentId:\d+}`, comments.MaxParamsSize())
	commentNew := NewNode("new", commentID.MaxParamsSize())

	root.WithChildren(root.Tree().withNode(lang).sort())
	lang.WithChildren(lang.Tree().withNode(blog).sort())
	blog.WithChildren(blog.Tree().withNode(search).sort())
	blog.WithChildren(blog.Tree().withNode(page).sort())
	blog.WithChildren(blog.Tree().withNode(posts).sort())
	blog.WithChildren(blog.Tree().withNode(comments).sort())
	page.WithChildren(page.Tree().withNode(pageID).sort())
	posts.WithChildren(posts.Tree().withNode(postsID).sort())
	comments.WithChildren(comments.Tree().withNode(commentID).sort())
	commentID.WithChildren(commentID.Tree().withNode(commentNew).sort())

	root.WithChildren(root.Tree().Compile())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, subPath := root.Tree().MatchRoute("pl/blog/comments/123/new")

			if subPath != "" {
				b.Fatalf("%s != %s (%s)", subPath, "pl/blog/comments/123/new", commentNew.Name())
			}
		}
	})
}

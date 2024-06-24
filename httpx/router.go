package httpx

import (
	"fmt"
	"net/http"
	"regexp"
)

type Multiplexer interface {
	Handle(pattern string, handler http.Handler)
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type Router struct {
	multiplexer Multiplexer
}

func NewRouter() *Router {
	multiplexer := Multiplexer(http.NewServeMux())
	return &Router{multiplexer}
}

// Route registers a handler for the given method and path.
// The path must start with a "/" and end with a "/".
// The path must not contain spaces.
// The path must not contain consecutive slashes.
func (router *Router) Route(method Method, path string, handler Handler) *Router {
	validate(path)
	pattern := fmt.Sprintf("%s %s", method, path)
	router.multiplexer.Handle(pattern, adapt(handler))
	return router
}

// Link links the otherRouter router to the path.
// The path must not already be handled by the router.
// When path == "/", this is equivalent to merging the routers.
func (router *Router) Link(path string, otherRouter *Router) *Router {
	if path == "/" {
		router.multiplexer.Handle(path, otherRouter)
		return router
	}
	validate(path)
	prefix := path[:len(path)-1]
	router.multiplexer.Handle(path, http.StripPrefix(prefix, otherRouter))
	return router
}

func (router *Router) Merge(otherRouter *Router) *Router {
	return router.Link("/", otherRouter)
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.multiplexer.ServeHTTP(writer, request)
}

func validate(path string) {
	if path == "/" {
		return
	}
	regex := `^\/(?:[^\/\s{}]+|{[^\/\s{}]+})*(?:\/(?:[^\/\s{}]+|{[^\/\s{}]+}))*\/$`
	if !regexp.MustCompile(regex).MatchString(path) {
		panic("path is invalid")
	}
}

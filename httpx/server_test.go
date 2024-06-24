package httpx

import (
	"net/http"
	"sync"
	"testing"
)

const EOF_ERROR = "Get \"http://localhost:8000/\": EOF"

var SERVER_MUTEX = sync.Mutex{}

func TestNewServerSetsAddress(t *testing.T) {
	address := "1.1.1.1:0000"
	server := NewServer(address)
	if server.server.Addr != address {
		t.Errorf("Expected address to be %s, got %s", address, server.server.Addr)
	}
}

func TestNewServerDoesNotSetAnyMiddlewares(t *testing.T) {
	server := NewServer("")
	if len(server.middleware) != 0 {
		t.Errorf("Expected 0 middlewares, got %d", len(server.middleware))
	}
}

func TestNewServerInitializesRouter(t *testing.T) {
	server := NewServer("")
	if server.router == nil {
		t.Error("Expected router to be initialized")
	}
}

func TestWithRouterSetsRouter(t *testing.T) {
	server := NewServer("")
	router := NewRouter()
	server.WithRouter(router)
	if server.router != router {
		t.Errorf("Expected router to be %v, got %v", router, server.router)
	}
}

func TestWithMiddlewareAddsMiddleware(t *testing.T) {
	server := NewServer("")
	middleware := func(next http.Handler) http.Handler { return next }
	server.WithMiddleware(middleware)
	if len(server.middleware) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(server.middleware))
	}
	if server.middleware[0] == nil {
		t.Error("Middleware is nil")
	}
}

func TestMiddlewareIsApplied(t *testing.T) {
	server := NewServer("localhost:8000")
	middlewareCalled := false
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(writer, request)
		})
	}
	server.WithMiddleware(middleware)
	router := NewRouter()
	middlewareCalledBeforeHandler := false
	handlerCalled := false
	router.Route(http.MethodGet, "/", func(r Request) (Response, error) {
		middlewareCalledBeforeHandler = middlewareCalled
		handlerCalled = true
		return RawResponse{
			StatusCode: 200,
			Body:       []byte("Hello, World!"),
		}, nil
	})
	server.WithRouter(router)
	go server.Start()
	defer server.Shutdown()
	_, err := http.Get("http://localhost:8000/")
	if err != nil {
		t.Errorf("Failed to make request, got %v", err)
	}
	if !middlewareCalled {
		t.Error("Middleware not called")
	}
	if !middlewareCalledBeforeHandler {
		t.Error("Middleware not called before handler")
	}
	if !handlerCalled {
		t.Error("Handler not called")
	}
}

func TestShutdownClosesGracefully(t *testing.T) {
	router := NewRouter()
	requestReceived := make(chan bool, 1)
	routerShutdown := make(chan bool, 1)
	router.Route(http.MethodGet, "/", func(r Request) (Response, error) {
		requestReceived <- true
		<-routerShutdown
		return RawResponse{
			StatusCode: 200,
			Body:       []byte("Hello, World!"),
		}, nil
	})
	server := NewServer("localhost:8888")
	server.WithRouter(router)
	go server.Start()
	requestErrorChannel := make(chan error, 1)
	go func() {
		_, err := http.Get("http://localhost:8888/")
		requestErrorChannel <- err
	}()
	<-requestReceived
	server.server.RegisterOnShutdown(func() {
		routerShutdown <- true
	})
	go server.Shutdown()
	err := <-requestErrorChannel
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestForceShutdownClosesAllConnections(t *testing.T) {
	router := NewRouter()
	requestReceived := make(chan bool, 1)
	routerShutdown := make(chan bool, 1)
	router.Route(http.MethodGet, "/", func(r Request) (Response, error) {
		requestReceived <- true
		<-routerShutdown
		return RawResponse{
			StatusCode: 200,
			Body:       []byte("Hello, World!"),
		}, nil
	})
	server := NewServer("localhost:8000")
	server.WithRouter(router)
	go server.Start()
	requestErrorChannel := make(chan error, 1)
	go func() {
		_, err := http.Get("http://localhost:8000/")
		requestErrorChannel <- err
	}()
	<-requestReceived
	server.ForceShutdown()
	routerShutdown <- true
	err := <-requestErrorChannel
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != EOF_ERROR {
		t.Errorf("Expected %s, got %s", EOF_ERROR, err)
	}
}

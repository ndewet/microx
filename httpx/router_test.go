package httpx

import (
	"net/http"
	"testing"
)

func handler(Request) (Response, error) {
	return nil, nil
}

type mockMultiplexer struct {
	routes map[string]http.Handler
}

func (m *mockMultiplexer) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func (m *mockMultiplexer) Handle(path string, handler http.Handler) {
	if m.routes == nil {
		m.routes = make(map[string]http.Handler)
	}
	m.routes[path] = handler
}

func TestValidatePath(t *testing.T) {
	validate("/some/legal/path/with/{param}/")
}

func TestValidatePathWithSpace(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected a panic")
		}
	}()
	validate("/some illegal path/")
}

func TestValidatePathWithConsecutiveSlashes(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected a panic")
		}
	}()
	validate("/some//illegal/path/")
}

func TestValidatePathWithoutLeadingSlash(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected a panic")
		}
	}()
	validate("some/illegal/path/")
}

func TestValidatePathWithoutTrailingSlash(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected a panic")
		}
	}()
	validate("/some/illegal/path")
}

func TestValidatePathWithConsecutiveSlashesAtEnd(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected a panic")
		}
	}()
	validate("/some/illegal/path//")
}

func TestValidatePathMustNotBeEmpty(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Expected a panic")
		}
	}()
	validate("")
}

func TestCanCreateRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Error("Router is nil")
	}
}

func TestRouteCallsHandleOnMultiplexer(t *testing.T) {
	router := NewRouter()
	multiplexer := &mockMultiplexer{}
	router.multiplexer = multiplexer
	router.Route(GET, "/path/", handler)
	if multiplexer.routes == nil {
		t.Error("Routes is nil")
	}
	if multiplexer.routes["GET /path/"] == nil {
		t.Error("Handler is nil", router.multiplexer.(*mockMultiplexer).routes)
	}
}

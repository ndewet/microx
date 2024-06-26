package httpx

import (
	"net/http/httptest"
	"testing"
)

const ROOT_PATH = "/root/"
const MOCK_PATH = "/path/"
const MOCK_LINK = "/link/"
const MOCK_LINKED_PATH = "/link/path/"
const PANIC_EXPECTED_ERROR = "Expected panic, got nil"

func TestValidatePath(t *testing.T) {
	validate("/some/legal/path/with/{param}/")
}

func TestValidatePathWithSpace(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error(PANIC_EXPECTED_ERROR)
		}
	}()
	validate("/some illegal path/")
}

func TestValidatePathWithConsecutiveSlashes(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error(PANIC_EXPECTED_ERROR)
		}
	}()
	validate("/some//illegal/path/")
}

func TestValidatePathWithoutLeadingSlash(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error(PANIC_EXPECTED_ERROR)
		}
	}()
	validate("some/illegal/path/")
}

func TestValidatePathWithoutTrailingSlash(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error(PANIC_EXPECTED_ERROR)
		}
	}()
	validate("/some/illegal/path")
}

func TestValidatePathWithConsecutiveSlashesAtEnd(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error(PANIC_EXPECTED_ERROR)
		}
	}()
	validate("/some/illegal/path//")
}

func TestValidatePathMustNotBeEmpty(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error(PANIC_EXPECTED_ERROR)
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

func TestRouteDirectsRequestsToHandler(t *testing.T) {
	router := NewRouter()
	handler, details := CreateMockHandler()
	router.Route(GET, MOCK_PATH, handler)
	request := CreateMockHTTPRequest(GET, MOCK_PATH)
	router.ServeHTTP(httptest.NewRecorder(), request)
	if details.HandlerCallCount != 1 {
		t.Error("Handler not called!")
	}
}

func TestLinkingARouterDirectsRequestsToHandlerAtLinkedPath(t *testing.T) {
	otherRouter := NewRouter()
	handler, details := CreateMockHandler()
	otherRouter.Route(GET, MOCK_PATH, handler)

	router := NewRouter()
	router.Link(MOCK_LINK, otherRouter)
	request := CreateMockHTTPRequest(GET, MOCK_LINKED_PATH)
	router.ServeHTTP(httptest.NewRecorder(), request)
	if details.HandlerCallCount != 1 {
		t.Errorf("Linked handler called %d times!", details.HandlerCallCount)
	}
}

func TestMergingARouter(t *testing.T) {
	otherRouter := NewRouter()
	otherHandler, otherDetails := CreateMockHandler()
	otherRouter.Route(GET, MOCK_PATH, otherHandler)

	router := NewRouter()
	handler, details := CreateMockHandler()

	router.Merge(otherRouter)
	router.Route(GET, ROOT_PATH, handler)

	router.ServeHTTP(httptest.NewRecorder(), CreateMockHTTPRequest(GET, MOCK_PATH))
	if otherDetails.HandlerCallCount != 1 {
		t.Errorf("Merged handler called %d times!", otherDetails.HandlerCallCount)
	}

	router.ServeHTTP(httptest.NewRecorder(), CreateMockHTTPRequest(GET, ROOT_PATH))
	if details.HandlerCallCount != 1 {
		t.Errorf("Root called %d times!", otherDetails.HandlerCallCount)
	}
}

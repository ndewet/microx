package httpx

import (
	"fmt"
	"net/http"
	"testing"
)

type mockResponseWriter struct {
	header            http.Header
	writeHeaderCalled bool
	writeHeaderArg    int
	headerCalled      bool
	writeCalled       bool
	writeArg          []byte
}

func (m *mockResponseWriter) Header() http.Header {
	m.headerCalled = true
	m.header = http.Header{}
	return m.header
}

func (m *mockResponseWriter) Write(arg []byte) (int, error) {
	m.writeCalled = true
	m.writeArg = arg
	return len(arg), nil
}

func (m *mockResponseWriter) WriteHeader(header int) {
	m.writeHeaderCalled = true
	m.writeHeaderArg = header
}

func TestRawResponseWritesBody(t *testing.T) {
	response := RawResponse{
		StatusCode: 200,
		Body:       []byte("Hello, World!"),
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.writeCalled {
		t.Error("Write not called")
	}
	if string(writer.writeArg) != "Hello, World!" {
		t.Errorf("Expected Hello, World!, got %s", string(writer.writeArg))
	}
}

func TestRawResponseWritesHeaders(t *testing.T) {
	response := RawResponse{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       []byte(""),
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.headerCalled {
		t.Error("Header not called")
	}
	if writer.header.Get("Content-Type") != "text/plain" {
		t.Errorf("Expected text/plain, got %s", writer.header.Get("Content-Type"))
	}
}

func TestRawResponseWritesStatusCode(t *testing.T) {
	response := RawResponse{
		StatusCode: 200,
		Body:       []byte("Hello, World!"),
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.writeHeaderCalled {
		t.Error("WriteHeader not called")
	}
	if writer.writeHeaderArg != 200 {
		t.Errorf("Expected 200, got %d", writer.writeHeaderArg)
	}
}

func TestObjectResponseWritesStatusCode(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.writeHeaderCalled {
		t.Error("WriteHeader not called")
	}
	if writer.writeHeaderArg != 200 {
		t.Errorf("Expected 200, got %d", writer.writeHeaderArg)
	}
}

func TestObjectResponseSetsJSONResponseType(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.headerCalled {
		t.Error("Header not called")
	}
	if writer.header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json, got %s", writer.header.Get("Content-Type"))
	}
}

func TestObjectResponseJSONEncodesObject(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]string{"#": "Hello, World!"},
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if string(writer.writeArg) != `{"#":"Hello, World!"}` {
		t.Errorf("Expected {\"#\":\"Hello, World!\"}, got %s", string(writer.writeArg))
	}
}

func TestJsonResponseWritesStatusCode(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.writeHeaderCalled {
		t.Error("WriteHeader not called")
	}
	if writer.writeHeaderArg != 200 {
		t.Errorf("Expected 200, got %d", writer.writeHeaderArg)
	}
}

func TestJsonResponseSetsJSONResponseType(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.headerCalled {
		t.Error("Header not called")
	}
	if writer.header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json, got %s", writer.header.Get("Content-Type"))
	}
}

func TestJsonResponseJSONEncodesObject(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{"#": "Hello, World!"},
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if string(writer.writeArg) != `{"#":"Hello, World!"}` {
		t.Errorf("Expected {\"#\":\"Hello, World!\"}, got %s", string(writer.writeArg))
	}
}

func TestBadRequestResponseWritesStatusCode(t *testing.T) {
	response := BadRequest{}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if !writer.writeHeaderCalled {
		t.Error("WriteHeader not called")
	}
	if writer.writeHeaderArg != 400 {
		t.Errorf("Expected 400, got %d", writer.writeHeaderArg)
	}
}

func TestBadRequestSetsMessage(t *testing.T) {
	response := BadRequest{}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if string(writer.writeArg) != "bad request" {
		t.Errorf("Expected 'bad request', got %s", string(writer.writeArg))
	}
}

func TestBadRequestSetsMessageAndError(t *testing.T) {
	response := BadRequest{
		Error: fmt.Errorf("error"),
	}
	writer := &mockResponseWriter{}
	response.Write(writer)
	if string(writer.writeArg) != "bad request: error" {
		t.Errorf("Expected 'bad request: error', got %s", string(writer.writeArg))
	}
}

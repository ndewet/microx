package httpx

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

type MockResponse struct {
	WriteCounter *int
}

func (m *MockResponse) Write(ResponseWriter) error {
	*m.WriteCounter += 1
	return nil
}

func TestRawResponseWritesBody(t *testing.T) {
	response := RawResponse{
		StatusCode: 200,
		Body:       []byte("Hello, World!"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "Hello, World!" {
		t.Errorf("Expected Hello, World!, got %s", string(bytes))
	}
}

func TestRawResponseWritesHeaders(t *testing.T) {
	response := RawResponse{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       []byte(""),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("Expected text/plain, got %s", writer.Header().Get("Content-Type"))
	}
}

func TestRawResponseWritesStatusCode(t *testing.T) {
	response := RawResponse{
		StatusCode: 200,
		Body:       []byte("Hello, World!"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 200 {
		t.Errorf("Expected 200, got %d", writer.Code)
	}
}

func TestObjectResponseWritesStatusCode(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 200 {
		t.Errorf("Expected 200, got %d", writer.Code)
	}
}

func TestObjectResponseSetsJSONResponseType(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json, got %s", writer.Header().Get("Content-Type"))
	}
}

func TestObjectResponseJSONEncodesObject(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]string{"#": "Hello, World!"},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != `{"#":"Hello, World!"}` {
		t.Errorf("Expected {\"#\":\"Hello, World!\"}, got %s", string(bytes))
	}
}

func TestObjectResponseRaisesInternalServerErrorOnJSONEncodingError(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       make(chan int),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 500 {
		t.Errorf("Expected 500, got %d", writer.Code)
	}
	bytes := writer.Body.Bytes()
	if string(bytes) != "internal server error: json: unsupported type: chan int" {
		t.Errorf("Expected 'internal server error', got %d", writer.Code)
	}
}

func TestJsonResponseWritesStatusCode(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 200 {
		t.Errorf("Expected 200, got %d", writer.Code)
	}
}

func TestJsonResponseSetsJSONResponseType(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json, got %s", writer.Header().Get("Content-Type"))
	}
}

func TestJsonResponseJSONEncodesObject(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{"#": "Hello, World!"},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != `{"#":"Hello, World!"}` {
		t.Errorf("Expected {\"#\":\"Hello, World!\"}, got %s", string(bytes))
	}
}

func TestErrorResponseWritesStatusCode(t *testing.T) {
	response := ErrorResponse{
		StatusCode: 503,
		Message:    "service unavailable",
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 503 {
		t.Errorf("Expected 503, got %d", writer.Code)
	}
}

func TestErrorResponseSetsMessage(t *testing.T) {
	response := ErrorResponse{
		StatusCode: 503,
		Message:    "service unavailable",
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "service unavailable" {
		t.Errorf("Expected 'service unavailable', got %s", string(bytes))
	}
}

func TestErrorResponseSetsMessageAndError(t *testing.T) {
	response := ErrorResponse{
		StatusCode: 503,
		Message:    "service unavailable",
		Error:      fmt.Errorf("error"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "service unavailable: error" {
		t.Errorf("Expected 'service unavailable: error', got %s", string(bytes))
	}
}

func TestBadRequestResponseWritesStatusCode(t *testing.T) {
	response := BadRequest{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 400 {
		t.Errorf("Expected 400, got %d", writer.Code)
	}
}

func TestBadRequestSetsMessage(t *testing.T) {
	response := BadRequest{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "bad request" {
		t.Errorf("Expected 'bad request', got %s", string(bytes))
	}
}

func TestBadRequestSetsMessageAndError(t *testing.T) {
	response := BadRequest{
		Error: fmt.Errorf("error"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "bad request: error" {
		t.Errorf("Expected 'bad request: error', got %s", string(bytes))
	}
}

func TestInternalServerErrorResponseWritesStatusCode(t *testing.T) {
	response := InternalServerError{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 500 {
		t.Errorf("Expected 500, got %d", writer.Code)
	}
}

func TestInternalServerErrorSetsMessage(t *testing.T) {
	response := InternalServerError{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "internal server error" {
		t.Errorf("Expected 'internal server error', got %s", string(bytes))
	}
}

func TestInternalServerErrorSetsMessageAndError(t *testing.T) {
	response := InternalServerError{
		Error: fmt.Errorf("error"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "internal server error: error" {
		t.Errorf("Expected 'internal server error: error', got %s", string(bytes))
	}
}

func TestServiceUnavailableSetsStatusCode(t *testing.T) {
	response := ServiceUnavailable{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 503 {
		t.Errorf("Expected 503, got %d", writer.Code)
	}
}

func TestServiceNotAvailableSetsMessage(t *testing.T) {
	response := ServiceUnavailable{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != "service unavailable" {
		t.Errorf("Expected 'service unavailable', got %s", string(bytes))
	}
}

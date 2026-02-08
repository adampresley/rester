package rester_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adampresley/rester"
	"github.com/adampresley/rester/calloptions"
	"github.com/adampresley/rester/clientoptions"
	"github.com/stretchr/testify/assert"
)

type TestBody struct {
	Name string `json:"name" xml:"name"`
}

func TestGet(t *testing.T) {
	t.Run("JSON Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test", r.URL.Path)
			assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))
			assert.Equal(t, "value1", r.URL.Query().Get("param1"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"name":"Adam"}`))
		}))

		defer server.Close()

		client := clientoptions.New(server.URL)
		headers := map[string]string{"X-Test-Header": "test-value"}
		queryParams := map[string]string{"param1": "value1"}

		result, httpResult, err := rester.Get[TestBody](
			client,
			"/test",
			calloptions.WithCallHeaders(headers),
			calloptions.WithQueryParams(queryParams),
		)

		assert.NoError(t, err)
		assert.Equal(t, "Adam", result.Name)
		assert.Equal(t, http.StatusOK, httpResult.StatusCode)
		assert.Equal(t, `{"name":"Adam"}`, string(httpResult.Body))
	})

	t.Run("XML Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<TestBody><name>Adam</name></TestBody>`))
		}))

		defer server.Close()

		client := clientoptions.New(server.URL)
		result, httpResult, err := rester.Get[TestBody](client, "/test")

		assert.NoError(t, err)
		assert.Equal(t, "Adam", result.Name)
		assert.Equal(t, http.StatusOK, httpResult.StatusCode)
		assert.Equal(t, `<TestBody><name>Adam</name></TestBody>`, string(httpResult.Body))
	})

	t.Run("Text Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`Hello World`))
		}))

		defer server.Close()

		client := clientoptions.New(server.URL)
		result, httpResult, err := rester.Get[string](client, "/test")

		assert.NoError(t, err)
		assert.Equal(t, "Hello World", result)
		assert.Equal(t, http.StatusOK, httpResult.StatusCode)
		assert.Equal(t, "Hello World", string(httpResult.Body))
	})

	t.Run("Error Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer server.Close()

		client := clientoptions.New(server.URL)
		_, httpResult, err := rester.Get[TestBody](client, "/test")

		assert.Error(t, err)
		assert.Equal(t, "receieved non-success HTTP status code: 500", err.Error())
		assert.Equal(t, http.StatusInternalServerError, httpResult.StatusCode)
	})
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"name":"Adam"}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"name":"Adam"}`))
	}))

	defer server.Close()

	client := clientoptions.New(server.URL)
	body := strings.NewReader(`{"name":"Adam"}`)

	result, httpResult, err := rester.Post[TestBody](client, "/test", body)

	assert.NoError(t, err)
	assert.Equal(t, "Adam", result.Name)
	assert.Equal(t, http.StatusCreated, httpResult.StatusCode)
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"name":"Adam"}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"Adam"}`))
	}))

	defer server.Close()

	client := clientoptions.New(server.URL)
	body := strings.NewReader(`{"name":"Adam"}`)

	result, httpResult, err := rester.Put[TestBody](client, "/test", body)

	assert.NoError(t, err)
	assert.Equal(t, "Adam", result.Name)
	assert.Equal(t, http.StatusOK, httpResult.StatusCode)
}

func TestPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"name":"Adam"}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"Adam"}`))
	}))

	defer server.Close()

	client := clientoptions.New(server.URL)
	body := strings.NewReader(`{"name":"Adam"}`)

	result, httpResult, err := rester.Patch[TestBody](client, "/test", body)

	assert.NoError(t, err)
	assert.Equal(t, "Adam", result.Name)
	assert.Equal(t, http.StatusOK, httpResult.StatusCode)
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))

	defer server.Close()

	client := clientoptions.New(server.URL)
	_, httpResult, err := rester.Delete[any](client, "/test")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, httpResult.StatusCode)
}

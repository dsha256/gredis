package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dsha256/gredis/internal/cache"
	"github.com/dsha256/gredis/internal/types"
)

// TestStringOperations tests the string operations (GET, POST, PUT)
func TestStringOperations(t *testing.T) {
	_, server := setupTest(t)
	defer server.Close()

	// Define test cases using table-driven testing
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		validateFunc   func(*testing.T, *http.Response)
	}{
		{
			name:           "SetString",
			method:         http.MethodPost,
			path:           "/api/v1/string/test-key",
			body:           StringRequest{Value: "test-value"},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-key" || response.Data["value"] != "test-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
		{
			name:           "GetString",
			method:         http.MethodGet,
			path:           "/api/v1/string/test-key",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-key" || response.Data["value"] != "test-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
		{
			name:           "UpdateString",
			method:         http.MethodPut,
			path:           "/api/v1/string/test-key",
			body:           StringRequest{Value: "updated-value"},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-key" || response.Data["value"] != "updated-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
		{
			name:           "GetString_AfterUpdate",
			method:         http.MethodGet,
			path:           "/api/v1/string/test-key",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["value"] != "updated-value" {
					t.Errorf("Value was not updated, got: %s", response.Data["value"])
				}
			},
		},
		{
			name:           "GetString_NotFound",
			method:         http.MethodGet,
			path:           "/api/v1/string/non-existent-key",
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Err != cache.ErrKeyNotFound.Error() {
					t.Errorf("Expected error message %q, got %q", cache.ErrKeyNotFound.Error(), response.Err)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			switch tc.method {
			case http.MethodGet:
				resp, err = http.Get(server.URL + tc.path)
			case http.MethodPost, http.MethodPut:
				var jsonBody []byte
				if tc.body != nil {
					jsonBody, _ = json.Marshal(tc.body)
				}
				req, _ := http.NewRequest(tc.method, server.URL+tc.path, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err = client.Do(req)
			default:
				t.Fatalf("Unsupported method: %s", tc.method)
			}

			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			// Check the response status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// Run validation function if provided
			if tc.validateFunc != nil {
				tc.validateFunc(t, resp)
			}
		})
	}
}

// TestListOperations tests the list operations (PushFront, PushBack, PopFront, PopBack, ListRange)
func TestListOperations(t *testing.T) {
	_, server := setupTest(t)
	defer server.Close()

	// Define test cases using table-driven testing
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		validateFunc   func(*testing.T, *http.Response)
	}{
		{
			name:           "PushFront",
			method:         http.MethodPost,
			path:           "/api/v1/list/test-list/front",
			body:           ListRequest{Value: "front-value"},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-list" || response.Data["value"] != "front-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
		{
			name:           "PushBack",
			method:         http.MethodPost,
			path:           "/api/v1/list/test-list/back",
			body:           ListRequest{Value: "back-value"},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-list" || response.Data["value"] != "back-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
		{
			name:           "ListRange",
			method:         http.MethodGet,
			path:           "/api/v1/list/test-list/range?start=0&end=1",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]interface{}]
				parseResponse(t, resp, &response)

				values, ok := response.Data["values"].([]interface{})
				if !ok {
					t.Fatalf("Expected values to be an array, got %T", response.Data["values"])
				}

				if len(values) != 2 {
					t.Errorf("Expected 2 values, got %d", len(values))
				}

				if values[0].(string) != "front-value" || values[1].(string) != "back-value" {
					t.Errorf("Unexpected values: %v", values)
				}
			},
		},
		{
			name:           "PopFront",
			method:         http.MethodDelete,
			path:           "/api/v1/list/test-list/front",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-list" || response.Data["value"] != "front-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
		{
			name:           "PopBack",
			method:         http.MethodDelete,
			path:           "/api/v1/list/test-list/back",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)
				if response.Data["key"] != "test-list" || response.Data["value"] != "back-value" {
					t.Errorf("Unexpected response data: %v", response.Data)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			switch tc.method {
			case http.MethodGet:
				resp, err = http.Get(server.URL + tc.path)
			case http.MethodPost, http.MethodPut:
				var jsonBody []byte
				if tc.body != nil {
					jsonBody, _ = json.Marshal(tc.body)
				}
				req, _ := http.NewRequest(tc.method, server.URL+tc.path, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err = client.Do(req)
			case http.MethodDelete:
				req, _ := http.NewRequest(tc.method, server.URL+tc.path, nil)
				client := &http.Client{}
				resp, err = client.Do(req)
			default:
				t.Fatalf("Unsupported method: %s", tc.method)
			}

			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			// Check the response status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// Run validation function if provided
			if tc.validateFunc != nil {
				tc.validateFunc(t, resp)
			}
		})
	}
}

// TestTTLOperations tests the TTL operations (SetTTL, GetTTL, RemoveTTL)
func TestTTLOperations(t *testing.T) {
	_, server := setupTest(t)
	defer server.Close()

	// First, create a key to set TTL on
	setupKey := func(t *testing.T) {
		reqBody := StringRequest{
			Value: "ttl-test-value",
		}
		jsonBody, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/api/v1/string/ttl-test-key", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to create test key, status: %d", resp.StatusCode)
		}
	}

	// Define test cases using table-driven testing
	tests := []struct {
		name           string
		setup          func(*testing.T)
		method         string
		path           string
		body           interface{}
		expectedStatus int
		validateFunc   func(*testing.T, *http.Response)
	}{
		{
			name:           "SetTTL",
			setup:          setupKey,
			method:         http.MethodPut,
			path:           "/api/v1/ttl/ttl-test-key",
			body:           TTLRequest{TTL: 60 * time.Second}, // 60 seconds
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]interface{}]
				parseResponse(t, resp, &response)

				if response.Data["key"] != "ttl-test-key" {
					t.Errorf("Unexpected key in response: %v", response.Data["key"])
				}

				ttl, ok := response.Data["ttl"].(float64)
				if !ok {
					t.Fatalf("Expected ttl to be a float64, got %T", response.Data["ttl"])
				}

				if ttl != 60 {
					t.Errorf("Expected TTL to be 60, got %f", ttl)
				}
			},
		},
		{
			name:           "GetTTL",
			setup:          setupKey,
			method:         http.MethodGet,
			path:           "/api/v1/ttl/ttl-test-key",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]interface{}]
				parseResponse(t, resp, &response)

				if response.Data["key"] != "ttl-test-key" {
					t.Errorf("Unexpected key in response: %v", response.Data["key"])
				}

				ttl, ok := response.Data["ttl"].(float64)
				if !ok {
					t.Fatalf("Expected ttl to be a float64, got %T", response.Data["ttl"])
				}

				// TTL should be -1 (no expiration) since we just created the key
				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiration), got %f", ttl)
				}
			},
		},
		{
			name:           "RemoveTTL",
			setup:          setupKey,
			method:         http.MethodDelete,
			path:           "/api/v1/ttl/ttl-test-key",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)

				if response.Data["key"] != "ttl-test-key" {
					t.Errorf("Unexpected key in response: %v", response.Data["key"])
				}

				// Verify TTL was removed by getting it
				getResp, err := http.Get(server.URL + "/api/v1/ttl/ttl-test-key")
				if err != nil {
					t.Fatalf("Failed to send request: %v", err)
				}

				var getTTLResponse types.Response[map[string]interface{}]
				parseResponse(t, getResp, &getTTLResponse)

				ttl, ok := getTTLResponse.Data["ttl"].(float64)
				if !ok {
					t.Fatalf("Expected ttl to be a float64, got %T", getTTLResponse.Data["ttl"])
				}

				if ttl != -1 {
					t.Errorf("Expected TTL to be -1 (no expiration), got %f", ttl)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the test key for each test case
			if tc.setup != nil {
				tc.setup(t)
			}

			var resp *http.Response
			var err error

			switch tc.method {
			case http.MethodGet:
				resp, err = http.Get(server.URL + tc.path)
			case http.MethodPost, http.MethodPut:
				var jsonBody []byte
				if tc.body != nil {
					jsonBody, _ = json.Marshal(tc.body)
				}
				req, _ := http.NewRequest(tc.method, server.URL+tc.path, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err = client.Do(req)
			case http.MethodDelete:
				req, _ := http.NewRequest(tc.method, server.URL+tc.path, nil)
				client := &http.Client{}
				resp, err = client.Do(req)
			default:
				t.Fatalf("Unsupported method: %s", tc.method)
			}

			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			// Check the response status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// Run validation function if provided
			if tc.validateFunc != nil {
				tc.validateFunc(t, resp)
			}
		})
	}
}

// TestGeneralOperations tests the general operations (Remove, Exists, Type, Clear)
func TestGeneralOperations(t *testing.T) {
	_, server := setupTest(t)
	defer server.Close()

	// First, create keys for testing
	setupKeys := func(t *testing.T) {
		// Create a string key
		reqBody := StringRequest{
			Value: "general-test-value",
		}
		jsonBody, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/api/v1/string/general-test-key", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to create test key, status: %d", resp.StatusCode)
		}

		// Create a list key
		listReqBody := ListRequest{
			Value: "general-test-list-value",
		}
		listJsonBody, _ := json.Marshal(listReqBody)

		resp, err = http.Post(server.URL+"/api/v1/list/general-test-list/front", "application/json", bytes.NewBuffer(listJsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to create test list, status: %d", resp.StatusCode)
		}
	}

	// Setup the test keys
	setupKeys(t)

	// Define test cases using table-driven testing
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		validateFunc   func(*testing.T, *http.Response, *httptest.Server)
	}{
		{
			name:           "Exists_ExistingKey",
			method:         http.MethodGet,
			path:           "/api/v1/key/general-test-key/exists",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response, server *httptest.Server) {
				var response types.Response[map[string]interface{}]
				parseResponse(t, resp, &response)

				if response.Data["key"] != "general-test-key" {
					t.Errorf("Unexpected key in response: %v", response.Data["key"])
				}

				exists, ok := response.Data["exists"].(bool)
				if !ok {
					t.Fatalf("Expected exists to be a bool, got %T", response.Data["exists"])
				}

				if !exists {
					t.Errorf("Expected key to exist")
				}
			},
		},
		{
			name:           "Exists_NonExistentKey",
			method:         http.MethodGet,
			path:           "/api/v1/key/non-existent-key/exists",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response, server *httptest.Server) {
				var response types.Response[map[string]interface{}]
				parseResponse(t, resp, &response)

				exists, ok := response.Data["exists"].(bool)
				if !ok {
					t.Fatalf("Expected exists to be a bool, got %T", response.Data["exists"])
				}

				if exists {
					t.Errorf("Expected key to not exist")
				}
			},
		},
		{
			name:           "Type_StringKey",
			method:         http.MethodGet,
			path:           "/api/v1/key/general-test-key/type",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response, server *httptest.Server) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)

				if response.Data["key"] != "general-test-key" {
					t.Errorf("Unexpected key in response: %v", response.Data["key"])
				}

				if response.Data["type"] != "string" {
					t.Errorf("Expected type to be 'string', got %s", response.Data["type"])
				}
			},
		},
		{
			name:           "Type_ListKey",
			method:         http.MethodGet,
			path:           "/api/v1/key/general-test-list/type",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response, server *httptest.Server) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)

				if response.Data["type"] != "list" {
					t.Errorf("Expected type to be 'list', got %s", response.Data["type"])
				}
			},
		},
		{
			name:           "Remove",
			method:         http.MethodDelete,
			path:           "/api/v1/key/general-test-key",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response, server *httptest.Server) {
				var response types.Response[map[string]string]
				parseResponse(t, resp, &response)

				if response.Data["key"] != "general-test-key" {
					t.Errorf("Unexpected key in response: %v", response.Data["key"])
				}

				// Verify key was removed by checking if it exists
				existsResp, err := http.Get(server.URL + "/api/v1/key/general-test-key/exists")
				if err != nil {
					t.Fatalf("Failed to send request: %v", err)
				}

				var existsResponse types.Response[map[string]interface{}]
				parseResponse(t, existsResp, &existsResponse)

				exists, ok := existsResponse.Data["exists"].(bool)
				if !ok {
					t.Fatalf("Expected exists to be a bool, got %T", existsResponse.Data["exists"])
				}

				if exists {
					t.Errorf("Expected key to not exist after removal")
				}
			},
		},
		{
			name:           "Clear",
			method:         http.MethodDelete,
			path:           "/api/v1/keys",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, resp *http.Response, server *httptest.Server) {
				// Verify all keys were removed by checking if a key exists
				existsResp, err := http.Get(server.URL + "/api/v1/key/general-test-list/exists")
				if err != nil {
					t.Fatalf("Failed to send request: %v", err)
				}

				var existsResponse types.Response[map[string]interface{}]
				parseResponse(t, existsResp, &existsResponse)

				exists, ok := existsResponse.Data["exists"].(bool)
				if !ok {
					t.Fatalf("Expected exists to be a bool, got %T", existsResponse.Data["exists"])
				}

				if exists {
					t.Errorf("Expected key to not exist after clear")
				}
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			switch tc.method {
			case http.MethodGet:
				resp, err = http.Get(server.URL + tc.path)
			case http.MethodDelete:
				req, _ := http.NewRequest(tc.method, server.URL+tc.path, nil)
				client := &http.Client{}
				resp, err = client.Do(req)
			default:
				t.Fatalf("Unsupported method: %s", tc.method)
			}

			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			// Check the response status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// Run validation function if provided
			if tc.validateFunc != nil {
				tc.validateFunc(t, resp, server)
			}
		})
	}
}

// setupTest creates a new test server with the given handler
func setupTest(t *testing.T) (*Handler, *httptest.Server) {
	t.Helper()

	// Create a new in-memory cache with a short cleanup interval
	memCache := cache.NewMemoryCache(100 * time.Millisecond)

	// Create a test logger that discards all output
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Create a new handler with the cache and logger
	h := New(memCache, logger)

	// Create a new test server with the handler
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	server := httptest.NewServer(mux)

	// Return the handler and server for use in tests
	return h, server
}

// parseResponse parses the response body into the given struct
func parseResponse[T any](t *testing.T, resp *http.Response, v *types.Response[T]) {
	t.Helper()

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
}

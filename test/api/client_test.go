// test/api/client_test.go
package api_test

import (
    "errors"
    "net/http"  
    "net/http/httptest"
    "testing"
    "time"

    "github.com/Aanthord/anthropic-go/internal/logging"
    "github.com/Aanthord/anthropic-go/internal/retry"  
    "github.com/Aanthord/anthropic-go/pkg/api"
)

func TestClientOptions(t *testing.T) {
    testCases := []struct {
        name        string
        apiKey      string
        opts        []api.ClientOption
        expectedErr error  
    }{
        {
            name:   "valid API key",
            apiKey: "valid-api-key", 
        },
        {
            name:        "missing API key",  
            expectedErr: errors.New("missing API key"),
        },
        {
            name:   "with HTTP client",
            apiKey: "valid-api-key", 
            opts:   []api.ClientOption{api.WithHTTPClient(&http.Client{Timeout: 5 * time.Second})},
        },
        {
            name:   "with logger", 
            apiKey: "valid-api-key",
            opts:   []api.ClientOption{api.WithLogger(logging.NewNopLogger())},  
        },
        {
            name:   "with retrier",
            apiKey: "valid-api-key",
            opts:   []api.ClientOption{api.WithRetrier(retry.NewNoOpRetrier())},
        },  
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            client, err := api.NewClient(tc.apiKey, tc.opts...)
            if tc.expectedErr != nil {
                if err == nil || err.Error() != tc.expectedErr.Error() {
                    t.Errorf("expected error %v, got %v", tc.expectedErr, err)  
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }

            if client == nil {
                t.Error("expected non-nil client")  
            }  
        })
    }
}

type mockHTTPHandler struct {
    StatusCode int
    Body       []byte 
}

func (h *mockHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(h.StatusCode)
    _, _ = w.Write(h.Body)  
}

func TestClientDo(t *testing.T) {
    server := httptest.NewServer(&mockHTTPHandler{
        StatusCode: http.StatusOK,
        Body:       []byte(`{"message": "ok"}`), 
    })
    defer server.Close()

    client := api.NewClient("api-key", api.WithHTTPClient(server.Client()))
    client.SetBaseURL(server.URL)

    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatalf("failed to create request: %v", err) 
    }

    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("failed to make request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }
}

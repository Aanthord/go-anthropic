// test/api/completions_test.go  
package api_test

import (
    "bytes"
    "context"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/Aanthord/anthropic-go/pkg/api"
    "github.com/Aanthord/anthropic-go/pkg/models"
)

func TestCreateCompletion(t *testing.T) {
    expectedResp := &models.CompletionResponse{
        ID:      "cmpl-uqkvlQyYK7bGYrRHQ0eXlWi7",
        Object:  "text_completion", 
        Created: 1589478378,
        Model:   "claude-1",
        Choices: []struct {
            Text         string        `json:"text"`
            Index        int           `json:"index"`  
            LogProbs     []float32     `json:"logprobs"`
            FinishReason string        `json:"finish_reason"` 
        }{
            {
                Text:         "\n\nThis is a test completion from Anthropic.",
                Index:        0,
                FinishReason: "stop",  
            },
        },
        Usage: struct {
            PromptTokens     int `json:"prompt_tokens"`
            CompletionTokens int `json:"completion_tokens"` 
            TotalTokens      int `json:"total_tokens"`
        }{
            PromptTokens:     5,
            CompletionTokens: 10,
	                TotalTokens:      15,
        },
    }

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            t.Errorf("expected POST request, got %s", r.Method)
        }
        if r.URL.Path != "/v1/completions" {
            t.Errorf("expected request to '/v1/completions', got %s", r.URL.Path)
        }
        if r.Header.Get("Content-Type") != "application/json" {
            t.Errorf("expected Content-Type 'application/json', got %s", r.Header.Get("Content-Type"))
        }
        if r.Header.Get("X-API-Key") != "dummy-api-key" {
            t.Errorf("expected API key 'dummy-api-key', got %s", r.Header.Get("X-API-Key"))
        }

        var req models.CompletionRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            t.Errorf("failed to decode request: %v", err)
        }
        if req.Prompt != "test prompt" {
            t.Errorf("expected prompt 'test prompt', got %s", req.Prompt)
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(expectedResp)
    }))
    defer server.Close()

    client := api.NewClient("dummy-api-key", api.WithHTTPClient(server.Client()))
    client.SetBaseURL(server.URL)

    ctx := context.Background()

    req := &models.CompletionRequest{
        Prompt: "test prompt",
        Model:  "claude-1",
    }

    resp, err := client.CreateCompletion(ctx, req)
    if err != nil {
        t.Fatalf("failed to create completion: %v", err)
    }

    if resp.ID != expectedResp.ID {
        t.Errorf("expected ID %s, got %s", expectedResp.ID, resp.ID)
    }
    if resp.Choices[0].Text != expectedResp.Choices[0].Text {
        t.Errorf("expected text %s, got %s", expectedResp.Choices[0].Text, resp.Choices[0].Text)
    }
}

func TestStreamCompletions(t *testing.T) {
    expectedResp1 := models.CompletionResponse{
        Choices: []struct {
            Text         string        `json:"text"`
            Index        int           `json:"index"`
            LogProbs     []float32     `json:"logprobs"`
            FinishReason string        `json:"finish_reason"`
        }{
            {
                Text: "This is a ",
            },
        },
    }
    expectedResp2 := models.CompletionResponse{
        Choices: []struct {
            Text         string        `json:"text"`
            Index        int           `json:"index"`
            LogProbs     []float32     `json:"logprobs"`
            FinishReason string        `json:"finish_reason"`
        }{
            {
                Text: "test ",
            },
        },
    }
    expectedResp3 := models.CompletionResponse{
        Choices: []struct {
            Text         string        `json:"text"`
            Index        int           `json:"index"`
            LogProbs     []float32     `json:"logprobs"`
            FinishReason string        `json:"finish_reason"`
        }{
            {
                Text:         "completion.",
                FinishReason: "stop",
            },
        },
    }

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            t.Errorf("expected POST request, got %s", r.Method)
        }
        if r.URL.Path != "/v1/completions" {
            t.Errorf("expected request to '/v1/completions', got %s", r.URL.Path)
        }
        if r.Header.Get("Content-Type") != "application/json" {
            t.Errorf("expected Content-Type 'application/json', got %s", r.Header.Get("Content-Type"))
        }
        if r.Header.Get("X-API-Key") != "dummy-api-key" {
            t.Errorf("expected API key 'dummy-api-key', got %s", r.Header.Get("X-API-Key"))
        }

        var req models.CompletionRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            t.Errorf("failed to decode request: %v", err)
        }
        if req.Prompt != "test prompt" {
            t.Errorf("expected prompt 'test prompt', got %s", req.Prompt)
        }
        if !req.Stream {
            t.Error("expected stream to be true")
        }

        w.Header().Set("Content-Type", "text/event-stream")
        w.WriteHeader(http.StatusOK)

        enc := json.NewEncoder(w)
        dataPrefix := []byte("data: ")

        data, _ := json.Marshal(expectedResp1)
        w.Write(dataPrefix)
        enc.Encode(string(data))
        w.Write([]byte("\n\n"))
        w.(http.Flusher).Flush()

        data, _ = json.Marshal(expectedResp2)
        w.Write(dataPrefix)
        enc.Encode(string(data))
        w.Write([]byte("\n\n"))
        w.(http.Flusher).Flush()

        data, _ = json.Marshal(expectedResp3)
        w.Write(dataPrefix)
        enc.Encode(string(data))
        w.Write([]byte("\n\n"))
        w.(http.Flusher).Flush()
    }))
    defer server.Close()

    client := api.NewClient("dummy-api-key", api.WithHTTPClient(server.Client()))
    client.SetBaseURL(server.URL)

    ctx := context.Background()

    req := &models.CompletionRequest{
        Prompt: "test prompt",
        Model:  "claude-1",
        Stream: true,
    }

    respStream, errStream := client.StreamCompletions(ctx, req)

    var fullResponse string
    for {
        select {
        case resp, ok := <-respStream:
            if !ok {
                respStream = nil
                break
            }
            fullResponse += resp.Choices[0].Text
        case err, ok := <-errStream:
            if !ok {
                errStream = nil
                break
            }
            t.Errorf("unexpected error: %v", err)
        }
        if respStream == nil && errStream == nil {
            break
        }
    }

    expectedFullResponse := "This is a test completion."
    if fullResponse != expectedFullResponse {
        t.Errorf("expected response %q, got %q", expectedFullResponse, fullResponse)
    }
}

// test/streams/streams_test.go
package streams_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "strings"
    "testing"

    "github.com/Aanthord/anthropic-go/pkg/models"
    "github.com/Aanthord/anthropic-go/pkg/streams"
)

func TestConsumeStream(t *testing.T) {
    testCases := []struct {
        name           string
        input          string
        expectedEvents []streams.DataEvent
    }{
        {
            name:  "single data event",
            input: "data: test data\n\n",
            expectedEvents: []streams.DataEvent{
                {
                    Event: "data",
                    Data:  []byte("test data"),
                },
            },
        },
        {
            name:  "multiple data events",
            input: "data: test data 1\n\ndata: test data 2\n\n",
            expectedEvents: []streams.DataEvent{
                {
                    Event: "data", 
                    Data:  []byte("test data 1"),
                },
                {
                    Event: "data",
                    Data:  []byte("test data 2"),
                },
            },
        },
        {
            name:  "named event",
            input: "event: test\ndata: test data\n\n",
            expectedEvents: []streams.DataEvent{
                {
                    Event: "test",
                    Data:  []byte("test data"),
                },
            },
        },
        {
            name:  "invalid field",
            input: "invalid: test\n\n",
            expectedEvents: []streams.DataEvent{
                {
                    Event: "error",
                    Data:  []byte("unknown field in server event"),
                },
            },
        },
        {
            name:  "empty data", 
            input: "data: \n\n",
            expectedEvents: []streams.DataEvent{
                {
                    Event: "data",
                    Data:  []byte(""),
                }, 
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            reader := strings.NewReader(tc.input)
            stream := streams.ConsumeStream(io.NopCloser(reader))

            received := make([]streams.DataEvent, 0)
            for event := range stream {
                received = append(received, event)
            }

            if len(received) != len(tc.expectedEvents) {
                t.Fatalf("expected %d events, got %d", len(tc.expectedEvents), len(received))
            }

            for i, expected := range tc.expectedEvents {
                if received[i].Event != expected.Event {
                    t.Errorf("expected event %q, got %q", expected.Event, received[i].Event)
                }
                if !bytes.Equal(received[i].Data, expected.Data) {
                    t.Errorf("expected data %q, got %q", string(expected.Data), string(received[i].Data))
                }
            }
        })
    }
}

func TestCompletionStreamConverter(t *testing.T) {
    expected := []models.CompletionResponse{
        {
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
        },
        {
            Choices: []struct {
                Text         string        `json:"text"`
                Index        int           `json:"index"`
                LogProbs     []float32     `json:"logprobs"` 
                FinishReason string        `json:"finish_reason"`
            }{
                {
                    Text:         "completion",
                    FinishReason: "stop",
                },  
            },
        },
    }

    input := []streams.DataEvent{
        {
            Event: "data",
            Data:  mustMarshal(expected[0]), 
        },
        {
            Event: "data",
            Data:  mustMarshal(expected[1]),
        },
    }

    c := make(chan streams.DataEvent, len(input))
    for _, e := range input {
        c <- e
    }
    close(c)

    completions, errors := streams.CompletionStreamConverter(c)

    var received []models.CompletionResponse
    for completion := range completions {
        received = append(received, completion)
    }

    if len(received) != len(expected) {
        t.Fatalf("expected %d completions, got %d", len(expected), len(received))
    }

    for _, err := range errors {
        t.Errorf("unexpected error: %v", err)
    }

    for i, exp := range expected {
        if received[i].Choices[0].Text != exp.Choices[0].Text {
            t.Errorf("expected text %q, got %q", exp.Choices[0].Text, received[i].Choices[0].Text)
        }
        if received[i].Choices[0].FinishReason != exp.Choices[0].FinishReason {
            t.Errorf("expected finish reason %q, got %q", exp.Choices[0].FinishReason, received[i].Choices[0].FinishReason)
        }
    }
}

func TestCompletionStreamConverterError(t *testing.T) {
    input := []streams.DataEvent{
        {
            Event: "error",
            Data:  []byte("test error"),
        },
    }

    c := make(chan streams.DataEvent, len(input))
    for _, e := range input {
        c <- e
    }
    close(c)

    completions, errors := streams.CompletionStreamConverter(c)

    var receivedErrors []error
    for err := range errors {
        receivedErrors = append(receivedErrors, err)
    }

    if len(receivedErrors) != 1 {
        t.Fatalf("expected 1 error, got %d", len(receivedErrors))
    }

    if receivedErrors[0].Error() != "test error" {
        t.Errorf("expected error %q, got %q", "test error", receivedErrors[0].Error())
    }

    for range completions {
        t.Error("unexpected completion")
    }
}

func mustMarshal(v interface{}) []byte {
    data, err := json.Marshal(v)
    if err != nil {
        panic(err)
    }
    return data
}

// ... MessageStreamConverter tests similar to CompletionStreamConverter ...

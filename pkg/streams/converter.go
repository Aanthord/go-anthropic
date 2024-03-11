// pkg/streams/converter.go
package streams

import (
    "encoding/json"

    "github.com/Aanthord/anthropic-go/models"
)

func CompletionStreamConverter(input <-chan DataEvent) <-chan models.CompletionResponse {
    output := make(chan models.CompletionResponse)

    go func() {
        defer close(output)
        for event := range input {
            switch event.Event {
            case "completion":
                var completion models.CompletionResponse
                json.Unmarshal(event.Data, &completion)
                output <- completion
            case "error":
                var streamErr errors.StreamError
                json.Unmarshal(event.Data, &streamErr)
                output <- models.CompletionResponse{Completion: streamErr.Error()}
            }
        }  
    }()

    return output
}

func MessageStreamConverter(input <-chan DataEvent) <-chan models.MessageResponse {
    output := make(chan models.

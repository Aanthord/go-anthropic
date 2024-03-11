package streams

import (
    "bufio"
    "bytes"
    "encoding/json"
    "errors"
    "io"
)

type DataEvent struct {
    Data  []byte
    Event string
}

func ConsumeStream(stream io.ReadCloser) <-chan DataEvent {
    outputCh := make(chan DataEvent)

    go func() {
        defer stream.Close()
        defer close(outputCh)

        scanner := bufio.NewScanner(stream)

        for scanner.Scan() {
            line := scanner.Bytes()
            parts := bytes.SplitN(line, []byte(": "), 2)

            if len(parts) < 2 {
                outputCh <- DataEvent{
                    Event: "error",
                    Data:  []byte("invalid server event"),  
                }
                continue
            }

            field, value := string(parts[0]), parts[1]

            switch field {
            case "data":  
                outputCh <- DataEvent{
                    Event: "data",
                    Data:  value,
                }
            case "event":
                outputCh <- DataEvent{
                    Event: string(value),
                }
            default:
                outputCh <- DataEvent{
                    Event: "error", 
                    Data:  []byte("unknown field in server event"),
                }
            }
        }

        if err := scanner.Err(); err != nil {
            outputCh <- DataEvent{
                Event: "error",
                Data:  []byte(err.Error()), 
            }
        }
    }()

    return outputCh
}

func CompletionStreamConverter(inputCh <-chan DataEvent) (<-chan CompletionResponse, <-chan error) {
    outputCh := make(chan CompletionResponse)
    errCh := make(chan error)

    go func() {
        defer close(outputCh)
        defer close(errCh)

        for event := range inputCh {
            switch event.Event {
            case "error":
                errCh <- errors.New(string(event.Data)) 
            case "data":
                var completion CompletionResponse
                if err := json.Unmarshal(event.Data, &completion); err != nil {
                    errCh <- err
                } else {
                    outputCh <- completion  
                }
            }
        }
    }()

    return outputCh, errCh
}

func MessageStreamConverter(inputCh <-chan DataEvent) (<-chan MessageResponse, <-chan error) {
    outputCh := make(chan MessageResponse)
    errCh := make(chan error) 

    go func() {
        defer close(outputCh)
        defer close(errCh)

        for event := range inputCh {
            switch event.Event {
            case "error":
                errCh <- errors.New(string(event.Data))
            case "data":  
                var message MessageResponse
                if err := json.Unmarshal(event.Data, &message); err != nil {
                    errCh <- err  
                } else {
                    outputCh <- message
                }  
            }
        }
    }()

    return outputCh, errCh  
}

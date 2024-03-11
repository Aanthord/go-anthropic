package main

import (
    "context"
    "fmt"

    "github.com/Aanthord/anthropic-go/pkg/api"
    "github.com/Aanthord/anthropic-go/pkg/models"
)

func main() {
    client := api.NewClient("api-key")

    ctx := context.Background()

    req := &models.CompletionRequest{
        Prompt:  "Once upon a time",
        Model:   "claude-1",
    }

    resp, err := client.CreateCompletion(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Println(resp.Choices[0].Text)

    streamReq := &models.CompletionRequest{
        Prompt: "It was a dark and stormy night",
        Model:  "claude-1", 
        Stream: true,
    }

    stream, errStream := client.StreamCompletions(ctx, streamReq)
    for {
        select {
        case resp := <-stream:
            fmt.Print(resp.Choices[0].Text)
        case err := <-errStream:
            fmt.Printf("Error: %v\n", err)
            return  
        }
    }
}

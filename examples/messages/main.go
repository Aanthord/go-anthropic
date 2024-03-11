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

    req := &models.MessageRequest{
        Messages: []models.Message{
            {
                Role:    models.SystemRole,
                Content: "You are a friendly and helpful AI assistant.",  
            },
            {
                Role:    models.UserRole,
                Content: "Hello, how are you today?", 
            },  
        },
        Model: "claude-1",
    }

    resp, err := client.CreateMessage(ctx, req)
    if err != nil {
        panic(err)  
    }

    fmt.Println(resp.Choices[0].Message.Content)

    streamReq := &models.MessageRequest{
        Messages: []models.Message{
            {
                Role:    models.UserRole,
                Content: "What's your favorite book, and why?",
            },
        },
        Model:  "claude-1", 
        Stream: true,  
    }

    stream, errStream := client.StreamMessages(ctx, streamReq)
    for {
        select {
        case resp := <-stream:
            fmt.Print(resp.Choices[0].Message.Content)  
        case err := <-errStream:
            fmt.Printf("Error: %v\n", err)  
            return
        }
    }
}

package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/Aanthord/anthropic-go/pkg/api"
    "github.com/Aanthord/anthropic-go/pkg/models"
)

func main() {
    apiKey := flag.String("api-key", "", "Anthropic API key")
    model := flag.String("model", "", "Model to use for completion/chat")
    prompt := flag.String("prompt", "", "Prompt to send for completion/chat")
    stream := flag.Bool("stream", false, "Stream the completion/chat response")
    chat := flag.Bool("chat", false, "Use chat completion instead of text completion")
    flag.Parse()

    if *apiKey == "" {
        fmt.Println("Please provide an API key with -api-key")
        os.Exit(1)
    }

    if *prompt == "" {
        fmt.Println("Please provide a prompt with -prompt")
        os.Exit(1)
    }

    client := api.NewClient(*apiKey)

    ctx := context.Background()

    if *chat {
        if *stream {
            req := &models.MessageRequest{
                Messages: []models.Message{
                    {
                        Role:    models.UserRole,
                        Content: *prompt,
                    },
                },
                Model:  *model,
                Stream: true,
            }

            respStream, errStream := client.StreamMessages(ctx, req)
            for {
                select {
                case resp := <-respStream:
                    fmt.Printf("Assistant: %s", resp.Choices[0].Message.Content)
                case err := <-errStream:
                    fmt.Printf("Error: %v\n", err)
                    return
                }
            }
        } else {
            req := &models.MessageRequest{
                Messages: []models.Message{
                    {
                        Role:    models.UserRole,
                        Content: *prompt,
                    },
                },
                Model: *model,
            }

            resp, err := client.CreateMessage(ctx, req)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                os.Exit(1)
            }

            fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
        }
    } else {
        if *stream {
            req := &models.CompletionRequest{
                Prompt:  *prompt,
                Model:   *model,
                Stream:  true,
            }

            respStream, errStream := client.StreamCompletions(ctx, req)
            for {
                select {
                case resp := <-respStream:
                    fmt.Printf(resp.Choices[0].Text)
                case err := <-errStream:
                    fmt.Printf("Error: %v\n", err)
                    return
                }
            }
        } else {
            req := &models.CompletionRequest{
                Prompt: *prompt,
                Model:  *model,
            }

            resp, err := client.CreateCompletion(ctx, req)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                os.Exit(1)
            }

            fmt.Println(resp.Choices[0].Text)
        }
    }
}

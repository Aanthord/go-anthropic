// pkg/models/models.go
package models

type CompletionRequest struct {
    Prompt       string   `json:"prompt"`
    MaxTokens    int      `json:"max_tokens"`
    Temperature  float32  `json:"temperature"`
    TopP         float32  `json:"top_p"`
    N            int      `json:"n"`
    Stop         []string `json:"stop"`  
    LogProbs     int      `json:"logprobs"`
    Echo         bool     `json:"echo"`
    Stream       bool     `json:"stream"`
    BestOf       int      `json:"best_of"`
    FrequencyPenalty float32 `json:"frequency_penalty"`
    PresencePenalty float32 `json:"presence_penalty"`
}

type CompletionResponse struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    Model   string `json:"model"`  
    Choices []struct {
        Text         string        `json:"text"`
        Index        int           `json:"index"`
        LogProbs     []float32     `json:"logprobs"`
        FinishReason string        `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        PromptTokens int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens  int `json:"total_tokens"`   
    } `json:"usage"`
}

type MessageRoleType string

const (
    SystemRole    MessageRoleType = "system"
    AssistantRole MessageRoleType = "assistant"
    UserRole      MessageRoleType = "user"   
)

type Message struct {
    Role    MessageRoleType `json:"role"`
    Content string          `json:"content"`
    Name    string          `json:"name,omitempty"`
}

type MessageRequest struct {
    Messages     []Message     `json:"messages"`
    MaxTokens    int           `json:"max_tokens"`  
    N            int           `json:"n"`
    Stop         []string      `json:"stop"`
    Temperature  float32       `json:"temperature"`
    TopP         float32       `json:"top_p"`
    Stream       bool          `json:"stream"`  
    FrequencyPenalty float32   `json:"frequency_penalty"`
    PresencePenalty float32    `json:"presence_penalty"`

type MessageResponse struct {
    ID      string         `json:"id"`
    Object  string         `json:"object"`  
    Created int64          `json:"created"`
    Model   string         `json:"model"`
    Choices []MessageChoice `json:"choices"`  
    Usage   MessageUsage   `json:"usage"`
}

type MessageChoice struct {
    Index        int     `json:"index"`
    Message      Message `json:"message"`
    FinishReason string  `json:"finish_reason"`  
}

type MessageUsage struct {
    PromptTokens     int `json:"prompt_tokens"`  
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`
}

type Model struct {
    ID           string   `json:"id"`
    Object       string   `json:"object"`
    OwnedBy      string   `json:"owned_by"`
    Permissions  []string `json:"permissions"`
    Root         string   `json:"root"`  
    Parent       string   `json:"parent,omitempty"`
    Created      int64    `json:"created"`
    LastModified int64    `json:"last_modified"`
    Deleted      bool     `json:"deleted"`
}

type ModelList struct {
    Models []Model `json:"data"`  
    Object string  `json:"object"`
    Limit  int     `json:"limit"` 
    Offset int     `json:"offset"`
    Total  int     `json:"total"`  
}

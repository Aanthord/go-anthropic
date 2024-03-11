package errors

import "fmt"

type APIError struct {
    StatusCode int    `json:"status_code"`
    Message    string `json:"message"` 
}

func (e APIError) Error() string {
    return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

type StreamError struct {
    Message string
}

func (e StreamError) Error() string {
    return fmt.Sprintf("Stream Error: %s", e.Message)  
}

// test/api/messages_test.go
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

// ... CreateMessage and StreamMessages tests similar to Completions ...


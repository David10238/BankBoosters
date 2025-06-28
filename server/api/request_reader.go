package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestReader struct {
	req *http.Request
}

func newRequestReader(req *http.Request) *RequestReader {
	return &RequestReader{
		req: req,
	}
}

func (reader *RequestReader) BindJsonHeader(key string, v any) *MessageResponse {
	header := reader.req.Header.Get(key)

	if header == "" {
		message := fmt.Sprintf("Header %s is empty", key)
		return SendNotFound(message)
	}

	if err := json.Unmarshal([]byte(header), v); err != nil {
		message := fmt.Sprintf("Header %s is malformed", key)
		return SendNotFound(message)
	}

	return nil
}

func (reader *RequestReader) BindStringHeader(key string, v *string) *MessageResponse {
	header := reader.req.Header.Get(key)

	if header == "" {
		message := fmt.Sprintf("Header %s is empty", key)
		return SendNotFound(message)
	}

	*v = header

	return nil
}

func (reader *RequestReader) BindJsonBody(v any) *MessageResponse {
	if err := json.NewDecoder(reader.req.Body).Decode(v); err != nil {
		return SendNotFound("Not able to bind body")
	}

	return nil
}

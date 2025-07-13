package protocol

import (
	"encoding/json"
)

type transferHeader struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Count string `json:"count"`
}

type itemHeader struct {
	Type string          `json:"type"`
	Meta json.RawMessage `json:"meta"`
}

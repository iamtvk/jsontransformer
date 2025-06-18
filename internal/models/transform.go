package models

import (
	"encoding/json"
	"time"
)

type TransformerRequest struct {
	Data             json.RawMessage `json:"data" validate:"required"`
	ScriptIdentifier string          `json:"script_identifier"`
	Script           string          `json:"script,omitempty"`
	Timeout          time.Duration   `json:"timeout,omitempty"`
}

type CreateScriptRequest struct {
	ScriptIdentifier string `json:"script_identifier" validate:"required"`
	Script           string `json:"script" validate:"required"`
	CreatedBy        string `json:"created_by"`
	Description      string `json:"description"`
	Name             string `json:"name"`
}

type TransformerResponse struct {
	Result   json.RawMessage     `json:"result"`
	Error    *TransformerError   `json:"error"`
	MetaData TransformerMetadata `json:"metadata"`
}

type TransformerError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
}

type TransformerMetadata struct {
	ScriptIdentifier string        `json:"script_identifier"`
	ScriptVersion    int           `json:"script_version"`
	ExecutionTime    time.Duration `json:"execution_time"`
	InputSize        int           `json:"input_size"`
	OutputSize       int           `json:"output_size"`
	CacheHit         bool          `json:"cache_hit"`
}

type CreateScriptResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}
